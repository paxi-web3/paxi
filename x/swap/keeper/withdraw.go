package keeper

import (
	"encoding/json"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/x/swap/types"
)

func (k Keeper) WithdrawLiquidity(ctx sdk.Context, msg *types.MsgWithdrawLiquidity) error {
	// Parse and validate the creator address
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return fmt.Errorf("invalid creator: %w", err)
	}

	// Parse and validate LP token amount
	lpAmount, ok := sdkmath.NewIntFromString(msg.LpAmount)
	if !ok || !lpAmount.IsPositive() {
		return fmt.Errorf("invalid LP amount: %s", msg.LpAmount)
	}

	// Retrieve the pool
	pool, found := k.GetPool(ctx, msg.Prc20)
	if !found || pool.TotalShares.IsZero() {
		return fmt.Errorf("pool not found or empty")
	}

	// Retrieve user's position
	position, found := k.GetPosition(ctx, msg.Prc20, creator)
	iCurrentLpAmount, ok := sdkmath.NewIntFromString(position.LpAmount)
	if !ok {
		return fmt.Errorf("invalid LP amount in position")
	}
	if !found || iCurrentLpAmount.LT(lpAmount) {
		return fmt.Errorf("insufficient LP token balance")
	}

	// Calculate user's share ratio and withdrawable assets
	shareRatio := sdkmath.LegacyNewDecFromInt(lpAmount).Quo(sdkmath.LegacyNewDecFromInt(pool.TotalShares))
	paxiOut := sdkmath.LegacyNewDecFromInt(pool.ReservePaxi).Mul(shareRatio).TruncateInt()
	prc20Out := sdkmath.LegacyNewDecFromInt(pool.ReservePRC20).Mul(shareRatio).TruncateInt()

	// Update pool reserves and total shares
	pool.ReservePaxi = pool.ReservePaxi.Sub(paxiOut)
	pool.ReservePRC20 = pool.ReservePRC20.Sub(prc20Out)
	pool.TotalShares = pool.TotalShares.Sub(lpAmount)
	k.SetPool(ctx, pool)

	// Burn LP tokens: transfer from user to module, then burn
	lpDenom := types.LPTokenDenom(msg.Prc20)
	lpCoin := sdk.NewCoin(lpDenom, lpAmount)
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, creator, types.ModuleName, sdk.NewCoins(lpCoin))
	if err != nil {
		return fmt.Errorf("failed to collect LP token: %w", err)
	}
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(lpCoin))
	if err != nil {
		return fmt.Errorf("burn failed: %w", err)
	}

	// Send PAXI back to the user
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creator,
		sdk.NewCoins(sdk.NewCoin(types.DefaultDenom, paxiOut)))
	if err != nil {
		return fmt.Errorf("failed to send PAXI: %w", err)
	}

	// Send PRC20 back to the user (via wasm execute)
	err = k.transferPRC20FromModule(ctx, msg.Prc20, creator, prc20Out)
	if err != nil {
		return fmt.Errorf("failed to send PRC20: %w", err)
	}

	// Update or delete user's position record
	if iCurrentLpAmount.Equal(lpAmount) {
		k.DeletePosition(ctx, msg.Prc20, creator)
	} else {
		position.LpAmount = iCurrentLpAmount.Sub(lpAmount).String()
		k.SetPosition(ctx, position)
	}

	return nil
}

func (k Keeper) transferPRC20FromModule(ctx sdk.Context, contract string, to sdk.AccAddress, amount sdkmath.Int) error {
	msg := map[string]interface{}{
		"transfer": map[string]interface{}{
			"recipient": to.String(),
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
		nil)
	return err
}
