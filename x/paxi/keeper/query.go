package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/x/paxi/types"
)

type queryServer struct {
	k Keeper
}

func NewQueryServer(k Keeper) types.QueryServer {
	return &queryServer{k: k}
}

func (q *queryServer) LockedVesting(ctx context.Context, req *types.QueryLockedVestingRequest) (*types.QueryLockedVestingResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	lockedVesting := q.k.GetLockedVesting(sdkCtx)

	return &types.QueryLockedVestingResponse{
		Amount: sdk.NewCoin(types.DefaultDenom, lockedVesting.TruncateInt()),
	}, nil
}
