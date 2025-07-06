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

	// Validate that the PRC20 contract exists and is created from the allowed CodeID
	contractInfo := k.wasmQueryKeeper.GetContractInfo(ctx, sdk.MustAccAddressFromBech32(msg.Prc20))
	if contractInfo == nil || contractInfo.CodeID != params.CodeID {
		return fmt.Errorf("PRC20 contract %s not found or invalid CodeID: %d", msg.Prc20, contractInfo.CodeID)
	}

	// Retrieve the pool or initialize a new one
	pool, found := k.GetPool(ctx, msg.Prc20)
	if !found {
		pool = types.Pool{
			Prc20:        msg.Prc20,
			ReservePaxi:  sdkmath.ZeroInt(),
			ReservePRC20: sdkmath.ZeroInt(),
		}
	}

	paxiAmt := msg.PaxiAmount.Amount
	if msg.PaxiAmount.Denom != types.DefaultDenom {
		return fmt.Errorf("only %s accepted as base token", types.DefaultDenom)
	}

	prc20Amt, ok := sdkmath.NewIntFromString(msg.Prc20Amount)
	if !ok {
		return fmt.Errorf("invalid prc20 amount: %s", msg.Prc20Amount)
	}

	// Check minimum liquidity requirement
	if paxiAmt.LT(sdkmath.NewInt(int64(params.MinLidquidity))) || prc20Amt.LT(sdkmath.NewInt(int64(params.MinLidquidity))) {
		return fmt.Errorf("below minimum liquidity: %d", params.MinLidquidity)
	}

	// Transfer PAXI tokens from user to module account
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, creator, types.ModuleName, sdk.NewCoins(sdk.NewCoin(types.DefaultDenom, paxiAmt)))
	if err != nil {
		return err
	}

	// Transfer PRC20 tokens from user to module account using contract
	err = k.transferPRC20(ctx, msg.Creator, msg.Prc20, prc20Amt)
	if err != nil {
		return err
	}

	// Update the pool reserves
	pool.ReservePaxi = pool.ReservePaxi.Add(paxiAmt)
	pool.ReservePRC20 = pool.ReservePRC20.Add(prc20Amt)
	k.SetPool(ctx, pool)

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
