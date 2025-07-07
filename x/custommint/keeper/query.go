package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/x/custommint/types"
)

type queryServer struct {
	types.UnimplementedQueryServer
	k Keeper
}

func NewQueryServer(k Keeper) types.QueryServer {
	return &queryServer{k: k}
}

func (q *queryServer) TotalMinted(ctx context.Context, req *types.QueryTotalMintedRequest) (*types.QueryTotalMintedResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	total := q.k.GetTotalMinted(sdkCtx)
	coin := sdk.NewCoin(types.DefaultDenom, total)

	return &types.QueryTotalMintedResponse{
		Amount: &coin,
	}, nil
}

func (q *queryServer) TotalBurned(ctx context.Context, req *types.QueryTotalBurnedRequest) (*types.QueryTotalBurnedResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	total := q.k.GetTotalBurned(sdkCtx)
	coin := sdk.NewCoin(types.DefaultDenom, total)

	return &types.QueryTotalBurnedResponse{
		Amount: &coin,
	}, nil
}

func (q *queryServer) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := q.k.GetParams(sdkCtx)

	return &types.QueryParamsResponse{
		BurnThreshold:       params.BurnRatio.String(),
		BurnRatio:           params.BurnRatio.String(),
		BlocksPerYear:       params.BlocksPerYear,
		FirstYearInflation:  params.FirstYearInflation.String(),
		SecondYearInflation: params.SecondYearInflation.String(),
		OtherYearInflation:  params.OtherYearInflation.String(),
	}, nil
}
