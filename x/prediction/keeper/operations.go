package keeper

import (
	"fmt"
	"math"
	"strings"

	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/x/prediction/types"
)

func (k Keeper) CreateMarket(ctx sdk.Context, msg *types.MsgCreateMarket) (uint64, error) {
	if err := msg.ValidateBasic(); err != nil {
		return 0, errors.Wrap(types.ErrInvalidRequest, err.Error())
	}

	outcomeType := strings.ToUpper(strings.TrimSpace(msg.OutcomeType))
	if outcomeType != "BINARY" {
		return 0, errors.Wrap(types.ErrInvalidRequest, "outcome_type must be BINARY")
	}

	if len(msg.Outcomes) != 2 {
		return 0, errors.Wrap(types.ErrInvalidRequest, "binary market must have exactly two outcomes")
	}
	first := normalizeOutcome(msg.Outcomes[0])
	second := normalizeOutcome(msg.Outcomes[1])
	if !(first == "YES" && second == "NO") {
		return 0, errors.Wrap(types.ErrInvalidRequest, "outcomes must be [YES, NO]")
	}

	params := k.GetParams(ctx)
	if err := params.Validate(); err != nil {
		return 0, err
	}

	creatorAddr := sdk.MustAccAddressFromBech32(msg.Creator)
	createBond, err := parsePositiveInt(params.CreateMarketBond, "create_market_bond")
	if err != nil {
		return 0, err
	}
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		creatorAddr,
		types.ModuleName,
		sdk.NewCoins(sdk.NewCoin(params.CreateMarketBondDenom, createBond)),
	); err != nil {
		return 0, err
	}

	id := k.NextMarketID(ctx)
	market := &types.Market{
		Id:                      id,
		Creator:                 msg.Creator,
		Resolver:                msg.Resolver,
		Title:                   strings.TrimSpace(msg.Title),
		Description:             msg.Description,
		Rule:                    msg.Rule,
		ResolutionSource:        msg.ResolutionSource,
		OutcomeType:             outcomeType,
		Outcomes:                []string{"YES", "NO"},
		CollateralType:          msg.CollateralType,
		CollateralDenom:         msg.CollateralDenom,
		CollateralContractAddr:  msg.CollateralContractAddr,
		OpenTime:                msg.OpenTime,
		CloseTime:               msg.CloseTime,
		ResolveTime:             msg.ResolveTime,
		FeeBps:                  params.MarketFeeBps,
		Status:                  types.MarketStatus_MARKET_STATUS_OPEN,
		WinningOutcome:          types.Outcome_OUTCOME_UNSPECIFIED,
		TotalYesShares:          sdkmath.ZeroInt().String(),
		TotalNoShares:           sdkmath.ZeroInt().String(),
		CreateMarketBond:        createBond.String(),
		CreateMarketBondDenom:   params.CreateMarketBondDenom,
		ResolverFeeSharePercent: params.ResolverFeeSharePercent,
		LastTradePrice:          "",
		BestBidPrice:            "",
		BestAskPrice:            "",
		TotalTradeVolume:        sdkmath.ZeroInt().String(),
	}

	if err := k.ValidateMarketInvariants(market); err != nil {
		return 0, err
	}

	k.SetMarket(ctx, market)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventMarketCreated,
			sdk.NewAttribute(types.AttributeKeyMarketID, intToStr(id)),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Creator),
			sdk.NewAttribute("resolver", msg.Resolver),
			sdk.NewAttribute(types.AttributeKeyBond, createBond.String()),
			sdk.NewAttribute(types.AttributeKeyBondDenom, params.CreateMarketBondDenom),
			sdk.NewAttribute("market_fee_bps", intToStr(params.MarketFeeBps)),
			sdk.NewAttribute("resolver_fee_share_percent", intToStr(params.ResolverFeeSharePercent)),
		),
	)

	return id, nil
}

func (k Keeper) PlaceOrder(ctx sdk.Context, msg *types.MsgPlaceOrder) (uint64, error) {
	if err := msg.ValidateBasic(); err != nil {
		return 0, errors.Wrap(types.ErrInvalidRequest, err.Error())
	}

	market, found := k.GetMarket(ctx, msg.MarketId)
	if !found {
		return 0, types.ErrMarketNotFound
	}
	k.maybeCloseExpiredMarket(ctx, market)
	if market.Status != types.MarketStatus_MARKET_STATUS_OPEN {
		return 0, errors.Wrap(types.ErrInvalidMarketStatus, "market must be OPEN")
	}
	if ctx.BlockTime().Unix() < market.OpenTime {
		return 0, errors.Wrap(types.ErrInvalidMarketStatus, "market is not open yet")
	}
	if ctx.BlockTime().Unix() >= market.CloseTime {
		return 0, errors.Wrap(types.ErrInvalidMarketStatus, "market is closed")
	}

	params := k.GetParams(ctx)
	if err := params.Validate(); err != nil {
		return 0, err
	}
	if params.MaxOrderLifetimeBh > math.MaxInt64 {
		return 0, errors.Wrap(types.ErrInvalidRequest, "max_order_lifetime_bh too large")
	}
	if msg.ExpireBh <= ctx.BlockHeight() {
		return 0, errors.Wrap(types.ErrInvalidRequest, "expire_bh must be greater than current block height")
	}
	maxExpireBh := ctx.BlockHeight() + int64(params.MaxOrderLifetimeBh)
	if msg.ExpireBh > maxExpireBh {
		return 0, errors.Wrap(types.ErrInvalidRequest, "expire_bh exceeds max_order_lifetime_bh")
	}

	traderAddr := sdk.MustAccAddressFromBech32(msg.Trader)
	if err := k.enforceOpenOrderLimit(ctx, traderAddr, msg.MarketId, params); err != nil {
		return 0, errors.Wrap(types.ErrInvalidRequest, err.Error())
	}

	amount, _ := sdkmath.NewIntFromString(msg.Amount)
	if err := k.ensurePlaceOrderCapacity(ctx, market, traderAddr, msg.Side, msg.OrderType, msg.LimitPrice, msg.WorstPrice, amount); err != nil {
		return 0, errors.Wrap(types.ErrInsufficientFunds, err.Error())
	}

	id := k.NextOrderID(ctx)
	order := &types.Order{
		Id:              id,
		MarketId:        msg.MarketId,
		Trader:          msg.Trader,
		Side:            msg.Side,
		OrderType:       msg.OrderType,
		Amount:          amount.String(),
		FilledAmount:    sdkmath.ZeroInt().String(),
		RemainingAmount: amount.String(),
		LimitPrice:      msg.LimitPrice,
		WorstPrice:      msg.WorstPrice,
		Status:          types.OrderStatus_ORDER_STATUS_OPEN,
		CreatedBh:       ctx.BlockHeight(),
		ExpireBh:        msg.ExpireBh,
		ClosedBh:        0,
	}

	if err := k.assertOrderInvariant(order); err != nil {
		return 0, err
	}
	if err := k.incrementOpenOrderCounts(ctx, traderAddr, order.MarketId); err != nil {
		return 0, err
	}

	k.SetOrder(ctx, order)
	k.refreshMarketBookPrices(ctx, market)
	k.SetMarket(ctx, market)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventOrderPlaced,
			sdk.NewAttribute(types.AttributeKeyMarketID, intToStr(order.MarketId)),
			sdk.NewAttribute(types.AttributeKeyOrderID, intToStr(order.Id)),
			sdk.NewAttribute(types.AttributeKeyAddress, order.Trader),
			sdk.NewAttribute(types.AttributeKeyOrderSide, order.Side.String()),
			sdk.NewAttribute(types.AttributeKeyOrderType, order.OrderType.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, order.Amount),
		),
	)

	return id, nil
}

func (k Keeper) CancelOrder(ctx sdk.Context, msg *types.MsgCancelOrder) error {
	if err := msg.ValidateBasic(); err != nil {
		return errors.Wrap(types.ErrInvalidRequest, err.Error())
	}

	order, found := k.GetOrder(ctx, msg.MarketId, msg.OrderId)
	if !found {
		return types.ErrOrderNotFound
	}
	if order.Trader != msg.Trader {
		return types.ErrUnauthorized
	}
	if err := k.expireOrderIfNeeded(ctx, order); err != nil {
		return err
	}
	if order.Status != types.OrderStatus_ORDER_STATUS_OPEN && order.Status != types.OrderStatus_ORDER_STATUS_PARTIALLY_FILLED {
		return types.ErrInvalidOrderStatus
	}
	filled, err := parseNonNegativeInt(order.FilledAmount, "filled_amount")
	if err != nil {
		return err
	}

	prevStatus := order.Status
	order.Status = types.OrderStatus_ORDER_STATUS_CANCELLED
	if err := k.onOrderStatusTransition(ctx, order, prevStatus); err != nil {
		return err
	}
	if filled.IsZero() {
		k.DeleteOrder(ctx, order.MarketId, order.Id)
	} else {
		k.SetOrder(ctx, order)
	}
	if market, found := k.GetMarket(ctx, order.MarketId); found {
		k.refreshMarketBookPrices(ctx, market)
		k.SetMarket(ctx, market)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventOrderCancelled,
			sdk.NewAttribute(types.AttributeKeyMarketID, intToStr(order.MarketId)),
			sdk.NewAttribute(types.AttributeKeyOrderID, intToStr(order.Id)),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Trader),
		),
	)

	return nil
}

func (k Keeper) SplitPosition(ctx sdk.Context, msg *types.MsgSplitPosition) error {
	if err := msg.ValidateBasic(); err != nil {
		return errors.Wrap(types.ErrInvalidRequest, err.Error())
	}

	market, found := k.GetMarket(ctx, msg.MarketId)
	if !found {
		return types.ErrMarketNotFound
	}
	k.maybeCloseExpiredMarket(ctx, market)
	if market.Status != types.MarketStatus_MARKET_STATUS_OPEN {
		return errors.Wrap(types.ErrInvalidMarketStatus, "market must be OPEN")
	}
	if ctx.BlockTime().Unix() < market.OpenTime {
		return errors.Wrap(types.ErrInvalidMarketStatus, "market is not open yet")
	}
	if ctx.BlockTime().Unix() >= market.CloseTime {
		return errors.Wrap(types.ErrInvalidMarketStatus, "market is closed")
	}

	amount, _ := sdkmath.NewIntFromString(msg.Amount)
	traderAddr := sdk.MustAccAddressFromBech32(msg.Trader)
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	if err := k.transferCollateralBetweenAccounts(ctx, market, traderAddr, moduleAddr, amount); err != nil {
		return err
	}

	pos := k.getPositionOrDefault(ctx, market.Id, traderAddr)
	yes, lockedYes, no, lockedNo, err := k.mustPositionInts(pos)
	if err != nil {
		return err
	}
	yes = yes.Add(amount)
	no = no.Add(amount)
	k.mustSetPositionInts(pos, yes, lockedYes, no, lockedNo)
	if err := k.assertPositionInvariant(pos); err != nil {
		return err
	}
	k.SetPosition(ctx, pos)

	yesTotal, noTotal, err := k.getMarketShareInts(market)
	if err != nil {
		return err
	}
	market.TotalYesShares = yesTotal.Add(amount).String()
	market.TotalNoShares = noTotal.Add(amount).String()
	if err := k.ValidateMarketInvariants(market); err != nil {
		return err
	}
	k.SetMarket(ctx, market)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventSplit,
			sdk.NewAttribute(types.AttributeKeyMarketID, intToStr(market.Id)),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Trader),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	)

	return nil
}

func (k Keeper) MergePosition(ctx sdk.Context, msg *types.MsgMergePosition) error {
	if err := msg.ValidateBasic(); err != nil {
		return errors.Wrap(types.ErrInvalidRequest, err.Error())
	}

	market, found := k.GetMarket(ctx, msg.MarketId)
	if !found {
		return types.ErrMarketNotFound
	}
	k.maybeCloseExpiredMarket(ctx, market)
	if market.Status != types.MarketStatus_MARKET_STATUS_OPEN {
		return errors.Wrap(types.ErrInvalidMarketStatus, "market must be OPEN")
	}
	if ctx.BlockTime().Unix() < market.OpenTime {
		return errors.Wrap(types.ErrInvalidMarketStatus, "market is not open yet")
	}
	if ctx.BlockTime().Unix() >= market.CloseTime {
		return errors.Wrap(types.ErrInvalidMarketStatus, "market is closed")
	}

	amount, _ := sdkmath.NewIntFromString(msg.Amount)
	traderAddr := sdk.MustAccAddressFromBech32(msg.Trader)
	pos := k.getPositionOrDefault(ctx, market.Id, traderAddr)
	yes, lockedYes, no, lockedNo, err := k.mustPositionInts(pos)
	if err != nil {
		return err
	}

	freeYes := yes.Sub(lockedYes)
	freeNo := no.Sub(lockedNo)
	if freeYes.LT(amount) {
		return errors.Wrap(types.ErrInsufficientFunds, "insufficient YES shares")
	}
	if freeNo.LT(amount) {
		return errors.Wrap(types.ErrInsufficientFunds, "insufficient NO shares")
	}

	yes = yes.Sub(amount)
	no = no.Sub(amount)
	k.mustSetPositionInts(pos, yes, lockedYes, no, lockedNo)
	if err := k.assertPositionInvariant(pos); err != nil {
		return err
	}
	k.SetPosition(ctx, pos)

	yesTotal, noTotal, err := k.getMarketShareInts(market)
	if err != nil {
		return err
	}
	if yesTotal.LT(amount) || noTotal.LT(amount) {
		return fmt.Errorf("invalid market share totals")
	}
	market.TotalYesShares = yesTotal.Sub(amount).String()
	market.TotalNoShares = noTotal.Sub(amount).String()
	if err := k.ValidateMarketInvariants(market); err != nil {
		return err
	}
	k.SetMarket(ctx, market)

	if err := k.transferCollateralFromModule(ctx, market, traderAddr, amount); err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventMerge,
			sdk.NewAttribute(types.AttributeKeyMarketID, intToStr(market.Id)),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Trader),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	)

	return nil
}

func (k Keeper) ResolveMarket(ctx sdk.Context, msg *types.MsgResolveMarket) error {
	if err := msg.ValidateBasic(); err != nil {
		return errors.Wrap(types.ErrInvalidRequest, err.Error())
	}

	market, found := k.GetMarket(ctx, msg.MarketId)
	if !found {
		return types.ErrMarketNotFound
	}
	if market.Resolver != msg.Resolver {
		return types.ErrUnauthorized
	}
	if market.Status != types.MarketStatus_MARKET_STATUS_OPEN && market.Status != types.MarketStatus_MARKET_STATUS_CLOSED {
		return errors.Wrap(types.ErrInvalidMarketStatus, "market must be OPEN or CLOSED")
	}
	// resolve_time = 0 means resolver can resolve at any time.
	if market.ResolveTime > 0 && ctx.BlockTime().Unix() < market.ResolveTime {
		return fmt.Errorf("cannot resolve before resolve_time")
	}

	creatorAddr := sdk.MustAccAddressFromBech32(market.Creator)
	createBond, err := parseNonNegativeInt(market.CreateMarketBond, "create_market_bond")
	if err != nil {
		return err
	}
	if createBond.IsPositive() {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx,
			types.ModuleName,
			creatorAddr,
			sdk.NewCoins(sdk.NewCoin(market.CreateMarketBondDenom, createBond)),
		); err != nil {
			return err
		}
	}

	market.Status = types.MarketStatus_MARKET_STATUS_RESOLVED
	market.WinningOutcome = msg.WinningOutcome
	if strings.TrimSpace(msg.ResolutionSource) != "" {
		market.ResolutionSource = msg.ResolutionSource
	}
	market.CreateMarketBond = sdkmath.ZeroInt().String()

	if err := k.ValidateMarketInvariants(market); err != nil {
		return err
	}

	k.SetMarket(ctx, market)
	winStr, _ := parseOutcome(msg.WinningOutcome)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventMarketResolved,
			sdk.NewAttribute(types.AttributeKeyMarketID, intToStr(market.Id)),
			sdk.NewAttribute(types.AttributeKeyWinning, winStr),
			sdk.NewAttribute(types.AttributeKeyBond, createBond.String()),
			sdk.NewAttribute(types.AttributeKeyBondDenom, market.CreateMarketBondDenom),
		),
	)

	return nil
}

func (k Keeper) RequestResolve(ctx sdk.Context, msg *types.MsgRequestResolve) error {
	if err := msg.ValidateBasic(); err != nil {
		return errors.Wrap(types.ErrInvalidRequest, err.Error())
	}

	market, found := k.GetMarket(ctx, msg.MarketId)
	if !found {
		return types.ErrMarketNotFound
	}
	if market.Creator != msg.Creator {
		return types.ErrUnauthorized
	}
	if market.Status != types.MarketStatus_MARKET_STATUS_OPEN && market.Status != types.MarketStatus_MARKET_STATUS_CLOSED {
		return errors.Wrap(types.ErrInvalidMarketStatus, "market must be OPEN or CLOSED")
	}

	req := &types.ResolveRequest{
		MarketId:         msg.MarketId,
		Creator:          msg.Creator,
		RequestedOutcome: msg.RequestedOutcome,
		RequestedSource:  msg.RequestedSource,
		RequestedBh:      ctx.BlockHeight(),
	}
	k.SetResolveRequest(ctx, req)

	outcomeStr, _ := parseOutcome(msg.RequestedOutcome)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventResolveRequested,
			sdk.NewAttribute(types.AttributeKeyMarketID, intToStr(msg.MarketId)),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyWinning, outcomeStr),
			sdk.NewAttribute("requested_source", msg.RequestedSource),
			sdk.NewAttribute("requested_bh", intToStr(uint64(ctx.BlockHeight()))),
		),
	)

	return nil
}

func (k Keeper) VoidMarket(ctx sdk.Context, msg *types.MsgVoidMarket) error {
	if err := msg.ValidateBasic(); err != nil {
		return errors.Wrap(types.ErrInvalidRequest, err.Error())
	}

	market, found := k.GetMarket(ctx, msg.MarketId)
	if !found {
		return types.ErrMarketNotFound
	}
	if market.Resolver != msg.Resolver {
		return types.ErrUnauthorized
	}
	if market.Status != types.MarketStatus_MARKET_STATUS_OPEN && market.Status != types.MarketStatus_MARKET_STATUS_CLOSED {
		return errors.Wrap(types.ErrInvalidMarketStatus, "market must be OPEN or CLOSED")
	}

	creatorAddr := sdk.MustAccAddressFromBech32(market.Creator)
	createBond, err := parseNonNegativeInt(market.CreateMarketBond, "create_market_bond")
	if err != nil {
		return err
	}
	if createBond.IsPositive() {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx,
			types.ModuleName,
			creatorAddr,
			sdk.NewCoins(sdk.NewCoin(market.CreateMarketBondDenom, createBond)),
		); err != nil {
			return err
		}
	}

	market.Status = types.MarketStatus_MARKET_STATUS_VOIDED
	market.WinningOutcome = types.Outcome_OUTCOME_UNSPECIFIED
	market.CreateMarketBond = sdkmath.ZeroInt().String()
	k.SetMarket(ctx, market)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventMarketVoided,
			sdk.NewAttribute(types.AttributeKeyMarketID, intToStr(market.Id)),
			sdk.NewAttribute("reason", msg.Reason),
			sdk.NewAttribute(types.AttributeKeyBond, createBond.String()),
			sdk.NewAttribute(types.AttributeKeyBondDenom, market.CreateMarketBondDenom),
		),
	)

	return nil
}
