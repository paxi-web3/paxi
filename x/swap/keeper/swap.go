package keeper

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/x/swap/types"
)

func (k Keeper) Swap(ctx sdk.Context, msg *types.MsgSwap) error {
	defer func() {
		if r := recover(); r != nil {
			ctx.Logger().Error("swap panic recovered", "err", r)
			panic(fmt.Errorf("swap module panic recovered: %v", r))
		}
	}()

	params := k.GetParams(ctx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}

	offerAmt, ok := sdkmath.NewIntFromString(msg.OfferAmount)
	if !ok || offerAmt.IsZero() {
		return fmt.Errorf("invalid offer amount")
	}

	minReceive, ok := sdkmath.NewIntFromString(msg.MinReceive)
	if !ok {
		return fmt.Errorf("invalid min_receive")
	}

	pool, found := k.GetPool(ctx, msg.Prc20)
	if !found || pool.ReservePaxi.IsZero() || pool.ReservePRC20.IsZero() {
		return fmt.Errorf("pool not found or empty")
	}

	feeBps := sdkmath.NewInt(int64(params.SwapFeeBPS))
	if feeBps.GTE(sdkmath.NewInt(types.BPSUnit)) {
		return fmt.Errorf("invalid swap fee BPS: %d", feeBps.Int64())
	}

	var inputReserve, outputReserve sdkmath.Int
	var recvAmount sdkmath.Int

	if msg.OfferDenom == types.DefaultDenom {
		// Swap: PAXI -> PRC20
		inputReserve = pool.ReservePaxi
		outputReserve = pool.ReservePRC20

		recvAmount = getAmountOut(offerAmt, inputReserve, outputReserve, feeBps)
		if recvAmount.IsZero() {
			return fmt.Errorf("swap too small: results in zero output")
		}
		if recvAmount.LT(minReceive) {
			return fmt.Errorf("slippage too high")
		}
		if pool.ReservePRC20.LT(recvAmount) {
			return fmt.Errorf("insufficient PRC20 liquidity")
		}

		// Enforce module‐level max‐swap ratio
		maxOut := outputReserve.
			MulRaw(int64(types.MaxSwapRatioBPS)).
			QuoRaw(types.BPSUnit)
		if recvAmount.GT(maxOut) {
			return fmt.Errorf(
				"requested output %s exceeds module max %d BPS of reserve (%s)",
				recvAmount.String(), types.MaxSwapRatioBPS, maxOut.String(),
			)
		}

		// Transfer PAXI in, PRC20 out
		err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, creator, types.ModuleName,
			sdk.NewCoins(sdk.NewCoin(types.DefaultDenom, offerAmt)))
		if err != nil {
			return fmt.Errorf("failed to send PAXI: %w", err)
		}
		err = k.transferPRC20FromModule(ctx, msg.Prc20, creator, recvAmount)
		if err != nil {
			return fmt.Errorf("failed to send PRC20: %w", err)
		}

		// Update reserves
		pool.ReservePaxi = pool.ReservePaxi.Add(offerAmt)
		pool.ReservePRC20 = pool.ReservePRC20.Sub(recvAmount)

	} else if msg.OfferDenom == msg.Prc20 {
		// Swap: PRC20 -> PAXI
		inputReserve = pool.ReservePRC20
		outputReserve = pool.ReservePaxi

		recvAmount = getAmountOut(offerAmt, inputReserve, outputReserve, feeBps)
		if recvAmount.IsZero() {
			return fmt.Errorf("swap too small: results in zero output")
		}
		if recvAmount.LT(minReceive) {
			return fmt.Errorf("slippage too high")
		}
		if pool.ReservePaxi.LT(recvAmount) {
			return fmt.Errorf("insufficient PAXI liquidity")
		}

		// Enforce module‐level max‐swap ratio
		maxOut := outputReserve.
			MulRaw(int64(types.MaxSwapRatioBPS)).
			QuoRaw(types.BPSUnit)
		if recvAmount.GT(maxOut) {
			return fmt.Errorf(
				"requested output %s exceeds module max %d BPS of reserve (%s)",
				recvAmount.String(), types.MaxSwapRatioBPS, maxOut.String(),
			)
		}

		// Transfer PRC20 in, PAXI out
		err = k.transferPRC20(ctx, msg.Creator, msg.Prc20, offerAmt)
		if err != nil {
			return fmt.Errorf("failed to send PRC20: %w", err)
		}
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creator,
			sdk.NewCoins(sdk.NewCoin(types.DefaultDenom, recvAmount)))
		if err != nil {
			return fmt.Errorf("failed to send PAXI: %w", err)
		}

		pool.ReservePRC20 = pool.ReservePRC20.Add(offerAmt)
		pool.ReservePaxi = pool.ReservePaxi.Sub(recvAmount)

	} else {
		return fmt.Errorf("invalid offer denom: %s", msg.OfferDenom)
	}

	// Swap fee stays in pool, benefits LPs via higher reserve ratios
	k.SetPool(ctx, pool)
	return nil
}

// getAmountOut calculates the output token amount after applying swap fee and AMM formula
func getAmountOut(amountIn, reserveIn, reserveOut, feeBps sdkmath.Int) sdkmath.Int {
	feeDenom := sdkmath.NewInt(types.BPSUnit)
	inputWithFee := amountIn.Mul(feeDenom.Sub(feeBps)).Quo(feeDenom)
	numerator := inputWithFee.Mul(reserveOut)
	denominator := reserveIn.Add(inputWithFee)
	return numerator.Quo(denominator)
}
