package keeper

import (
	"encoding/json"
	"fmt"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/x/swap/types"
)

func (k Keeper) WithdrawLiquidity(ctx sdk.Context, msg *types.MsgWithdrawLiquidity) error {
	defer func() {
		if r := recover(); r != nil {
			ctx.Logger().Error("swap panic recovered", "err", r)
			panic(fmt.Errorf("swap module panic recovered: %v", r))
		}
	}()

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

	if paxiOut.IsZero() && prc20Out.IsZero() {
		return fmt.Errorf("withdrawal too small, results in zero output")
	}

	// Send PAXI and PRC20 back to user
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creator, sdk.NewCoins(sdk.NewCoin(types.DefaultDenom, paxiOut))); err != nil {
		return fmt.Errorf("failed to send PAXI: %w", err)
	}
	if err := k.transferPRC20FromModule(ctx, msg.Prc20, creator, prc20Out); err != nil {
		return fmt.Errorf("failed to send PRC20: %w", err)
	}

	// Update pool reserves and total shares
	if pool.ReservePaxi.LT(paxiOut) || pool.ReservePRC20.LT(prc20Out) || pool.TotalShares.LT(lpAmount) {
		return fmt.Errorf("pool reserve insufficient for withdrawal")
	}
	pool.ReservePaxi = pool.ReservePaxi.Sub(paxiOut)
	pool.ReservePRC20 = pool.ReservePRC20.Sub(prc20Out)
	pool.TotalShares = pool.TotalShares.Sub(lpAmount)

	// Update or delete user's position
	newLpAmount := iCurrentLpAmount.Sub(lpAmount)
	if newLpAmount.IsZero() {
		k.DeletePosition(ctx, msg.Prc20, creator)
	} else {
		position.LpAmount = newLpAmount.String()
		k.SetPosition(ctx, position)
	}

	// If pool is empty after withdrawal, delete it
	if pool.TotalShares.IsZero() {
		k.DeletePool(ctx, pool.Prc20)
	} else {
		k.SetPool(ctx, pool)
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

	// Set safty boundary to prevent infinite loop
	safeGas := uint64(500_000)

	// Detect if this is a real tx or simulation
	isSimulate := ctx.IsCheckTx() || ctx.IsReCheckTx()

	// Set safe gas only when it's not simulating
	execCtx := ctx
	if !isSimulate {
		execCtx = ctx.WithGasMeter(storetypes.NewGasMeter(safeGas)) // only limit in DeliverTx
	}

	_, err = k.wasmKeeper.Execute(
		execCtx,
		sdk.MustAccAddressFromBech32(contract),
		k.accountKeeper.GetModuleAddress(types.ModuleName),
		bz,
		nil)
	return err
}
