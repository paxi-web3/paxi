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
	params := k.GetParams(ctx)

	// Validate PRC20 contract
	contractAddr := sdk.MustAccAddressFromBech32(msg.Prc20)
	contractInfo := k.wasmQueryKeeper.GetContractInfo(ctx, contractAddr)
	if contractInfo == nil || contractInfo.CodeID != params.CodeID {
		return fmt.Errorf("PRC20 contract %s not found or invalid CodeID: %d", msg.Prc20, contractInfo.CodeID)
	}

	// Parse amounts
	paxiAmt := msg.PaxiAmount.Amount
	if msg.PaxiAmount.Denom != types.DefaultDenom {
		return fmt.Errorf("only %s accepted as base token", types.DefaultDenom)
	}

	prc20Amt, ok := sdkmath.NewIntFromString(msg.Prc20Amount)
	if !ok {
		return fmt.Errorf("invalid prc20 amount: %s", msg.Prc20Amount)
	}

	if paxiAmt.LT(sdkmath.NewInt(int64(params.MinLidquidity))) || prc20Amt.LT(sdkmath.NewInt(int64(params.MinLidquidity))) {
		return fmt.Errorf("below minimum liquidity: %d", params.MinLidquidity)
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
	lpDenom := types.LPTokenDenom(msg.Prc20)

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
		// Subsequent liquidity providers will mint LP tokens in proportion
		share1 := paxiAmt.Mul(pool.TotalShares).Quo(pool.ReservePaxi)
		share2 := prc20Amt.Mul(pool.TotalShares).Quo(pool.ReservePRC20)
		lpToMint = sdkmath.MinInt(share1, share2)

		pool.ReservePaxi = pool.ReservePaxi.Add(paxiAmt)
		pool.ReservePRC20 = pool.ReservePRC20.Add(prc20Amt)
		pool.TotalShares = pool.TotalShares.Add(lpToMint)
	}

	// Mint LP token to user
	lpCoin := sdk.NewCoin(lpDenom, lpToMint)
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(lpCoin)); err != nil {
		return fmt.Errorf("failed to mint LP token: %w", err)
	}
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creator, sdk.NewCoins(lpCoin)); err != nil {
		return fmt.Errorf("failed to send LP token: %w", err)
	}

	k.SetPool(ctx, pool)

	// Update LP token ownership
	pos, found := k.GetPosition(ctx, msg.Prc20, creator)
	if !found {
		pos = types.ProviderPosition{
			Creator:     creator.String(),
			Prc20:       msg.Prc20,
			LpAmount:    lpToMint.String(),
			DepositedLp: lpToMint.String(),
		}
	} else {
		lpAmount, _ := sdkmath.NewIntFromString(pos.LpAmount)
		pos.LpAmount = lpAmount.Add(lpToMint).String()
		depositedLp, _ := sdkmath.NewIntFromString(pos.DepositedLp)
		pos.DepositedLp = depositedLp.Add(lpToMint).String()
	}
	k.SetPosition(ctx, pos)

	return nil
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
		sdk.MustAccAddressFromBech32(from),
		bz,
		nil, // no attached funds
	)
	return err
}
