package keeper

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/paxi-web3/paxi/x/swap/types"
)

type queryServer struct {
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
		CodeId:       params.CodeID,
		SwapFeeBps:   params.SwapFeeBPS,
		MinLiquidity: params.MinLiquidity,
	}, nil
}

func (q *queryServer) Position(ctx context.Context, req *types.QueryPositionRequest) (*types.QueryPositionResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

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

	return &types.QueryPositionResponse{
		Position:      &pos,
		ExpectedPaxi:  expectedPaxi.String(),
		ExpectedPrc20: expectedPRC20.String(),
	}, nil
}

func (q *queryServer) Pool(ctx context.Context, req *types.QueryPoolRequest) (*types.QueryPoolResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	pool, found := q.k.GetPool(sdkCtx, req.Prc20)
	if !found || pool.ReservePaxi.IsZero() || pool.ReservePRC20.IsZero() {
		return nil, fmt.Errorf("pool not found or empty")
	}

	pricePaxiPerPRC20 := sdkmath.LegacyNewDecFromInt(pool.ReservePaxi).
		Quo(sdkmath.LegacyNewDecFromInt(pool.ReservePRC20))
	pricePRC20PerPaxi := sdkmath.LegacyNewDecFromInt(pool.ReservePRC20).
		Quo(sdkmath.LegacyNewDecFromInt(pool.ReservePaxi))

	return &types.QueryPoolResponse{
		Prc20:             req.Prc20,
		ReservePaxi:       pool.ReservePaxi.String(),
		ReservePrc20:      pool.ReservePRC20.String(),
		PricePaxiPerPrc20: pricePaxiPerPRC20.String(),
		PricePrc20PerPaxi: pricePRC20PerPaxi.String(),
		TotalShares:       pool.TotalShares.String(),
	}, nil
}

func (q *queryServer) AllPools(ctx context.Context, req *types.QueryAllPoolsRequest) (*types.QueryAllPoolsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Max limit is 100
	if req.Pagination != nil && req.Pagination.Limit > 100 {
		req.Pagination.Limit = 100
	}

	store := sdkCtx.KVStore(q.k.storeKey)
	poolStore := prefix.NewStore(store, types.PoolPrefix)

	var pools []*types.QueryPoolResponse
	pageRes, err := query.Paginate(poolStore, req.Pagination, func(key []byte, value []byte) error {
		var pool types.PoolProto
		if err := q.k.cdc.Unmarshal(value, &pool); err != nil {
			return err
		}

		reservePaxi, ok := sdkmath.NewIntFromString(pool.ReservePaxi)
		if !ok || reservePaxi.IsZero() {
			// Skip empty pools
			return nil
		}
		reservePrc20, ok := sdkmath.NewIntFromString(pool.ReservePrc20)
		if !ok || reservePrc20.IsZero() {
			return nil
		}

		paxiDec, err := sdkmath.LegacyNewDecFromStr(pool.ReservePaxi)
		if err != nil {
			return fmt.Errorf("invalid ReservePaxi: %w", err)
		}
		prc20Dec, err := sdkmath.LegacyNewDecFromStr(pool.ReservePrc20)
		if err != nil {
			return fmt.Errorf("invalid ReservePrc20: %w", err)
		}

		pools = append(pools, &types.QueryPoolResponse{
			Prc20:             pool.Prc20,
			ReservePaxi:       pool.ReservePaxi,
			ReservePrc20:      pool.ReservePrc20,
			PricePaxiPerPrc20: paxiDec.Quo(prc20Dec).String(),
			PricePrc20PerPaxi: prc20Dec.Quo(paxiDec).String(),
			TotalShares:       pool.TotalShares,
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryAllPoolsResponse{
		Pools:      pools,
		Pagination: pageRes,
	}, nil
}
