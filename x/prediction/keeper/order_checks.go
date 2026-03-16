package keeper

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/x/prediction/types"
)

func maxUnitPriceFromFields(orderType types.OrderType, limitPrice, worstPrice string) (sdkmath.Int, error) {
	switch orderType {
	case types.OrderType_ORDER_TYPE_LIMIT:
		return types.ParsePriceTicks(limitPrice, "limit_price")
	case types.OrderType_ORDER_TYPE_MARKET:
		return types.ParsePriceTicks(worstPrice, "worst_price")
	default:
		return sdkmath.Int{}, fmt.Errorf("invalid order_type")
	}
}

func precheckSellShares(pos *types.Position, side types.OrderSide, required sdkmath.Int) error {
	yes, lockedYes, no, lockedNo, err := mustPositionIntsStatic(pos)
	if err != nil {
		return err
	}

	switch side {
	case types.OrderSide_ORDER_SIDE_SELL_YES:
		free := yes.Sub(lockedYes)
		if free.LT(required) {
			return fmt.Errorf("insufficient YES shares: required=%s available=%s", required.String(), free.String())
		}
	case types.OrderSide_ORDER_SIDE_SELL_NO:
		free := no.Sub(lockedNo)
		if free.LT(required) {
			return fmt.Errorf("insufficient NO shares: required=%s available=%s", required.String(), free.String())
		}
	}

	return nil
}

func mustPositionIntsStatic(pos *types.Position) (yes sdkmath.Int, lockedYes sdkmath.Int, no sdkmath.Int, lockedNo sdkmath.Int, err error) {
	yes, err = parseNonNegativeInt(pos.YesShares, "yes_shares")
	if err != nil {
		return
	}
	lockedYes, err = parseNonNegativeInt(pos.LockedYesShares, "locked_yes_shares")
	if err != nil {
		return
	}
	no, err = parseNonNegativeInt(pos.NoShares, "no_shares")
	if err != nil {
		return
	}
	lockedNo, err = parseNonNegativeInt(pos.LockedNoShares, "locked_no_shares")
	return
}

func maxBuyCostWithFee(unitPrice sdkmath.Int, amount sdkmath.Int, feeBps uint64) sdkmath.Int {
	cost := unitPrice.Mul(amount)
	feeUnit := unitPrice.MulRaw(int64(feeBps)).QuoRaw(int64(types.BPSDenominator))
	return cost.Add(feeUnit.Mul(amount))
}

func (k Keeper) ensurePlaceOrderCapacity(
	ctx sdk.Context,
	market *types.Market,
	trader sdk.AccAddress,
	side types.OrderSide,
	orderType types.OrderType,
	limitPrice string,
	worstPrice string,
	amount sdkmath.Int,
) error {
	if isBuySide(side) {
		maxUnitPrice, err := maxUnitPriceFromFields(orderType, limitPrice, worstPrice)
		if err != nil {
			return err
		}
		// Conservative check: buyer can afford worst-case execution + fee path.
		required := maxBuyCostWithFee(maxUnitPrice, amount, market.FeeBps)
		return k.ensureCollateralBalance(ctx, market, trader, required)
	}

	if isSellSide(side) {
		pos := k.getPositionOrDefault(ctx, market.Id, trader)
		return precheckSellShares(pos, side, amount)
	}

	return fmt.Errorf("invalid order side")
}
