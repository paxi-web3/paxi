package keeper

import (
	"encoding/json"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/x/swap/types"
)

// ProvideLiquidity processes a liquidity provision message. If the pool does not exist, it is created.
func (k Keeper) ProvideLiquidity(ctx sdk.Context, msg *types.MsgProvideLiquidity) error {
	var recoveredErr error
	defer func() {
		if r := recover(); r != nil {
			ctx.Logger().Error("swap panic recovered", "err", r)
			recoveredErr = fmt.Errorf("panic: %v", r)
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

	// Transfer tokens to module
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, creator, types.ModuleName, sdk.NewCoins(sdk.NewCoin(types.DefaultDenom, paxiAmt)))
	if err != nil {
		return err
	}
	err = k.transferPRC20(ctx, msg.Creator, msg.Prc20, prc20Amt)
	if err != nil {
		return err
	}

	// Load or init pool
	pool, found := k.GetPool(ctx, msg.Prc20)

	var lpToMint sdkmath.Int
	if !found {
		// Initial liquidity provider, direct minting LP = paxiAmt
		lpToMint = paxiAmt
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

		// Subsequent liquidity providers will mint LP tokens in proportion
		share1 := paxiAmt.Mul(pool.TotalShares).Quo(pool.ReservePaxi)
		share2 := prc20Amt.Mul(pool.TotalShares).Quo(pool.ReservePRC20)
		lpToMint = sdkmath.MinInt(share1, share2)

		if lpToMint.IsZero() {
			return fmt.Errorf("provided liquidity too small to mint LP token")
		}

		pool.ReservePaxi = pool.ReservePaxi.Add(paxiAmt)
		pool.ReservePRC20 = pool.ReservePRC20.Add(prc20Amt)
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
		lpAmount, _ := sdkmath.NewIntFromString(pos.LpAmount)
		pos.LpAmount = lpAmount.Add(lpToMint).String()
	}
	k.SetPosition(ctx, pos)

	return recoveredErr
}

// transferPRC20 performs a transfer_from call on a PRC20 contract
// to move tokens from the user to the swap module account.
func (k Keeper) transferPRC20(ctx sdk.Context, from string, contract string, amount sdkmath.Int) error {
	msg := map[string]interface{}{
		"transfer_from": map[string]interface{}{
			"owner":     from,
			"recipient": k.accountKeeper.GetModuleAddress(types.ModuleName).String(),
			"amount":    amount.String(),
		},
	}
	bz, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = k.wasmKeeper.Execute(ctx,
		sdk.MustAccAddressFromBech32(contract),
		k.accountKeeper.GetModuleAddress(types.ModuleName),
		bz,
		nil, // no attached funds
	)
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
