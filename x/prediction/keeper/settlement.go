package keeper

import (
	"fmt"

	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/x/prediction/types"
)

func (k Keeper) ApplyTradeBatch(ctx sdk.Context, msg *types.MsgApplyTradeBatch) (*types.MsgApplyTradeBatchResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, errors.Wrap(types.ErrInvalidRequest, err.Error())
	}

	params := k.GetParams(ctx)
	if uint64(len(msg.Trades)) > params.MaxBatchSize {
		return nil, fmt.Errorf("batch size exceeds max_batch_size")
	}

	market, found := k.GetMarket(ctx, msg.MarketId)
	if !found {
		return nil, types.ErrMarketNotFound
	}
	if msg.Sender != market.Resolver {
		return nil, types.ErrUnauthorized
	}
	k.maybeCloseExpiredMarket(ctx, market)
	if market.Status != types.MarketStatus_MARKET_STATUS_OPEN {
		return nil, errors.Wrap(types.ErrInvalidMarketStatus, "market must be OPEN for settlement")
	}
	if ctx.BlockTime().Unix() < market.OpenTime {
		return nil, fmt.Errorf("market is not open yet")
	}
	if ctx.BlockTime().Unix() >= market.CloseTime {
		return nil, errors.Wrap(types.ErrInvalidMarketStatus, "market is closed")
	}

	totalFee := sdkmath.ZeroInt()
	yesTotal, noTotal, err := k.getMarketShareInts(market)
	if err != nil {
		return nil, err
	}
	totalTradeVolume := sdkmath.ZeroInt()
	if market.TotalTradeVolume != "" {
		totalTradeVolume, err = parseNonNegativeInt(market.TotalTradeVolume, "total_trade_volume")
		if err != nil {
			return nil, err
		}
	}
	lastTradePrice := ""

	for i := range msg.Trades {
		trade := msg.Trades[i]
		if k.HasAppliedTrade(ctx, market.Id, trade.TradeId) {
			return nil, types.ErrDuplicateTrade
		}

		matchAmount, err := parsePositiveInt(trade.MatchAmount, "trade.match_amount")
		if err != nil {
			return nil, fmt.Errorf("trade[%d]: %w", i, err)
		}
		executionPrice, err := types.ParsePriceTicks(trade.ExecutionPrice, "trade.execution_price")
		if err != nil {
			return nil, fmt.Errorf("trade[%d]: %w", i, err)
		}

		orderA, found := k.GetOrder(ctx, market.Id, trade.OrderAId)
		if !found {
			return nil, types.ErrOrderNotFound
		}
		orderB, found := k.GetOrder(ctx, market.Id, trade.OrderBId)
		if !found {
			return nil, types.ErrOrderNotFound
		}

		if err := k.expireOrderIfNeeded(ctx, orderA); err != nil {
			return nil, err
		}
		if err := k.expireOrderIfNeeded(ctx, orderB); err != nil {
			return nil, err
		}

		if err := validateOrderMatchable(orderA, matchAmount, executionPrice); err != nil {
			return nil, fmt.Errorf("trade[%d] order_a: %w", i, err)
		}
		if err := validateOrderMatchable(orderB, matchAmount, executionPrice); err != nil {
			return nil, fmt.Errorf("trade[%d] order_b: %w", i, err)
		}

		feeUnit := executionPrice.MulRaw(int64(market.FeeBps)).QuoRaw(int64(types.BPSDenominator))
		moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)

		switch {
		case isBuyYesBuyNo(orderA.Side, orderB.Side):
			var yesOrder, noOrder *types.Order
			if orderA.Side == types.OrderSide_ORDER_SIDE_BUY_YES {
				yesOrder, noOrder = orderA, orderB
			} else {
				yesOrder, noOrder = orderB, orderA
			}

			yesBuyer := sdk.MustAccAddressFromBech32(yesOrder.Trader)
			noBuyer := sdk.MustAccAddressFromBech32(noOrder.Trader)
			if err := k.transferCollateralBetweenAccounts(ctx, market, yesBuyer, moduleAddr, executionPrice); err != nil {
				return nil, err
			}
			if err := k.transferCollateralBetweenAccounts(ctx, market, noBuyer, moduleAddr, executionPrice); err != nil {
				return nil, err
			}
			// BUY_YES <-> BUY_NO fee is deducted from collected trade collateral,
			// not charged as extra on top of buyer payment.
			if err := k.chargeAndDistributeFeeFromModule(ctx, market, feeUnit); err != nil {
				return nil, err
			}
			if err := k.chargeAndDistributeFeeFromModule(ctx, market, feeUnit); err != nil {
				return nil, err
			}

			yesPos := k.getPositionOrDefault(ctx, market.Id, yesBuyer)
			yesShares, yesLocked, yesNoShares, yesNoLocked, err := k.mustPositionInts(yesPos)
			if err != nil {
				return nil, err
			}
			yesShares = yesShares.Add(matchAmount)
			k.mustSetPositionInts(yesPos, yesShares, yesLocked, yesNoShares, yesNoLocked)
			if err := k.assertPositionInvariant(yesPos); err != nil {
				return nil, err
			}
			k.SetPosition(ctx, yesPos)

			noPos := k.getPositionOrDefault(ctx, market.Id, noBuyer)
			noYesShares, noYesLocked, noShares, noLocked, err := k.mustPositionInts(noPos)
			if err != nil {
				return nil, err
			}
			noShares = noShares.Add(matchAmount)
			k.mustSetPositionInts(noPos, noYesShares, noYesLocked, noShares, noLocked)
			if err := k.assertPositionInvariant(noPos); err != nil {
				return nil, err
			}
			k.SetPosition(ctx, noPos)

			yesTotal = yesTotal.Add(matchAmount)
			noTotal = noTotal.Add(matchAmount)
			totalFee = totalFee.Add(feeUnit).Add(feeUnit)

		case isBuySellSameOutcome(orderA.Side, orderB.Side):
			var buyOrder, sellOrder *types.Order
			if isBuySide(orderA.Side) {
				buyOrder, sellOrder = orderA, orderB
			} else {
				buyOrder, sellOrder = orderB, orderA
			}

			buyerAddr := sdk.MustAccAddressFromBech32(buyOrder.Trader)
			sellerAddr := sdk.MustAccAddressFromBech32(sellOrder.Trader)

			sellerPos := k.getPositionOrDefault(ctx, market.Id, sellerAddr)
			sellerYes, sellerLockedYes, sellerNo, sellerLockedNo, err := k.mustPositionInts(sellerPos)
			if err != nil {
				return nil, err
			}
			buyerPos := k.getPositionOrDefault(ctx, market.Id, buyerAddr)
			buyerYes, buyerLockedYes, buyerNo, buyerLockedNo, err := k.mustPositionInts(buyerPos)
			if err != nil {
				return nil, err
			}

			netToSeller := executionPrice.Sub(feeUnit)
			if netToSeller.IsNegative() {
				return nil, fmt.Errorf("trade[%d]: fee exceeds execution_price", i)
			}
			if netToSeller.IsPositive() {
				if err := k.transferCollateralBetweenAccounts(ctx, market, buyerAddr, sellerAddr, netToSeller); err != nil {
					return nil, err
				}
			}
			if err := k.chargeAndDistributeFee(ctx, market, buyerAddr, feeUnit); err != nil {
				return nil, err
			}

			switch {
			case buyOrder.Side == types.OrderSide_ORDER_SIDE_BUY_YES && sellOrder.Side == types.OrderSide_ORDER_SIDE_SELL_YES:
				if sellerYes.LT(matchAmount) {
					return nil, errors.Wrap(types.ErrInsufficientFunds, "seller YES shares")
				}
				sellerYes = sellerYes.Sub(matchAmount)
				buyerYes = buyerYes.Add(matchAmount)
			case buyOrder.Side == types.OrderSide_ORDER_SIDE_BUY_NO && sellOrder.Side == types.OrderSide_ORDER_SIDE_SELL_NO:
				if sellerNo.LT(matchAmount) {
					return nil, errors.Wrap(types.ErrInsufficientFunds, "seller NO shares")
				}
				sellerNo = sellerNo.Sub(matchAmount)
				buyerNo = buyerNo.Add(matchAmount)
			default:
				return nil, types.ErrInvalidOrderPair
			}

			k.mustSetPositionInts(sellerPos, sellerYes, sellerLockedYes, sellerNo, sellerLockedNo)
			if err := k.assertPositionInvariant(sellerPos); err != nil {
				return nil, err
			}
			k.SetPosition(ctx, sellerPos)

			k.mustSetPositionInts(buyerPos, buyerYes, buyerLockedYes, buyerNo, buyerLockedNo)
			if err := k.assertPositionInvariant(buyerPos); err != nil {
				return nil, err
			}
			k.SetPosition(ctx, buyerPos)

			totalFee = totalFee.Add(feeUnit)

		default:
			return nil, types.ErrInvalidOrderPair
		}

		prevStatusA := orderA.Status
		prevStatusB := orderB.Status
		if err := k.fillOrder(orderA, matchAmount); err != nil {
			return nil, err
		}
		if err := k.fillOrder(orderB, matchAmount); err != nil {
			return nil, err
		}
		if err := k.onOrderStatusTransition(ctx, orderA, prevStatusA); err != nil {
			return nil, err
		}
		if err := k.onOrderStatusTransition(ctx, orderB, prevStatusB); err != nil {
			return nil, err
		}
		k.SetOrder(ctx, orderA)
		k.SetOrder(ctx, orderB)
		k.SetAppliedTrade(ctx, market.Id, trade.TradeId)
		totalTradeVolume = totalTradeVolume.Add(matchAmount)
		lastTradePrice = executionPrice.String()

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTradeSettled,
				sdk.NewAttribute(types.AttributeKeyMarketID, intToStr(market.Id)),
				sdk.NewAttribute(types.AttributeKeyTradeID, trade.TradeId),
				sdk.NewAttribute(types.AttributeKeyOrderAID, intToStr(trade.OrderAId)),
				sdk.NewAttribute(types.AttributeKeyOrderBID, intToStr(trade.OrderBId)),
				sdk.NewAttribute(types.AttributeKeyShareAmount, matchAmount.String()),
				sdk.NewAttribute(types.AttributeKeyPrice, executionPrice.String()),
				sdk.NewAttribute(types.AttributeKeyFee, feeUnit.String()),
			),
		)
	}

	market.TotalYesShares = yesTotal.String()
	market.TotalNoShares = noTotal.String()
	market.TotalTradeVolume = totalTradeVolume.String()
	market.LastTradePrice = lastTradePrice
	k.refreshMarketBookPrices(ctx, market)
	if err := k.ValidateMarketInvariants(market); err != nil {
		return nil, err
	}
	k.SetMarket(ctx, market)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTradeBatchApplied,
			sdk.NewAttribute(types.AttributeKeyMarketID, intToStr(market.Id)),
			sdk.NewAttribute(types.AttributeKeyBatchID, msg.BatchId),
			sdk.NewAttribute(types.AttributeKeySettledTrades, intToStr(uint64(len(msg.Trades)))),
			sdk.NewAttribute(types.AttributeKeyFee, totalFee.String()),
		),
	)

	return &types.MsgApplyTradeBatchResponse{
		SettledCount: uint64(len(msg.Trades)),
		TotalFees:    totalFee.String(),
	}, nil
}

func validateOrderMatchable(order *types.Order, matchAmount sdkmath.Int, executionPrice sdkmath.Int) error {
	if order.Status != types.OrderStatus_ORDER_STATUS_OPEN && order.Status != types.OrderStatus_ORDER_STATUS_PARTIALLY_FILLED {
		return types.ErrInvalidOrderStatus
	}
	remaining, err := parseNonNegativeInt(order.RemainingAmount, "remaining_amount")
	if err != nil {
		return err
	}
	if remaining.LT(matchAmount) {
		return fmt.Errorf("order remaining_amount is insufficient")
	}

	switch order.OrderType {
	case types.OrderType_ORDER_TYPE_LIMIT:
		limitPrice, err := types.ParsePriceTicks(order.LimitPrice, "limit_price")
		if err != nil {
			return err
		}
		if isBuySide(order.Side) && executionPrice.GT(limitPrice) {
			return fmt.Errorf("execution_price exceeds buy limit_price")
		}
		if isSellSide(order.Side) && executionPrice.LT(limitPrice) {
			return fmt.Errorf("execution_price below sell limit_price")
		}
	case types.OrderType_ORDER_TYPE_MARKET:
		worstPrice, err := types.ParsePriceTicks(order.WorstPrice, "worst_price")
		if err != nil {
			return err
		}
		if isBuySide(order.Side) && executionPrice.GT(worstPrice) {
			return fmt.Errorf("execution_price exceeds buy worst_price")
		}
		if isSellSide(order.Side) && executionPrice.LT(worstPrice) {
			return fmt.Errorf("execution_price below sell worst_price")
		}
	default:
		return fmt.Errorf("invalid order_type")
	}

	return nil
}

func isBuyYesBuyNo(a, b types.OrderSide) bool {
	return (a == types.OrderSide_ORDER_SIDE_BUY_YES && b == types.OrderSide_ORDER_SIDE_BUY_NO) ||
		(a == types.OrderSide_ORDER_SIDE_BUY_NO && b == types.OrderSide_ORDER_SIDE_BUY_YES)
}

func isBuySellSameOutcome(a, b types.OrderSide) bool {
	return (a == types.OrderSide_ORDER_SIDE_BUY_YES && b == types.OrderSide_ORDER_SIDE_SELL_YES) ||
		(a == types.OrderSide_ORDER_SIDE_SELL_YES && b == types.OrderSide_ORDER_SIDE_BUY_YES) ||
		(a == types.OrderSide_ORDER_SIDE_BUY_NO && b == types.OrderSide_ORDER_SIDE_SELL_NO) ||
		(a == types.OrderSide_ORDER_SIDE_SELL_NO && b == types.OrderSide_ORDER_SIDE_BUY_NO)
}

func (k Keeper) chargeAndDistributeFee(ctx sdk.Context, market *types.Market, payer sdk.AccAddress, fee sdkmath.Int) error {
	if !fee.IsPositive() {
		return nil
	}

	resolverAddr := sdk.MustAccAddressFromBech32(market.Resolver)
	creatorAddr := sdk.MustAccAddressFromBech32(market.Creator)
	resolverFee := fee.MulRaw(int64(market.ResolverFeeSharePercent)).QuoRaw(100)
	creatorFee := fee.Sub(resolverFee)

	switch {
	case resolverAddr.Equals(creatorAddr):
		if err := k.transferCollateralBetweenAccounts(ctx, market, payer, resolverAddr, fee); err != nil {
			return err
		}
	default:
		if resolverFee.IsPositive() {
			if err := k.transferCollateralBetweenAccounts(ctx, market, payer, resolverAddr, resolverFee); err != nil {
				return err
			}
		}
		if creatorFee.IsPositive() {
			if err := k.transferCollateralBetweenAccounts(ctx, market, payer, creatorAddr, creatorFee); err != nil {
				return err
			}
		}
	}

	return nil
}

func (k Keeper) chargeAndDistributeFeeFromModule(ctx sdk.Context, market *types.Market, fee sdkmath.Int) error {
	if !fee.IsPositive() {
		return nil
	}

	resolverAddr := sdk.MustAccAddressFromBech32(market.Resolver)
	creatorAddr := sdk.MustAccAddressFromBech32(market.Creator)
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	resolverFee := fee.MulRaw(int64(market.ResolverFeeSharePercent)).QuoRaw(100)
	creatorFee := fee.Sub(resolverFee)

	switch {
	case resolverAddr.Equals(creatorAddr):
		return k.transferFeeFromModule(ctx, market, moduleAddr, resolverAddr, fee)
	default:
		if resolverFee.IsPositive() {
			if err := k.transferFeeFromModule(ctx, market, moduleAddr, resolverAddr, resolverFee); err != nil {
				return err
			}
		}
		if creatorFee.IsPositive() {
			if err := k.transferFeeFromModule(ctx, market, moduleAddr, creatorAddr, creatorFee); err != nil {
				return err
			}
		}
	}

	return nil
}

func (k Keeper) transferFeeFromModule(
	ctx sdk.Context,
	market *types.Market,
	moduleAddr sdk.AccAddress,
	to sdk.AccAddress,
	amount sdkmath.Int,
) error {
	if !amount.IsPositive() {
		return nil
	}

	switch market.CollateralType {
	case types.CollateralType_COLLATERAL_TYPE_NATIVE:
		coin := sdk.NewCoin(market.CollateralDenom, amount)
		return k.bankKeeper.SendCoins(ctx, moduleAddr, to, sdk.NewCoins(coin))
	case types.CollateralType_COLLATERAL_TYPE_PRC20:
		return k.transferPRC20FromModule(ctx, to, market.CollateralContractAddr, amount)
	default:
		return fmt.Errorf("unsupported collateral type")
	}
}

func (k Keeper) HasAppliedTrade(ctx sdk.Context, marketID uint64, tradeID string) bool {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.AppliedTradeStoreKey(marketID, tradeID))
	return err == nil && bz != nil
}

func (k Keeper) SetAppliedTrade(ctx sdk.Context, marketID uint64, tradeID string) {
	store := k.storeService.OpenKVStore(ctx)
	k.mustSet(store, types.AppliedTradeStoreKey(marketID, tradeID), []byte{1})
}
