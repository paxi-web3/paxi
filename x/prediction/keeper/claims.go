package keeper

import (
	"fmt"

	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/x/prediction/types"
)

func (k Keeper) ClaimPayout(ctx sdk.Context, msg *types.MsgClaimPayout) (*types.MsgClaimPayoutResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, errors.Wrap(types.ErrInvalidRequest, err.Error())
	}

	market, found := k.GetMarket(ctx, msg.MarketId)
	if !found {
		return nil, types.ErrMarketNotFound
	}
	if market.Status != types.MarketStatus_MARKET_STATUS_RESOLVED {
		return nil, errors.Wrap(types.ErrInvalidMarketStatus, "market is not resolved")
	}

	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	pos := k.getPositionOrDefault(ctx, market.Id, creator)
	if pos.ClaimedPayout {
		return nil, types.ErrAlreadyClaimed
	}

	yes, lockedYes, no, lockedNo, err := k.mustPositionInts(pos)
	if err != nil {
		return nil, err
	}

	yesTotal, noTotal, err := k.getMarketShareInts(market)
	if err != nil {
		return nil, err
	}

	burnYes := yes.Add(lockedYes)
	burnNo := no.Add(lockedNo)
	if yesTotal.LT(burnYes) || noTotal.LT(burnNo) {
		return nil, fmt.Errorf("invalid market share totals")
	}

	shareUnit := sdkmath.NewInt(types.CollateralUnit)
	payout := sdkmath.ZeroInt()
	switch market.WinningOutcome {
	case types.Outcome_OUTCOME_YES:
		payout = burnYes.Mul(shareUnit)
	case types.Outcome_OUTCOME_NO:
		payout = burnNo.Mul(shareUnit)
	default:
		return nil, types.ErrInvalidOutcome
	}
	if !payout.IsPositive() {
		return nil, errors.Wrap(types.ErrInsufficientFunds, "no winning shares")
	}

	if err := k.transferCollateralFromModule(ctx, market, creator, payout); err != nil {
		return nil, err
	}

	yesTotal = yesTotal.Sub(burnYes)
	noTotal = noTotal.Sub(burnNo)

	k.mustSetPositionInts(pos, sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt())
	pos.ClaimedPayout = true

	market.TotalYesShares = yesTotal.String()
	market.TotalNoShares = noTotal.String()

	if err := k.assertPositionInvariant(pos); err != nil {
		return nil, err
	}
	if err := k.ValidateMarketInvariants(market); err != nil {
		return nil, err
	}

	k.SetPosition(ctx, pos)
	k.SetMarket(ctx, market)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventPayoutClaimed,
			sdk.NewAttribute(types.AttributeKeyMarketID, intToStr(market.Id)),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyAmount, payout.String()),
		),
	)

	return &types.MsgClaimPayoutResponse{Payout: payout.String()}, nil
}

func (k Keeper) ClaimVoidRefund(ctx sdk.Context, msg *types.MsgClaimVoidRefund) (*types.MsgClaimVoidRefundResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, errors.Wrap(types.ErrInvalidRequest, err.Error())
	}

	market, found := k.GetMarket(ctx, msg.MarketId)
	if !found {
		return nil, types.ErrMarketNotFound
	}
	if market.Status != types.MarketStatus_MARKET_STATUS_VOIDED {
		return nil, errors.Wrap(types.ErrInvalidMarketStatus, "market is not voided")
	}

	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	pos := k.getPositionOrDefault(ctx, market.Id, creator)
	if pos.ClaimedVoidRefund {
		return nil, types.ErrAlreadyClaimed
	}

	yes, lockedYes, no, lockedNo, err := k.mustPositionInts(pos)
	if err != nil {
		return nil, err
	}

	burnYes := yes.Add(lockedYes)
	burnNo := no.Add(lockedNo)
	totalShares := burnYes.Add(burnNo)
	if !totalShares.IsPositive() {
		return nil, errors.Wrap(types.ErrInsufficientFunds, "no shares to refund")
	}

	shareUnit := sdkmath.NewInt(types.CollateralUnit)
	refund := totalShares.Mul(shareUnit).QuoRaw(2)
	if refund.IsPositive() {
		if err := k.transferCollateralFromModule(ctx, market, creator, refund); err != nil {
			return nil, err
		}
	}

	yesTotal, noTotal, err := k.getMarketShareInts(market)
	if err != nil {
		return nil, err
	}
	if yesTotal.LT(burnYes) || noTotal.LT(burnNo) {
		return nil, fmt.Errorf("invalid market share totals")
	}
	yesTotal = yesTotal.Sub(burnYes)
	noTotal = noTotal.Sub(burnNo)

	k.mustSetPositionInts(pos, sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt())
	pos.ClaimedVoidRefund = true

	market.TotalYesShares = yesTotal.String()
	market.TotalNoShares = noTotal.String()

	if err := k.assertPositionInvariant(pos); err != nil {
		return nil, err
	}
	if err := k.ValidateMarketInvariants(market); err != nil {
		return nil, err
	}

	k.SetPosition(ctx, pos)
	k.SetMarket(ctx, market)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventVoidRefundClaimed,
			sdk.NewAttribute(types.AttributeKeyMarketID, intToStr(market.Id)),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyAmount, refund.String()),
		),
	)

	return &types.MsgClaimVoidRefundResponse{Refund: refund.String()}, nil
}
