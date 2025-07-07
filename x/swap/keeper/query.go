package keeper

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/x/swap/types"
)

type queryServer struct {
	k Keeper
}

func NewQueryServer(k Keeper) types.QueryServer {
	return &queryServer{k: k}
}

func (q *queryServer) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := q.k.GetParams(sdkCtx)

	return &types.QueryParamsResponse{
		CodeId:     params.CodeID,
		SwapFeeBps: params.SwapFeeBPS,
	}, nil
}

func (q *queryServer) Position(ctx context.Context, req *types.QueryPositionRequest) (*types.QueryPositionResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	creator, err := sdk.AccAddressFromBech32(req.Creator)
	if err != nil {
		return nil, err
	}

	pos, found := q.k.GetPosition(sdkCtx, req.Prc20, creator)
	if !found {
		return &types.QueryPositionResponse{
			Position: &types.ProviderPosition{
				Creator:  req.Creator,
				Prc20:    req.Prc20,
				LpAmount: "0",
			},
			ExpectedPaxi:  "0",
			ExpectedPrc20: "0",
		}, nil
	}

	pool, found := q.k.GetPool(sdkCtx, req.Prc20)
	if !found || pool.TotalShares.IsZero() {
		return &types.QueryPositionResponse{
			Position:      &pos,
			ExpectedPaxi:  "0",
			ExpectedPrc20: "0",
		}, nil
	}

	// Calculate expected PAXI and PRC20 based on user's share in the pool
	lpAmount, err := sdkmath.LegacyNewDecFromStr(pos.LpAmount)
	if err != nil {
		return nil, fmt.Errorf("invalid LP amount: %w", err)
	}

	totalShares := sdkmath.LegacyNewDecFromInt(pool.TotalShares)
	reservePaxi := sdkmath.LegacyNewDecFromInt(pool.ReservePaxi)
	reservePRC20 := sdkmath.LegacyNewDecFromInt(pool.ReservePRC20)

	userShareRatio := lpAmount.Quo(totalShares)
	expectedPaxi := reservePaxi.Mul(userShareRatio).TruncateInt()
	expectedPRC20 := reservePRC20.Mul(userShareRatio).TruncateInt()

	depositedLp, err := sdkmath.LegacyNewDecFromStr(pos.DepositedLp)
	if err != nil {
		return nil, fmt.Errorf("invalid deposited LP amount: %w", err)
	}
	profitLp := lpAmount.Sub(depositedLp).TruncateInt()
	if profitLp.IsNegative() {
		profitLp = sdkmath.ZeroInt()
	}

	return &types.QueryPositionResponse{
		Position:      &pos,
		ExpectedPaxi:  expectedPaxi.String(),
		ExpectedPrc20: expectedPRC20.String(),
		ProfitLp:      profitLp.String(),
	}, nil
}
