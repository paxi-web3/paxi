package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/x/custommint/types"
)

type queryServer struct {
	k Keeper
}

func NewQueryServer(k Keeper) types.QueryServer {
	return &queryServer{k: k}
}

func (q *queryServer) TotalMinted(ctx context.Context, req *types.QueryTotalMintedRequest) (*types.QueryTotalMintedResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	total := q.k.GetTotalMinted(sdkCtx)

	return &types.QueryTotalMintedResponse{
		Amount: sdk.NewCoin(types.DefaultDenom, total),
	}, nil
}
