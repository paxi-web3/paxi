package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/x/customwasm/types"
)

type queryServer struct {
	types.UnimplementedQueryServer
	k Keeper
}

func NewQueryServer(k Keeper) types.QueryServer {
	return &queryServer{k: k}
}

func (q *queryServer) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := q.k.GetParams(sdkCtx)

	return &types.QueryParamsResponse{
		StoreCodeBaseGas:    params.StoreCodeBaseGas,
		StoreCodeMultiplier: params.StoreCodeMultiplier,
		InstBaseGas:         params.InstBaseGas,
		InstMultiplier:      params.InstMultiplier,
	}, nil
}
