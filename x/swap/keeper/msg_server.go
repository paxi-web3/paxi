package keeper

import (
	"context"
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/paxi-web3/paxi/x/swap/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(k Keeper) types.MsgServer {
	return &msgServer{Keeper: k}
}

func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.Authority != k.authority {
		return nil, sdkerrors.ErrUnauthorized
	}

	parsedParams := types.Params{
		CodeID:     msg.Params.CodeId,
		SwapFeeBPS: msg.Params.SwapFeeBps,
	}

	if err := parsedParams.Validate(); err != nil {
		return nil, err
	}

	store := k.storeService.OpenKVStore(ctx)
	bz, err := json.Marshal(&parsedParams)
	if err != nil {
		return nil, err
	}
	store.Set(types.KeyParams, bz)

	return &types.MsgUpdateParamsResponse{}, nil
}

func (m msgServer) ProvideLiquidity(goCtx context.Context, msg *types.MsgProvideLiquidity) (*types.MsgProvideLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := m.Keeper.ProvideLiquidity(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgProvideLiquidityResponse{}, nil
}

func (m msgServer) WithdrawLiquidity(goCtx context.Context, msg *types.MsgWithdrawLiquidity) (*types.MsgWithdrawLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := m.Keeper.WithdrawLiquidity(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgWithdrawLiquidityResponse{}, nil
}

func (m msgServer) Swap(goCtx context.Context, msg *types.MsgSwap) (*types.MsgSwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := m.Keeper.Swap(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &types.MsgSwapResponse{}, nil
}
