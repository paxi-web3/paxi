package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"
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

func (q *queryServer) CirculatingSupply(ctx context.Context, req *types.QueryCirculatingSupplyRequest) (*types.QueryCirculatingSupplyResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cirSupply := q.k.bankKeeper.GetSupply(ctx, types.DefaultDenom).Amount.Int64() - q.k.GetLockedVesting(sdkCtx).TruncateInt64()

	return &types.QueryCirculatingSupplyResponse{
		Amount: sdk.NewCoin(types.DefaultDenom, sdkmath.NewInt(cirSupply)),
	}, nil
}

func (q *queryServer) TotalSupply(ctx context.Context, req *types.QueryTotalSupplyRequest) (*types.QueryTotalSupplyResponse, error) {
	cirSupply := q.k.bankKeeper.GetSupply(ctx, types.DefaultDenom).Amount.Int64()

	return &types.QueryTotalSupplyResponse{
		Amount: sdk.NewCoin(types.DefaultDenom, sdkmath.NewInt(cirSupply)),
	}, nil
}
