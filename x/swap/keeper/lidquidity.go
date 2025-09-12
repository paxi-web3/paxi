package keeper

import (
	"encoding/json"
	"fmt"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/utils"
	"github.com/paxi-web3/paxi/x/swap/types"
)

// ProvideLiquidity processes a liquidity provision message. If the pool does not exist, it is created.
func (k Keeper) ProvideLiquidity(ctx sdk.Context, msg *types.MsgProvideLiquidity) error {
	defer func() {
		if r := recover(); r != nil {
			ctx.Logger().Error("swap panic recovered", "err", r)
			panic(fmt.Errorf("swap module panic recovered: %v", r))
		}
	}()

	params := k.GetParams(ctx)

	// Validate PRC20 contract
	contractAddr := sdk.MustAccAddressFromBech32(msg.Prc20)
	contractInfo := k.wasmQueryKeeper.GetContractInfo(ctx, contractAddr)
	if contractInfo == nil {
		return fmt.Errorf("PRC20 contract %s not found", msg.Prc20)
	}
	if contractInfo.CodeID != params.CodeID {
		return fmt.Errorf("PRC20 contract %s has invalid CodeID: %d", msg.Prc20, contractInfo.CodeID)
	}

	// Parse amounts
	paxiAmount, err := sdk.ParseCoinNormalized(msg.PaxiAmount)
	if err != nil {
		return fmt.Errorf("invalid paxi amount: %w", err)
	}
	paxiAmt := paxiAmount.Amount
	if paxiAmount.Denom != types.DefaultDenom {
		return fmt.Errorf("only %s accepted as base token", types.DefaultDenom)
	}

	prc20Amt, ok := sdkmath.NewIntFromString(msg.Prc20Amount)
	if !ok {
		return fmt.Errorf("invalid prc20 amount: %s", msg.Prc20Amount)
	}

	if paxiAmt.LT(sdkmath.NewInt(int64(params.MinLiquidity))) || prc20Amt.LT(sdkmath.NewInt(int64(params.MinLiquidity))) {
		return fmt.Errorf("below minimum liquidity: %d", params.MinLiquidity)
	}

	// Convert creator address
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}

	// Load or init pool
	pool, found := k.GetPool(ctx, msg.Prc20)

	var lpToMint sdkmath.Int
	if !found {
		// Transfer tokens to module
		err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, creator, types.ModuleName, sdk.NewCoins(sdk.NewCoin(types.DefaultDenom, paxiAmt)))
		if err != nil {
			return err
		}
		err = k.transferPRC20(ctx, msg.Creator, msg.Prc20, prc20Amt)
		if err != nil {
			return err
		}

		// Initial liquidity provider
		lpToMint = utils.IntSqrt(paxiAmt.Mul(prc20Amt))
		if lpToMint.IsZero() {
			return fmt.Errorf("provided liquidity too small to mint LP token")
		}

		// Create new pool
		pool = types.Pool{
			Prc20:        msg.Prc20,
			ReservePaxi:  paxiAmt,
			ReservePRC20: prc20Amt,
			TotalShares:  lpToMint,
		}
	} else {
		if pool.ReservePaxi.IsZero() || pool.ReservePRC20.IsZero() || pool.TotalShares.IsZero() {
			return fmt.Errorf("corrupted pool state, cannot join")
		}

		// Calculate actual needed amounts
		requiredPaxi := prc20Amt.Mul(pool.ReservePaxi).Quo(pool.ReservePRC20)
		requiredPrc20 := paxiAmt.Mul(pool.ReservePRC20).Quo(pool.ReservePaxi)
		usedPaxi := requiredPaxi
		usedPrc20 := prc20Amt

		if requiredPaxi.GT(paxiAmt) {
			usedPrc20 = requiredPrc20
			usedPaxi = paxiAmt
		}

		if usedPaxi.IsZero() || usedPrc20.IsZero() {
			return fmt.Errorf("provided liquidity too small to mint LP token")
		}

		// Subsequent liquidity providers will mint LP tokens in proportion
		share1 := usedPaxi.Mul(pool.TotalShares).Quo(pool.ReservePaxi)
		share2 := usedPrc20.Mul(pool.TotalShares).Quo(pool.ReservePRC20)
		lpToMint = sdkmath.MinInt(share1, share2)

		// Transfer tokens to module
		err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, creator, types.ModuleName, sdk.NewCoins(sdk.NewCoin(types.DefaultDenom, usedPaxi)))
		if err != nil {
			return err
		}
		err = k.transferPRC20(ctx, msg.Creator, msg.Prc20, usedPrc20)
		if err != nil {
			return err
		}

		if lpToMint.IsZero() {
			return fmt.Errorf("provided liquidity too small to mint LP token")
		}

		pool.ReservePaxi = pool.ReservePaxi.Add(usedPaxi)
		pool.ReservePRC20 = pool.ReservePRC20.Add(usedPrc20)
		pool.TotalShares = pool.TotalShares.Add(lpToMint)
	}

	// Update pool in KVStore
	k.SetPool(ctx, pool)

	// Update LP token ownership
	pos, found := k.GetPosition(ctx, msg.Prc20, creator)
	if !found {
		pos = types.ProviderPosition{
			Creator:  creator.String(),
			Prc20:    msg.Prc20,
			LpAmount: lpToMint.String(),
		}
	} else {
		lpAmount, ok := sdkmath.NewIntFromString(pos.LpAmount)
		if !ok {
			return fmt.Errorf("invalid lp amount in position: %s", pos.LpAmount)
		}
		pos.LpAmount = lpAmount.Add(lpToMint).String()
	}
	k.SetPosition(ctx, pos)

	return nil
}

// transferPRC20 performs a transfer_from call on a PRC20 contract
// to move tokens from the user to the swap module account.
func (k Keeper) transferPRC20(ctx sdk.Context, from string, contract string, amount sdkmath.Int) error {
	if !amount.IsPositive() {
		return fmt.Errorf("transferPRC20: amount must be positive")
	}

	type transferFrom struct {
		Owner     string `json:"owner"`
		Recipient string `json:"recipient"`
		Amount    string `json:"amount"`
	}
	type msgWrapper struct {
		TransferFrom transferFrom `json:"transfer_from"`
	}

	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	contractAddr, err := sdk.AccAddressFromBech32(contract)
	if err != nil {
		return fmt.Errorf("transferPRC20: invalid contract address: %w", err)
	}

	msg := msgWrapper{
		TransferFrom: transferFrom{
			Owner:     from,
			Recipient: moduleAddr.String(),
			Amount:    amount.String(),
		},
	}
	bz, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("transferPRC20: marshal msg: %w", err)
	}

	// Gas handling strategy
	const safeGas uint64 = 5_000_000
	if ctx.IsCheckTx() || ctx.IsReCheckTx() {
		// In mempool/simulation, do not cap; keep realistic gas accounting.
		_, err = k.wasmKeeper.Execute(ctx, contractAddr, moduleAddr, bz, nil)
		return err
	}

	parent := ctx.GasMeter()
	child := storetypes.NewGasMeter(safeGas)
	execCtx := ctx.WithGasMeter(child)

	// If child runs out of gas, it will panic → tx fails as expected.
	_, err = k.wasmKeeper.Execute(execCtx, contractAddr, moduleAddr, bz, nil)

	// Charge the actual gas consumed by child to the parent,
	// so total tx gas remains accurate and limited by the outer meter.
	// If the parent has insufficient remaining gas, this will panic → OOG, as expected.
	parent.ConsumeGas(child.GasConsumed(), "prc20 transfer_from")

	return err
}

// SetPool saves the pool to KVStore using protobuf
func (k Keeper) SetPool(ctx sdk.Context, pool types.Pool) {
	store := k.storeService.OpenKVStore(ctx)

	poolProto := pool.ToProto()
	bz, err := k.cdc.Marshal(&poolProto)
	if err != nil {
		panic(err)
	}

	err = store.Set(types.PoolStoreKey(pool.Prc20), bz)
	if err != nil {
		panic(err)
	}
}

// GetPool fetches a pool by PRC20 contract address using protobuf
func (k Keeper) GetPool(ctx sdk.Context, prc20 string) (types.Pool, bool) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.PoolStoreKey(prc20))
	if err != nil || bz == nil {
		return types.Pool{}, false
	}

	var poolProto types.PoolProto
	err = k.cdc.Unmarshal(bz, &poolProto)
	if err != nil {
		return types.Pool{}, false
	}

	pool, err := types.PoolFromProto(&poolProto)
	if err != nil {
		return types.Pool{}, false
	}

	return pool, true
}

// DeletePool removes a pool from KVStore by PRC20 contract address
func (k Keeper) DeletePool(ctx sdk.Context, prc20 string) {
	store := k.storeService.OpenKVStore(ctx)
	if err := store.Delete(types.PoolStoreKey(prc20)); err != nil {
		panic(err)
	}
}
