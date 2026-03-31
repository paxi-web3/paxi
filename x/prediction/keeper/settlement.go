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
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	collateralUnit := sdkmath.NewInt(types.CollateralUnit)
	appliedTradeStore := k.storeService.OpenKVStore(ctx)
	orderCache := make(map[uint64]*types.Order)
	dirtyOrders := make(map[uint64]struct{})
	deletedOrders := make(map[uint64]struct{})
	positionCache := make(map[string]*types.Position)
	dirtyPositions := make(map[string]struct{})
	newAppliedTrades := make(map[string]struct{})
	appliedTradeExistsCache := make(map[string]bool)

	type collateralAvailability struct {
		loaded         bool
		balance        sdkmath.Int
		allowance      sdkmath.Int
		balanceDelta   sdkmath.Int
		allowanceDelta sdkmath.Int
	}

	collateralCache := make(map[string]*collateralAvailability)

	getOrderCached := func(orderID uint64) (*types.Order, bool) {
		if _, deleted := deletedOrders[orderID]; deleted {
			return nil, false
		}
		if order, ok := orderCache[orderID]; ok {
			return order, true
		}
		order, found := k.GetOrder(ctx, market.Id, orderID)
		if !found {
			return nil, false
		}
		orderCache[orderID] = order
		return order, true
	}
	getPositionCached := func(addr sdk.AccAddress) *types.Position {
		key := addr.String()
		if pos, ok := positionCache[key]; ok {
			return pos
		}
		pos := k.getPositionOrDefault(ctx, market.Id, addr)
		positionCache[key] = pos
		return pos
	}
	markOrderDirty := func(order *types.Order) {
		orderCache[order.Id] = order
		dirtyOrders[order.Id] = struct{}{}
		delete(deletedOrders, order.Id)
	}
	markOrderDeleted := func(orderID uint64) {
		delete(orderCache, orderID)
		delete(dirtyOrders, orderID)
		deletedOrders[orderID] = struct{}{}
	}
	markPositionDirty := func(pos *types.Position) {
		positionCache[pos.Address] = pos
		dirtyPositions[pos.Address] = struct{}{}
	}
	cancelOrderDueToInsufficientCached := func(order *types.Order, reason string) error {
		filled, err := parseNonNegativeInt(order.FilledAmount, "filled_amount")
		if err != nil {
			return err
		}
		if err := k.cancelOrderDueToInsufficient(ctx, order, reason); err != nil {
			return err
		}
		if filled.IsZero() {
			markOrderDeleted(order.Id)
			return nil
		}
		orderCache[order.Id] = order
		delete(dirtyOrders, order.Id)
		return nil
	}
	getCollateralState := func(owner sdk.AccAddress) *collateralAvailability {
		key := owner.String()
		if state, ok := collateralCache[key]; ok {
			return state
		}
		state := &collateralAvailability{
			balance:        sdkmath.ZeroInt(),
			allowance:      sdkmath.ZeroInt(),
			balanceDelta:   sdkmath.ZeroInt(),
			allowanceDelta: sdkmath.ZeroInt(),
		}
		collateralCache[key] = state
		return state
	}
	loadCollateralState := func(owner sdk.AccAddress) (*collateralAvailability, error) {
		state := getCollateralState(owner)
		if state.loaded {
			return state, nil
		}
		switch market.CollateralType {
		case types.CollateralType_COLLATERAL_TYPE_NATIVE:
			state.balance = k.bankKeeper.GetBalance(ctx, owner, market.CollateralDenom).Amount.Add(state.balanceDelta)
		case types.CollateralType_COLLATERAL_TYPE_PRC20:
			balance, err := k.getPRC20Balance(ctx, market.CollateralContractAddr, owner)
			if err != nil {
				return nil, err
			}
			allowance, err := k.getPRC20Allowance(ctx, market.CollateralContractAddr, owner)
			if err != nil {
				return nil, err
			}
			state.balance = balance.Add(state.balanceDelta)
			state.allowance = allowance.Add(state.allowanceDelta)
		default:
			return nil, fmt.Errorf("unsupported collateral type")
		}
		state.loaded = true
		return state, nil
	}
	ensureCollateralBalanceCached := func(owner sdk.AccAddress, required sdkmath.Int) error {
		if !required.IsPositive() {
			return nil
		}
		state, err := loadCollateralState(owner)
		if err != nil {
			return err
		}
		if state.balance.LT(required) {
			switch market.CollateralType {
			case types.CollateralType_COLLATERAL_TYPE_NATIVE:
				return fmt.Errorf("insufficient funds: required=%s balance=%s", required.String(), state.balance.String())
			case types.CollateralType_COLLATERAL_TYPE_PRC20:
				return fmt.Errorf("insufficient prc20 balance: required=%s balance=%s", required.String(), state.balance.String())
			default:
				return fmt.Errorf("unsupported collateral type")
			}
		}
		if market.CollateralType == types.CollateralType_COLLATERAL_TYPE_PRC20 && state.allowance.LT(required) {
			return fmt.Errorf("insufficient prc20 allowance: required=%s allowance=%s", required.String(), state.allowance.String())
		}
		return nil
	}
	noteCollateralSent := func(owner sdk.AccAddress, amount sdkmath.Int) {
		if !amount.IsPositive() {
			return
		}
		state := getCollateralState(owner)
		if state.loaded {
			state.balance = state.balance.Sub(amount)
			if market.CollateralType == types.CollateralType_COLLATERAL_TYPE_PRC20 {
				state.allowance = state.allowance.Sub(amount)
			}
			return
		}
		state.balanceDelta = state.balanceDelta.Sub(amount)
		if market.CollateralType == types.CollateralType_COLLATERAL_TYPE_PRC20 {
			state.allowanceDelta = state.allowanceDelta.Sub(amount)
		}
	}
	noteCollateralReceived := func(owner sdk.AccAddress, amount sdkmath.Int) {
		if !amount.IsPositive() {
			return
		}
		state := getCollateralState(owner)
		if state.loaded {
			state.balance = state.balance.Add(amount)
			return
		}
		state.balanceDelta = state.balanceDelta.Add(amount)
	}
	hasAppliedTradeCached := func(tradeID string) bool {
		if exists, ok := appliedTradeExistsCache[tradeID]; ok {
			return exists
		}
		bz, err := appliedTradeStore.Get(types.AppliedTradeStoreKey(market.Id, tradeID))
		exists := err == nil && bz != nil
		appliedTradeExistsCache[tradeID] = exists
		return exists
	}

	for i := range msg.Trades {
		trade := msg.Trades[i]
		if _, exists := newAppliedTrades[trade.TradeId]; exists || hasAppliedTradeCached(trade.TradeId) {
			k.emitTradeSkipped(ctx, market.Id, trade.TradeId, "duplicate_trade_id")
			continue
		}

		matchAmount, err := parsePositiveInt(trade.MatchAmount, "trade.match_amount")
		if err != nil {
			k.emitTradeSkipped(ctx, market.Id, trade.TradeId, fmt.Sprintf("invalid_match_amount:%v", err))
			continue
		}
		yesExecutionPrice, noExecutionPrice, err := parseTradeDualExecutionPrices(trade)
		if err != nil {
			k.emitTradeSkipped(ctx, market.Id, trade.TradeId, fmt.Sprintf("invalid_execution_price:%v", err))
			continue
		}

		orderA, found := getOrderCached(trade.OrderAId)
		if !found {
			k.emitTradeSkipped(ctx, market.Id, trade.TradeId, "order_a_not_found")
			continue
		}
		orderB, found := getOrderCached(trade.OrderBId)
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

		skipTrade := false
		hasYesViewPrice := false
		yesViewPrice := sdkmath.ZeroInt()
		eventPrice := sdkmath.ZeroInt()
		hasDualPrice := false
		settledYesPrice := sdkmath.ZeroInt()
		settledNoPrice := sdkmath.ZeroInt()
		feeTotal := sdkmath.ZeroInt()
		tradeVolumeContribution := sdkmath.ZeroInt()

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
			effectiveYesExecutionPrice, effectiveNoExecutionPrice, err := normalizeBuyYesBuyNoExecutionPrices(
				yesOrder,
				noOrder,
				yesExecutionPrice,
				noExecutionPrice,
			)
			if err != nil {
				k.emitTradeSkipped(ctx, market.Id, trade.TradeId, fmt.Sprintf("buy_yes_buy_no_invalid_execution_price:%v", err))
				continue
			}

			priceForOrderA := effectiveYesExecutionPrice
			if orderA.Side == types.OrderSide_ORDER_SIDE_BUY_NO {
				priceForOrderA = effectiveNoExecutionPrice
			}
			priceForOrderB := effectiveYesExecutionPrice
			if orderB.Side == types.OrderSide_ORDER_SIDE_BUY_NO {
				priceForOrderB = effectiveNoExecutionPrice
			}
			if !effectiveYesExecutionPrice.Add(effectiveNoExecutionPrice).Equal(collateralUnit) {
				k.emitTradeSkipped(ctx, market.Id, trade.TradeId, "buy_yes_buy_no_price_sum_not_one")
				continue
			}
			if err := validateOrderMatchable(orderA, matchAmount, priceForOrderA); err != nil {
				k.emitTradeSkipped(ctx, market.Id, trade.TradeId, fmt.Sprintf("order_a_not_matchable:%v", err))
				continue
			}
			if err := validateOrderMatchable(orderB, matchAmount, priceForOrderB); err != nil {
				k.emitTradeSkipped(ctx, market.Id, trade.TradeId, fmt.Sprintf("order_b_not_matchable:%v", err))
				continue
			}

			yesBuyer := sdk.MustAccAddressFromBech32(yesOrder.Trader)
			noBuyer := sdk.MustAccAddressFromBech32(noOrder.Trader)
			yesTradeNotional := effectiveYesExecutionPrice.Mul(matchAmount)
			noTradeNotional := effectiveNoExecutionPrice.Mul(matchAmount)
			yesFeeTotal := yesTradeNotional.MulRaw(int64(market.FeeBps)).QuoRaw(int64(types.BPSDenominator))
			noFeeTotal := noTradeNotional.MulRaw(int64(market.FeeBps)).QuoRaw(int64(types.BPSDenominator))
			// For BUY_YES <-> BUY_NO, settlement fee is charged on top of notional
			// from each buyer so module collateral remains fully backed.
			yesBuyerRequired := yesTradeNotional.Add(yesFeeTotal)
			noBuyerRequired := noTradeNotional.Add(noFeeTotal)
			if err := ensureCollateralBalanceCached(yesBuyer, yesBuyerRequired); err != nil {
				if err := cancelOrderDueToInsufficientCached(yesOrder, err.Error()); err != nil {
					return nil, err
				}
				skipTrade = true
			}
			if err := ensureCollateralBalanceCached(noBuyer, noBuyerRequired); err != nil {
				if err := cancelOrderDueToInsufficientCached(noOrder, err.Error()); err != nil {
					return nil, err
				}
				skipTrade = true
			}
			if skipTrade {
				k.emitTradeSkipped(ctx, market.Id, trade.TradeId, "insufficient_balance")
				continue
			}
			if err := k.transferCollateralBetweenAccounts(ctx, market, yesBuyer, moduleAddr, yesBuyerRequired); err != nil {
				if err := cancelOrderDueToInsufficientCached(yesOrder, err.Error()); err != nil {
					return nil, err
				}
				k.emitTradeSkipped(ctx, market.Id, trade.TradeId, "transfer_from_yes_buyer_failed")
				continue
			}
			noteCollateralSent(yesBuyer, yesBuyerRequired)
			if err := k.transferCollateralBetweenAccounts(ctx, market, noBuyer, moduleAddr, noBuyerRequired); err != nil {
				if err := cancelOrderDueToInsufficientCached(noOrder, err.Error()); err != nil {
					return nil, err
				}
				k.emitTradeSkipped(ctx, market.Id, trade.TradeId, "transfer_from_no_buyer_failed")
				continue
			}
			noteCollateralSent(noBuyer, noBuyerRequired)

			yesPos := getPositionCached(yesBuyer)
			yesShares, yesLocked, yesNoShares, yesNoLocked, err := k.mustPositionInts(yesPos)
			if err != nil {
				return nil, err
			}
			yesShares = yesShares.Add(matchAmount)
			k.mustSetPositionInts(yesPos, yesShares, yesLocked, yesNoShares, yesNoLocked)
			if err := k.assertPositionInvariant(yesPos); err != nil {
				return nil, err
			}
			markPositionDirty(yesPos)

			noPos := getPositionCached(noBuyer)
			noYesShares, noYesLocked, noShares, noLocked, err := k.mustPositionInts(noPos)
			if err != nil {
				return nil, err
			}
			noShares = noShares.Add(matchAmount)
			k.mustSetPositionInts(noPos, noYesShares, noYesLocked, noShares, noLocked)
			if err := k.assertPositionInvariant(noPos); err != nil {
				return nil, err
			}
			markPositionDirty(noPos)

			yesTotal = yesTotal.Add(matchAmount)
			noTotal = noTotal.Add(matchAmount)
			if err := addOrderSpentCollateral(yesOrder, yesBuyerRequired); err != nil {
				return nil, err
			}
			if err := addOrderSpentCollateral(noOrder, noBuyerRequired); err != nil {
				return nil, err
			}
			feeTotal = yesFeeTotal.Add(noFeeTotal)
			totalFee = totalFee.Add(feeTotal)
			tradeVolumeContribution = yesTradeNotional.Add(noTradeNotional)
			hasYesViewPrice = true
			yesViewPrice = effectiveYesExecutionPrice
			eventPrice = effectiveYesExecutionPrice
			hasDualPrice = true
			settledYesPrice = effectiveYesExecutionPrice
			settledNoPrice = effectiveNoExecutionPrice

		case isBuySellSameOutcome(orderA.Side, orderB.Side):
			executionPrice := sdkmath.ZeroInt()
			switch {
			case (orderA.Side == types.OrderSide_ORDER_SIDE_BUY_YES && orderB.Side == types.OrderSide_ORDER_SIDE_SELL_YES) ||
				(orderA.Side == types.OrderSide_ORDER_SIDE_SELL_YES && orderB.Side == types.OrderSide_ORDER_SIDE_BUY_YES):
				executionPrice = yesExecutionPrice
			case (orderA.Side == types.OrderSide_ORDER_SIDE_BUY_NO && orderB.Side == types.OrderSide_ORDER_SIDE_SELL_NO) ||
				(orderA.Side == types.OrderSide_ORDER_SIDE_SELL_NO && orderB.Side == types.OrderSide_ORDER_SIDE_BUY_NO):
				executionPrice = noExecutionPrice
			default:
				k.emitTradeSkipped(ctx, market.Id, trade.TradeId, "invalid_order_pair")
				continue
			}
			executionPrice, err = normalizeLimitMarketExecutionPrice(orderA, orderB, executionPrice)
			if err != nil {
				k.emitTradeSkipped(ctx, market.Id, trade.TradeId, fmt.Sprintf("buy_sell_same_outcome_invalid_execution_price:%v", err))
				continue
			}
			if err := validateOrderMatchable(orderA, matchAmount, executionPrice); err != nil {
				k.emitTradeSkipped(ctx, market.Id, trade.TradeId, fmt.Sprintf("order_a_not_matchable:%v", err))
				continue
			}
			if err := validateOrderMatchable(orderB, matchAmount, executionPrice); err != nil {
				k.emitTradeSkipped(ctx, market.Id, trade.TradeId, fmt.Sprintf("order_b_not_matchable:%v", err))
				continue
			}

			var buyOrder, sellOrder *types.Order
			if isBuySide(orderA.Side) {
				buyOrder, sellOrder = orderA, orderB
			} else {
				buyOrder, sellOrder = orderB, orderA
			}

			buyerAddr := sdk.MustAccAddressFromBech32(buyOrder.Trader)
			sellerAddr := sdk.MustAccAddressFromBech32(sellOrder.Trader)
			isSelfTrade := buyerAddr.Equals(sellerAddr)

			sellerPos := getPositionCached(sellerAddr)
			sellerYes, sellerLockedYes, sellerNo, sellerLockedNo, err := k.mustPositionInts(sellerPos)
			if err != nil {
				return nil, err
			}
			buyerPos := sellerPos
			buyerYes, buyerLockedYes, buyerNo, buyerLockedNo := sellerYes, sellerLockedYes, sellerNo, sellerLockedNo
			if !isSelfTrade {
				buyerPos = getPositionCached(buyerAddr)
				buyerYes, buyerLockedYes, buyerNo, buyerLockedNo, err = k.mustPositionInts(buyerPos)
				if err != nil {
					return nil, err
				}
			}

			tradeNotional := executionPrice.Mul(matchAmount)
			feeTotal = tradeNotional.MulRaw(int64(market.FeeBps)).QuoRaw(int64(types.BPSDenominator))
			netToSeller := tradeNotional.Sub(feeTotal)
			if netToSeller.IsNegative() {
				k.emitTradeSkipped(ctx, market.Id, trade.TradeId, "fee_exceeds_execution_price")
				continue
			}
			// Buyer actual debit equals netToSeller + feeTotal = tradeNotional.
			buyerRequired := tradeNotional
			if err := ensureCollateralBalanceCached(buyerAddr, buyerRequired); err != nil {
				if err := cancelOrderDueToInsufficientCached(buyOrder, err.Error()); err != nil {
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
					if err := cancelOrderDueToInsufficientCached(sellOrder, "seller YES shares"); err != nil {
						return nil, err
					}
					skipTrade = true
				} else {
					if !isSelfTrade {
						sellerYes = sellerYes.Sub(matchAmount)
						buyerYes = buyerYes.Add(matchAmount)
					}
				}
				hasYesViewPrice = true
				yesViewPrice = executionPrice
			case buyOrder.Side == types.OrderSide_ORDER_SIDE_BUY_NO && sellOrder.Side == types.OrderSide_ORDER_SIDE_SELL_NO:
				freeNo := sellerNo.Sub(sellerLockedNo)
				if freeNo.LT(matchAmount) {
					if err := cancelOrderDueToInsufficientCached(sellOrder, "seller NO shares"); err != nil {
						return nil, err
					}
					skipTrade = true
				} else {
					if !isSelfTrade {
						sellerNo = sellerNo.Sub(matchAmount)
						buyerNo = buyerNo.Add(matchAmount)
					}
				}
				hasYesViewPrice = true
				yesViewPrice = collateralUnit.Sub(executionPrice)
			default:
				k.emitTradeSkipped(ctx, market.Id, trade.TradeId, "invalid_order_pair")
				continue
			}
			if skipTrade {
				k.emitTradeSkipped(ctx, market.Id, trade.TradeId, "insufficient_balance_or_shares")
				continue
			}
			if netToSeller.IsPositive() && !isSelfTrade {
				if err := k.transferCollateralBetweenAccounts(ctx, market, buyerAddr, sellerAddr, netToSeller); err != nil {
					if err := cancelOrderDueToInsufficientCached(buyOrder, err.Error()); err != nil {
						return nil, err
					}
					k.emitTradeSkipped(ctx, market.Id, trade.TradeId, "buyer_to_seller_transfer_failed")
					continue
				}
				noteCollateralSent(buyerAddr, netToSeller)
				noteCollateralReceived(sellerAddr, netToSeller)
			}
			if err := k.transferCollateralBetweenAccounts(ctx, market, buyerAddr, moduleAddr, feeTotal); err != nil {
				if err := cancelOrderDueToInsufficientCached(buyOrder, err.Error()); err != nil {
					return nil, err
				}
				k.emitTradeSkipped(ctx, market.Id, trade.TradeId, "fee_charge_failed")
				continue
			}
			noteCollateralSent(buyerAddr, feeTotal)

			if isSelfTrade {
				k.mustSetPositionInts(sellerPos, sellerYes, sellerLockedYes, sellerNo, sellerLockedNo)
				if err := k.assertPositionInvariant(sellerPos); err != nil {
					return nil, err
				}
				markPositionDirty(sellerPos)
			} else {
				k.mustSetPositionInts(sellerPos, sellerYes, sellerLockedYes, sellerNo, sellerLockedNo)
				if err := k.assertPositionInvariant(sellerPos); err != nil {
					return nil, err
				}
				markPositionDirty(sellerPos)

				k.mustSetPositionInts(buyerPos, buyerYes, buyerLockedYes, buyerNo, buyerLockedNo)
				if err := k.assertPositionInvariant(buyerPos); err != nil {
					return nil, err
				}
				markPositionDirty(buyerPos)
			}

			if err := addOrderSpentCollateral(buyOrder, tradeNotional); err != nil {
				return nil, err
			}
			if err := addOrderReceivedCollateral(sellOrder, netToSeller); err != nil {
				return nil, err
			}

			totalFee = totalFee.Add(feeTotal)
			tradeVolumeContribution = tradeNotional
			eventPrice = executionPrice

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
		markOrderDirty(orderA)
		markOrderDirty(orderB)
		newAppliedTrades[trade.TradeId] = struct{}{}
		totalTradeVolume = totalTradeVolume.Add(tradeVolumeContribution)
		if hasYesViewPrice {
			weightedYesPriceAmount = weightedYesPriceAmount.Add(yesViewPrice.Mul(matchAmount))
			weightedPriceAmountDenominator = weightedPriceAmountDenominator.Add(matchAmount)
		}
		settledCount++

		attrs := []sdk.Attribute{
			sdk.NewAttribute(types.AttributeKeyMarketID, intToStr(market.Id)),
			sdk.NewAttribute(types.AttributeKeyTradeID, trade.TradeId),
			sdk.NewAttribute(types.AttributeKeyOrderAID, intToStr(trade.OrderAId)),
			sdk.NewAttribute(types.AttributeKeyOrderBID, intToStr(trade.OrderBId)),
			sdk.NewAttribute(types.AttributeKeyShareAmount, matchAmount.String()),
			sdk.NewAttribute(types.AttributeKeyPrice, eventPrice.String()),
			sdk.NewAttribute(types.AttributeKeyFee, feeTotal.String()),
		}
		attrs = append(attrs, tradeSettledDualPriceAttrs(hasDualPrice, settledYesPrice, settledNoPrice)...)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTradeSettled,
				attrs...,
			),
		)
	}

	for orderID := range dirtyOrders {
		if _, deleted := deletedOrders[orderID]; deleted {
			continue
		}
		k.SetOrder(ctx, orderCache[orderID])
	}
	for addr := range dirtyPositions {
		k.SetPosition(ctx, positionCache[addr])
	}
	for tradeID := range newAppliedTrades {
		k.SetAppliedTrade(ctx, market.Id, tradeID)
	}

	market.TotalYesShares = yesTotal.String()
	market.TotalNoShares = noTotal.String()
	market.TotalTradeVolume = totalTradeVolume.String()
	if err := k.chargeAndDistributeFeeFromModule(ctx, market, totalFee); err != nil {
		return nil, err
	}
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
			sdk.NewAttribute(types.AttributeKeyTotalYesShares, market.TotalYesShares),
			sdk.NewAttribute(types.AttributeKeyTotalNoShares, market.TotalNoShares),
			sdk.NewAttribute(types.AttributeKeyTotalTradeVolume, market.TotalTradeVolume),
			sdk.NewAttribute(types.AttributeKeyLastYesTradePrice, market.LastYesTradePrice),
			sdk.NewAttribute(types.AttributeKeyLastNoTradePrice, market.LastNoTradePrice),
		),
	)

	return &types.MsgApplyTradeBatchResponse{
		SettledCount: settledCount,
		TotalFees:    totalFee.String(),
	}, nil
}

func parseTradeDualExecutionPrices(trade types.TradeMatch) (sdkmath.Int, sdkmath.Int, error) {
	yesPrice, err := types.ParsePriceTicks(trade.YesExecutionPrice, "trade.yes_execution_price")
	if err != nil {
		return sdkmath.Int{}, sdkmath.Int{}, err
	}
	noPrice, err := types.ParsePriceTicks(trade.NoExecutionPrice, "trade.no_execution_price")
	if err != nil {
		return sdkmath.Int{}, sdkmath.Int{}, err
	}
	return yesPrice, noPrice, nil
}

func tradeSettledDualPriceAttrs(enabled bool, yesPrice sdkmath.Int, noPrice sdkmath.Int) []sdk.Attribute {
	if !enabled {
		return nil
	}
	return []sdk.Attribute{
		sdk.NewAttribute(types.AttributeKeyYesPrice, yesPrice.String()),
		sdk.NewAttribute(types.AttributeKeyNoPrice, noPrice.String()),
	}
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

func addOrderSpentCollateral(order *types.Order, delta sdkmath.Int) error {
	return addOrderCollateral(&order.SpentCollateral, delta, "spent_collateral")
}

func addOrderReceivedCollateral(order *types.Order, delta sdkmath.Int) error {
	return addOrderCollateral(&order.ReceivedCollateral, delta, "received_collateral")
}

func addOrderCollateral(current *string, delta sdkmath.Int, field string) error {
	if delta.IsNegative() {
		return fmt.Errorf("%s delta cannot be negative", field)
	}
	if delta.IsZero() {
		if *current == "" {
			*current = sdkmath.ZeroInt().String()
		}
		return nil
	}

	base := sdkmath.ZeroInt()
	if *current != "" {
		parsed, err := parseNonNegativeInt(*current, field)
		if err != nil {
			return err
		}
		base = parsed
	}
	*current = base.Add(delta).String()
	return nil
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

func normalizeLimitMarketExecutionPrice(orderA *types.Order, orderB *types.Order, executionPrice sdkmath.Int) (sdkmath.Int, error) {
	orderAIsLimit := orderA.OrderType == types.OrderType_ORDER_TYPE_LIMIT
	orderBIsLimit := orderB.OrderType == types.OrderType_ORDER_TYPE_LIMIT

	switch {
	case orderAIsLimit && !orderBIsLimit:
		limitPrice, err := types.ParsePriceTicks(orderA.LimitPrice, "limit_price")
		if err != nil {
			return sdkmath.Int{}, err
		}
		return limitPrice, nil
	case orderBIsLimit && !orderAIsLimit:
		limitPrice, err := types.ParsePriceTicks(orderB.LimitPrice, "limit_price")
		if err != nil {
			return sdkmath.Int{}, err
		}
		return limitPrice, nil
	}

	return executionPrice, nil
}

func normalizeBuyYesBuyNoExecutionPrices(
	yesOrder *types.Order,
	noOrder *types.Order,
	yesExecutionPrice sdkmath.Int,
	noExecutionPrice sdkmath.Int,
) (sdkmath.Int, sdkmath.Int, error) {
	effectiveYes := yesExecutionPrice
	effectiveNo := noExecutionPrice
	unit := sdkmath.NewInt(types.CollateralUnit)
	yesIsLimit := yesOrder.OrderType == types.OrderType_ORDER_TYPE_LIMIT
	noIsLimit := noOrder.OrderType == types.OrderType_ORDER_TYPE_LIMIT

	switch {
	case yesIsLimit && !noIsLimit:
		yesLimitPrice, err := types.ParsePriceTicks(yesOrder.LimitPrice, "limit_price")
		if err != nil {
			return sdkmath.Int{}, sdkmath.Int{}, err
		}
		effectiveYes = yesLimitPrice
		effectiveNo = unit.Sub(yesLimitPrice)
	case noIsLimit && !yesIsLimit:
		noLimitPrice, err := types.ParsePriceTicks(noOrder.LimitPrice, "limit_price")
		if err != nil {
			return sdkmath.Int{}, sdkmath.Int{}, err
		}
		effectiveNo = noLimitPrice
		effectiveYes = unit.Sub(noLimitPrice)
	}

	if !effectiveYes.Add(effectiveNo).Equal(unit) {
		return sdkmath.Int{}, sdkmath.Int{}, fmt.Errorf("yes_execution_price + no_execution_price must equal collateral unit")
	}
	if err := validateExecutionPriceTicks(effectiveYes, "yes_execution_price"); err != nil {
		return sdkmath.Int{}, sdkmath.Int{}, err
	}
	if err := validateExecutionPriceTicks(effectiveNo, "no_execution_price"); err != nil {
		return sdkmath.Int{}, sdkmath.Int{}, err
	}

	return effectiveYes, effectiveNo, nil
}

func validateExecutionPriceTicks(price sdkmath.Int, field string) error {
	min := sdkmath.NewInt(types.MinPriceTicks)
	max := sdkmath.NewInt(types.MaxPriceTicks)
	if price.LT(min) || price.GT(max) {
		return fmt.Errorf("%s must be between %d and %d", field, types.MinPriceTicks, types.MaxPriceTicks)
	}
	if !price.ModRaw(types.PriceTickSize).IsZero() {
		return fmt.Errorf("%s must be a multiple of %d", field, types.PriceTickSize)
	}
	return nil
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
