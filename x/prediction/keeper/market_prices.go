package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/x/prediction/types"
)

// refreshMarketBookPrices recalculates best bid/ask from active LIMIT orders.
func (k Keeper) refreshMarketBookPrices(ctx sdk.Context, market *types.Market) {
	store := k.storeService.OpenKVStore(ctx)
	it, err := prefixIterator(store, types.OrderPrefix)
	if err != nil {
		panic(err)
	}
	defer it.Close()

	var (
		hasBid  bool
		hasAsk  bool
		bestBid sdkmath.Int
		bestAsk sdkmath.Int
	)

	for ; it.Valid(); it.Next() {
		order := &types.Order{}
		if err := k.cdc.Unmarshal(it.Value(), order); err != nil {
			continue
		}
		if order.MarketId != market.Id {
			continue
		}
		if !k.isOrderEffectivelyOpen(ctx, order) {
			continue
		}
		if order.OrderType != types.OrderType_ORDER_TYPE_LIMIT {
			continue
		}

		price, err := types.ParsePriceTicks(order.LimitPrice, "limit_price")
		if err != nil {
			continue
		}

		switch order.Side {
		case types.OrderSide_ORDER_SIDE_BUY_YES, types.OrderSide_ORDER_SIDE_BUY_NO:
			if !hasBid || price.GT(bestBid) {
				bestBid = price
				hasBid = true
			}
		case types.OrderSide_ORDER_SIDE_SELL_YES, types.OrderSide_ORDER_SIDE_SELL_NO:
			if !hasAsk || price.LT(bestAsk) {
				bestAsk = price
				hasAsk = true
			}
		}
	}

	if hasBid {
		market.BestBidPrice = bestBid.String()
	} else {
		market.BestBidPrice = ""
	}
	if hasAsk {
		market.BestAskPrice = bestAsk.String()
	} else {
		market.BestAskPrice = ""
	}
}
