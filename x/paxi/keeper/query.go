package keeper

import (
	"context"
	"fmt"

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
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	lockedVesting := q.k.GetLockedVesting(sdkCtx)
	coin := sdk.NewCoin(types.DefaultDenom, lockedVesting.TruncateInt())

	return &types.QueryLockedVestingResponse{
		Amount: &coin,
	}, nil
}

func (q *queryServer) CirculatingSupply(ctx context.Context, req *types.QueryCirculatingSupplyRequest) (*types.QueryCirculatingSupplyResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cirSupply := q.k.bankKeeper.GetSupply(ctx, types.DefaultDenom).Amount.Int64() - q.k.GetLockedVesting(sdkCtx).TruncateInt64()
	coin := sdk.NewCoin(types.DefaultDenom, sdkmath.NewInt(cirSupply))

	return &types.QueryCirculatingSupplyResponse{
		Amount: &coin,
	}, nil
}

func (q *queryServer) TotalSupply(ctx context.Context, req *types.QueryTotalSupplyRequest) (*types.QueryTotalSupplyResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	cirSupply := q.k.bankKeeper.GetSupply(ctx, types.DefaultDenom).Amount.Int64()
	coin := sdk.NewCoin(types.DefaultDenom, sdkmath.NewInt(cirSupply))

	return &types.QueryTotalSupplyResponse{
		Amount: &coin,
	}, nil
}

func (q *queryServer) LastBlockGasUsed(ctx context.Context, req *types.QueryLastBlockGasUsedRequest) (*types.QueryLastBlockGasUsedResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	gasUsed := q.k.blockStatusKeeper.GetLastBlockGasUsed()
	return &types.QueryLastBlockGasUsedResponse{
		GasUsed: gasUsed,
	}, nil
}

func (q *queryServer) TotalTxs(ctx context.Context, req *types.QueryTotalTxsRequest) (*types.QueryTotalTxsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	totalTxs := q.k.blockStatusKeeper.GetTotalTxs()
	return &types.QueryTotalTxsResponse{
		TotalTxs: totalTxs,
	}, nil
}

func (q *queryServer) UnlockSchedules(ctx context.Context, req *types.QueryUnlockSchedulesRequest) (*types.QueryUnlockSchedulesResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	unlockSchedules := q.k.GetUnlockSchedules(sdkCtx)
	return &types.QueryUnlockSchedulesResponse{
		UnlockSchedules: unlockSchedules,
	}, nil
}

func (q *queryServer) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := q.k.GetParams(sdkCtx)

	return &types.QueryParamsResponse{
		ExtraGasPerNewAccount: params.ExtraGasPerNewAccount,
	}, nil
}
