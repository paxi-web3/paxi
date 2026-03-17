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
	lastYesTradePrice := market.LastYesTradePrice
	lastNoTradePrice := market.LastNoTradePrice
	weightedYesPriceAmount := sdkmath.ZeroInt()
	weightedPriceAmountDenominator := sdkmath.ZeroInt()
	settledCount := uint64(0)

	for i := range msg.Trades {
		trade := msg.Trades[i]
		if k.HasAppliedTrade(ctx, market.Id, trade.TradeId) {
			k.emitTradeSkipped(ctx, market.Id, trade.TradeId, "duplicate_trade_id")
			continue
		}

		matchAmount, err := parsePositiveInt(trade.MatchAmount, "trade.match_amount")
		if err != nil {
			k.emitTradeSkipped(ctx, market.Id, trade.TradeId, fmt.Sprintf("invalid_match_amount:%v", err))
			continue
		}
		executionPrice, err := types.ParsePriceTicks(trade.ExecutionPrice, "trade.execution_price")
		if err != nil {
			k.emitTradeSkipped(ctx, market.Id, trade.TradeId, fmt.Sprintf("invalid_execution_price:%v", err))
			continue
		}

		orderA, found := k.GetOrder(ctx, market.Id, trade.OrderAId)
		if !found {
			k.emitTradeSkipped(ctx, market.Id, trade.TradeId, "order_a_not_found")
			continue
		}
		orderB, found := k.GetOrder(ctx, market.Id, trade.OrderBId)
		if !found {
			k.emitTradeSkipped(ctx, market.Id, trade.TradeId, "order_b_not_found")
			continue
		}

		if err := k.expireOrderIfNeeded(ctx, orderA); err != nil {
			return nil, err
		}
		if err := k.expireOrderIfNeeded(ctx, orderB); err != nil {
			return nil, err
		}

		if err := validateOrderMatchable(orderA, matchAmount, executionPrice); err != nil {
			k.emitTradeSkipped(ctx, market.Id, trade.TradeId, fmt.Sprintf("order_a_not_matchable:%v", err))
			continue
		}
		if err := validateOrderMatchable(orderB, matchAmount, executionPrice); err != nil {
			k.emitTradeSkipped(ctx, market.Id, trade.TradeId, fmt.Sprintf("order_b_not_matchable:%v", err))
			continue
		}

		tradeNotional := executionPrice.Mul(matchAmount)
		feeTotal := tradeNotional.MulRaw(int64(market.FeeBps)).QuoRaw(int64(types.BPSDenominator))
		moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
		skipTrade := false
		hasYesViewPrice := false
		yesViewPrice := sdkmath.ZeroInt()

		switch {
		case isBuyYesBuyNo(orderA.Side, orderB.Side):
			var yesOrder, noOrder *types.Order
			if orderA.Side == types.OrderSide_ORDER_SIDE_BUY_YES {
				yesOrder, noOrder = orderA, orderB
			} else {
				yesOrder, noOrder = orderB, orderA
			}
			// For BUY_YES <-> BUY_NO, do not allow market-market pairing.
			// At least one side must be a LIMIT order so execution pricing has
			// an anchored quote instead of two open-ended worst prices.
			if yesOrder.OrderType == types.OrderType_ORDER_TYPE_MARKET && noOrder.OrderType == types.OrderType_ORDER_TYPE_MARKET {
				k.emitTradeSkipped(ctx, market.Id, trade.TradeId, "buy_yes_buy_no_market_market_not_allowed")
				continue
			}

			yesBuyer := sdk.MustAccAddressFromBech32(yesOrder.Trader)
			noBuyer := sdk.MustAccAddressFromBech32(noOrder.Trader)
			// For BUY_YES <-> BUY_NO, settlement fee is charged on top of notional
			// from each buyer so module collateral remains fully backed.
			buyerRequired := tradeNotional.Add(feeTotal)
			if err := k.ensureCollateralBalance(ctx, market, yesBuyer, buyerRequired); err != nil {
				if err := k.cancelOrderDueToInsufficient(ctx, yesOrder, err.Error()); err != nil {
					return nil, err
				}
				skipTrade = true
			}
			if err := k.ensureCollateralBalance(ctx, market, noBuyer, buyerRequired); err != nil {
				if err := k.cancelOrderDueToInsufficient(ctx, noOrder, err.Error()); err != nil {
					return nil, err
				}
				skipTrade = true
			}
			if skipTrade {
				k.emitTradeSkipped(ctx, market.Id, trade.TradeId, "insufficient_balance")
				continue
			}
			if err := k.transferCollateralBetweenAccounts(ctx, market, yesBuyer, moduleAddr, tradeNotional); err != nil {
				if err := k.cancelOrderDueToInsufficient(ctx, yesOrder, err.Error()); err != nil {
					return nil, err
				}
				k.emitTradeSkipped(ctx, market.Id, trade.TradeId, "transfer_from_yes_buyer_failed")
				continue
			}
			if err := k.transferCollateralBetweenAccounts(ctx, market, noBuyer, moduleAddr, tradeNotional); err != nil {
				if err := k.cancelOrderDueToInsufficient(ctx, noOrder, err.Error()); err != nil {
					return nil, err
				}
				k.emitTradeSkipped(ctx, market.Id, trade.TradeId, "transfer_from_no_buyer_failed")
				continue
			}
			// Fee is charged from each buyer on top of collateral notional.
			if err := k.chargeAndDistributeFee(ctx, market, yesBuyer, feeTotal); err != nil {
				if err := k.cancelOrderDueToInsufficient(ctx, yesOrder, err.Error()); err != nil {
					return nil, err
				}
				k.emitTradeSkipped(ctx, market.Id, trade.TradeId, "yes_buyer_fee_charge_failed")
				continue
			}
			if err := k.chargeAndDistributeFee(ctx, market, noBuyer, feeTotal); err != nil {
				if err := k.cancelOrderDueToInsufficient(ctx, noOrder, err.Error()); err != nil {
					return nil, err
				}
				k.emitTradeSkipped(ctx, market.Id, trade.TradeId, "no_buyer_fee_charge_failed")
				continue
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
			totalFee = totalFee.Add(feeTotal).Add(feeTotal)
			hasYesViewPrice = true
			yesViewPrice = executionPrice

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

			netToSeller := tradeNotional.Sub(feeTotal)
			if netToSeller.IsNegative() {
				k.emitTradeSkipped(ctx, market.Id, trade.TradeId, "fee_exceeds_execution_price")
				continue
			}
			// Buyer actual debit equals netToSeller + feeTotal = tradeNotional.
			buyerRequired := tradeNotional
			if err := k.ensureCollateralBalance(ctx, market, buyerAddr, buyerRequired); err != nil {
				if err := k.cancelOrderDueToInsufficient(ctx, buyOrder, err.Error()); err != nil {
					return nil, err
				}
				skipTrade = true
			}
			if netToSeller.IsPositive() {
				// Transfer is executed after all prechecks.
			}

			switch {
			case buyOrder.Side == types.OrderSide_ORDER_SIDE_BUY_YES && sellOrder.Side == types.OrderSide_ORDER_SIDE_SELL_YES:
				freeYes := sellerYes.Sub(sellerLockedYes)
				if freeYes.LT(matchAmount) {
					if err := k.cancelOrderDueToInsufficient(ctx, sellOrder, "seller YES shares"); err != nil {
						return nil, err
					}
					skipTrade = true
				} else {
					sellerYes = sellerYes.Sub(matchAmount)
					buyerYes = buyerYes.Add(matchAmount)
				}
				hasYesViewPrice = true
				yesViewPrice = executionPrice
			case buyOrder.Side == types.OrderSide_ORDER_SIDE_BUY_NO && sellOrder.Side == types.OrderSide_ORDER_SIDE_SELL_NO:
				freeNo := sellerNo.Sub(sellerLockedNo)
				if freeNo.LT(matchAmount) {
					if err := k.cancelOrderDueToInsufficient(ctx, sellOrder, "seller NO shares"); err != nil {
						return nil, err
					}
					skipTrade = true
				} else {
					sellerNo = sellerNo.Sub(matchAmount)
					buyerNo = buyerNo.Add(matchAmount)
				}
				hasYesViewPrice = true
				yesViewPrice = sdkmath.NewInt(types.CollateralUnit).Sub(executionPrice)
			default:
				k.emitTradeSkipped(ctx, market.Id, trade.TradeId, "invalid_order_pair")
				continue
			}
			if skipTrade {
				k.emitTradeSkipped(ctx, market.Id, trade.TradeId, "insufficient_balance_or_shares")
				continue
			}
			if netToSeller.IsPositive() {
				if err := k.transferCollateralBetweenAccounts(ctx, market, buyerAddr, sellerAddr, netToSeller); err != nil {
					if err := k.cancelOrderDueToInsufficient(ctx, buyOrder, err.Error()); err != nil {
						return nil, err
					}
					k.emitTradeSkipped(ctx, market.Id, trade.TradeId, "buyer_to_seller_transfer_failed")
					continue
				}
			}
			if err := k.chargeAndDistributeFee(ctx, market, buyerAddr, feeTotal); err != nil {
				if err := k.cancelOrderDueToInsufficient(ctx, buyOrder, err.Error()); err != nil {
					return nil, err
				}
				k.emitTradeSkipped(ctx, market.Id, trade.TradeId, "fee_charge_failed")
				continue
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

			totalFee = totalFee.Add(feeTotal)

		default:
			k.emitTradeSkipped(ctx, market.Id, trade.TradeId, "invalid_order_pair")
			continue
		}

		prevStatusA := orderA.Status
		prevStatusB := orderB.Status
		if err := k.fillOrder(orderA, matchAmount); err != nil {
			k.emitTradeSkipped(ctx, market.Id, trade.TradeId, fmt.Sprintf("fill_order_a_failed:%v", err))
			continue
		}
		if err := k.fillOrder(orderB, matchAmount); err != nil {
			k.emitTradeSkipped(ctx, market.Id, trade.TradeId, fmt.Sprintf("fill_order_b_failed:%v", err))
			continue
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
		if hasYesViewPrice {
			weightedYesPriceAmount = weightedYesPriceAmount.Add(yesViewPrice.Mul(matchAmount))
			weightedPriceAmountDenominator = weightedPriceAmountDenominator.Add(matchAmount)
		}
		settledCount++

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTradeSettled,
				sdk.NewAttribute(types.AttributeKeyMarketID, intToStr(market.Id)),
				sdk.NewAttribute(types.AttributeKeyTradeID, trade.TradeId),
				sdk.NewAttribute(types.AttributeKeyOrderAID, intToStr(trade.OrderAId)),
				sdk.NewAttribute(types.AttributeKeyOrderBID, intToStr(trade.OrderBId)),
				sdk.NewAttribute(types.AttributeKeyShareAmount, matchAmount.String()),
				sdk.NewAttribute(types.AttributeKeyPrice, executionPrice.String()),
				sdk.NewAttribute(types.AttributeKeyFee, feeTotal.String()),
			),
		)
	}

	market.TotalYesShares = yesTotal.String()
	market.TotalNoShares = noTotal.String()
	market.TotalTradeVolume = totalTradeVolume.String()
	if weightedPriceAmountDenominator.IsPositive() {
		canonicalLastYes, canonicalLastNo := canonicalOutcomeLastPrices(
			weightedYesPriceAmount,
			weightedPriceAmountDenominator,
		)
		lastYesTradePrice = canonicalLastYes.String()
		lastNoTradePrice = canonicalLastNo.String()
	}
	market.LastYesTradePrice = lastYesTradePrice
	market.LastNoTradePrice = lastNoTradePrice
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
			sdk.NewAttribute(types.AttributeKeySettledTrades, intToStr(settledCount)),
			sdk.NewAttribute(types.AttributeKeyFee, totalFee.String()),
		),
	)

	return &types.MsgApplyTradeBatchResponse{
		SettledCount: settledCount,
		TotalFees:    totalFee.String(),
	}, nil
}

func canonicalOutcomeLastPrices(weightedYesPriceAmount sdkmath.Int, totalAmount sdkmath.Int) (sdkmath.Int, sdkmath.Int) {
	unit := sdkmath.NewInt(types.CollateralUnit)
	tick := sdkmath.NewInt(types.PriceTickSize)
	denominator := totalAmount.Mul(tick)
	roundedTickCount := weightedYesPriceAmount.Add(denominator.QuoRaw(2)).Quo(denominator)
	yes := roundedTickCount.Mul(tick)
	if yes.IsNegative() {
		yes = sdkmath.ZeroInt()
	}
	if yes.GT(unit) {
		yes = unit
	}
	no := unit.Sub(yes)
	if no.IsNegative() {
		no = sdkmath.ZeroInt()
	}
	return yes, no
}

func (k Keeper) cancelOrderDueToInsufficient(ctx sdk.Context, order *types.Order, reason string) error {
	if order.Status != types.OrderStatus_ORDER_STATUS_OPEN && order.Status != types.OrderStatus_ORDER_STATUS_PARTIALLY_FILLED {
		return nil
	}

	prevStatus := order.Status
	order.Status = types.OrderStatus_ORDER_STATUS_CANCELLED
	if err := k.onOrderStatusTransition(ctx, order, prevStatus); err != nil {
		return err
	}

	filled, err := parseNonNegativeInt(order.FilledAmount, "filled_amount")
	if err != nil {
		return err
	}
	if filled.IsZero() {
		k.DeleteOrder(ctx, order.MarketId, order.Id)
	} else {
		k.SetOrder(ctx, order)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventOrderCancelled,
			sdk.NewAttribute(types.AttributeKeyMarketID, intToStr(order.MarketId)),
			sdk.NewAttribute(types.AttributeKeyOrderID, intToStr(order.Id)),
			sdk.NewAttribute(types.AttributeKeyAddress, order.Trader),
			sdk.NewAttribute("reason", reason),
		),
	)

	return nil
}

func (k Keeper) emitTradeSkipped(ctx sdk.Context, marketID uint64, tradeID string, reason string) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"TradeSkipped",
			sdk.NewAttribute(types.AttributeKeyMarketID, intToStr(marketID)),
			sdk.NewAttribute(types.AttributeKeyTradeID, tradeID),
			sdk.NewAttribute("reason", reason),
		),
	)
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
