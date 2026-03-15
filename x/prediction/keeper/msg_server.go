package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/x/prediction/types"
)

type msgServer struct {
	Keeper
}

func NewMsgServerImpl(k Keeper) types.MsgServer {
	return &msgServer{Keeper: k}
}

func (m msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := msg.ValidateBasic(); err != nil {
		return nil, errors.Wrap(types.ErrInvalidRequest, err.Error())
	}
	if msg.Authority != m.authority {
		return nil, types.ErrUnauthorized
	}

	params := types.Params{
		MaxBatchSize:            msg.Params.MaxBatchSize,
		CreateMarketBond:        msg.Params.CreateMarketBond,
		CreateMarketBondDenom:   msg.Params.CreateMarketBondDenom,
		MarketFeeBps:            msg.Params.MarketFeeBps,
		ResolverFeeSharePercent: msg.Params.ResolverFeeSharePercent,
		MaxOrderLifetimeBh:      msg.Params.MaxOrderLifetimeBh,
		MaxOpenOrdersPerUser:    msg.Params.MaxOpenOrdersPerUser,
		MaxOpenOrdersPerMarket:  msg.Params.MaxOpenOrdersPerMarket,
		OrderPruneIntervalBh:    msg.Params.OrderPruneIntervalBh,
		OrderPruneRetainBh:      msg.Params.OrderPruneRetainBh,
		OrderPruneScanLimit:     msg.Params.OrderPruneScanLimit,
		OrderPruneDeleteLimit:   msg.Params.OrderPruneDeleteLimit,
	}
	if err := params.Validate(); err != nil {
		return nil, err
	}

	m.SetParams(ctx, params)
	return &types.MsgUpdateParamsResponse{}, nil
}

func (m msgServer) CreateMarket(goCtx context.Context, msg *types.MsgCreateMarket) (*types.MsgCreateMarketResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	marketID, err := m.Keeper.CreateMarket(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgCreateMarketResponse{MarketId: marketID}, nil
}

func (m msgServer) PlaceOrder(goCtx context.Context, msg *types.MsgPlaceOrder) (*types.MsgPlaceOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	orderID, err := m.Keeper.PlaceOrder(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgPlaceOrderResponse{OrderId: orderID}, nil
}

func (m msgServer) CancelOrder(goCtx context.Context, msg *types.MsgCancelOrder) (*types.MsgCancelOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.CancelOrder(ctx, msg); err != nil {
		return nil, err
	}
	return &types.MsgCancelOrderResponse{}, nil
}

func (m msgServer) SplitPosition(goCtx context.Context, msg *types.MsgSplitPosition) (*types.MsgSplitPositionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.SplitPosition(ctx, msg); err != nil {
		return nil, err
	}
	return &types.MsgSplitPositionResponse{}, nil
}

func (m msgServer) MergePosition(goCtx context.Context, msg *types.MsgMergePosition) (*types.MsgMergePositionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.MergePosition(ctx, msg); err != nil {
		return nil, err
	}
	return &types.MsgMergePositionResponse{}, nil
}

func (m msgServer) ApplyTradeBatch(goCtx context.Context, msg *types.MsgApplyTradeBatch) (*types.MsgApplyTradeBatchResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	return m.Keeper.ApplyTradeBatch(ctx, msg)
}

func (m msgServer) ResolveMarket(goCtx context.Context, msg *types.MsgResolveMarket) (*types.MsgResolveMarketResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.ResolveMarket(ctx, msg); err != nil {
		return nil, err
	}
	return &types.MsgResolveMarketResponse{}, nil
}

func (m msgServer) VoidMarket(goCtx context.Context, msg *types.MsgVoidMarket) (*types.MsgVoidMarketResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.VoidMarket(ctx, msg); err != nil {
		return nil, err
	}
	return &types.MsgVoidMarketResponse{}, nil
}

func (m msgServer) ClaimPayout(goCtx context.Context, msg *types.MsgClaimPayout) (*types.MsgClaimPayoutResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	return m.Keeper.ClaimPayout(ctx, msg)
}

func (m msgServer) ClaimVoidRefund(goCtx context.Context, msg *types.MsgClaimVoidRefund) (*types.MsgClaimVoidRefundResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	return m.Keeper.ClaimVoidRefund(ctx, msg)
}
