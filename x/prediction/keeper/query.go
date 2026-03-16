package keeper

import (
	"context"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/paxi-web3/paxi/x/prediction/types"
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
		MaxBatchSize:            params.MaxBatchSize,
		CreateMarketBond:        params.CreateMarketBond,
		CreateMarketBondDenom:   params.CreateMarketBondDenom,
		MarketFeeBps:            params.MarketFeeBps,
		ResolverFeeSharePercent: params.ResolverFeeSharePercent,
		MaxOrderLifetimeBh:      params.MaxOrderLifetimeBh,
		MaxOpenOrdersPerUser:    params.MaxOpenOrdersPerUser,
		MaxOpenOrdersPerMarket:  params.MaxOpenOrdersPerMarket,
		OrderPruneIntervalBh:    params.OrderPruneIntervalBh,
		OrderPruneRetainBh:      params.OrderPruneRetainBh,
		OrderPruneScanLimit:     params.OrderPruneScanLimit,
		OrderPruneDeleteLimit:   params.OrderPruneDeleteLimit,
	}, nil
}

func (q *queryServer) Market(ctx context.Context, req *types.QueryMarketRequest) (*types.QueryMarketResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.MarketId == 0 {
		return nil, fmt.Errorf("market_id must be > 0")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	market, found := q.k.GetMarket(sdkCtx, req.MarketId)
	if !found {
		return nil, types.ErrMarketNotFound
	}

	return &types.QueryMarketResponse{Market: market}, nil
}

func (q *queryServer) Markets(ctx context.Context, req *types.QueryMarketsRequest) (*types.QueryMarketsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	status, err := parseMarketStatusFilter(req.Status)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := sdkCtx.KVStore(q.k.storeKey)
	marketStore := prefix.NewStore(store, types.MarketPrefix)

	markets := make([]*types.Market, 0)
	pageRes, err := query.Paginate(marketStore, req.Pagination, func(_ []byte, value []byte) error {
		market := &types.Market{}
		if err := q.k.cdc.Unmarshal(value, market); err != nil {
			return err
		}
		if status != types.MarketStatus_MARKET_STATUS_UNSPECIFIED && market.Status != status {
			return nil
		}
		markets = append(markets, market)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryMarketsResponse{Markets: markets, Pagination: pageRes}, nil
}

func (q *queryServer) Order(ctx context.Context, req *types.QueryOrderRequest) (*types.QueryOrderResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.MarketId == 0 {
		return nil, fmt.Errorf("market_id must be > 0")
	}
	if req.OrderId == 0 {
		return nil, fmt.Errorf("order_id must be > 0")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	order, found := q.k.GetOrder(sdkCtx, req.MarketId, req.OrderId)
	if !found {
		return nil, types.ErrOrderNotFound
	}

	return &types.QueryOrderResponse{Order: order}, nil
}

func (q *queryServer) OrderById(ctx context.Context, req *types.QueryOrderByIdRequest) (*types.QueryOrderByIdResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.OrderId == 0 {
		return nil, fmt.Errorf("order_id must be > 0")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	order, found := q.k.GetOrderByID(sdkCtx, req.OrderId)
	if !found {
		return nil, types.ErrOrderNotFound
	}

	return &types.QueryOrderByIdResponse{Order: order}, nil
}

func (q *queryServer) OrdersByAddress(ctx context.Context, req *types.QueryOrdersByAddressRequest) (*types.QueryOrdersByAddressResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if _, err := sdk.AccAddressFromBech32(req.Address); err != nil {
		return nil, fmt.Errorf("invalid address: %w", err)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := sdkCtx.KVStore(q.k.storeKey)
	orderStore := prefix.NewStore(store, types.OrderPrefix)

	orders := make([]*types.Order, 0)
	pageRes, err := query.Paginate(orderStore, req.Pagination, func(_ []byte, value []byte) error {
		order := &types.Order{}
		if err := q.k.cdc.Unmarshal(value, order); err != nil {
			return err
		}
		if order.Trader == req.Address {
			orders = append(orders, order)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryOrdersByAddressResponse{Orders: orders, Pagination: pageRes}, nil
}

func (q *queryServer) OrdersByMarket(ctx context.Context, req *types.QueryOrdersByMarketRequest) (*types.QueryOrdersByMarketResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.MarketId == 0 {
		return nil, fmt.Errorf("market_id must be > 0")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := sdkCtx.KVStore(q.k.storeKey)
	marketBz := make([]byte, 8)
	binary.BigEndian.PutUint64(marketBz, req.MarketId)
	orderStore := prefix.NewStore(store, append(types.OrderPrefix, marketBz...))

	orders := make([]*types.Order, 0)
	pageRes, err := query.Paginate(orderStore, req.Pagination, func(_ []byte, value []byte) error {
		order := &types.Order{}
		if err := q.k.cdc.Unmarshal(value, order); err != nil {
			return err
		}
		orders = append(orders, order)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryOrdersByMarketResponse{Orders: orders, Pagination: pageRes}, nil
}

func (q *queryServer) ResolveRequest(ctx context.Context, req *types.QueryResolveRequestRequest) (*types.QueryResolveRequestResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.MarketId == 0 {
		return nil, fmt.Errorf("market_id must be > 0")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	resolveReq, found := q.k.GetResolveRequest(sdkCtx, req.MarketId)
	if !found {
		return &types.QueryResolveRequestResponse{}, nil
	}
	return &types.QueryResolveRequestResponse{ResolveRequest: resolveReq}, nil
}

func (q *queryServer) Position(ctx context.Context, req *types.QueryPositionRequest) (*types.QueryPositionResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.MarketId == 0 {
		return nil, fmt.Errorf("market_id must be > 0")
	}
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, fmt.Errorf("invalid address: %w", err)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	pos, found := q.k.GetPosition(sdkCtx, req.MarketId, addr)
	if !found {
		return &types.QueryPositionResponse{Position: &types.Position{
			MarketId:          req.MarketId,
			Address:           req.Address,
			YesShares:         sdkmath.ZeroInt().String(),
			LockedYesShares:   sdkmath.ZeroInt().String(),
			NoShares:          sdkmath.ZeroInt().String(),
			LockedNoShares:    sdkmath.ZeroInt().String(),
			ClaimedPayout:     false,
			ClaimedVoidRefund: false,
		}}, nil
	}

	return &types.QueryPositionResponse{Position: pos}, nil
}

func (q *queryServer) PositionsByAddress(ctx context.Context, req *types.QueryPositionsByAddressRequest) (*types.QueryPositionsByAddressResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if _, err := sdk.AccAddressFromBech32(req.Address); err != nil {
		return nil, fmt.Errorf("invalid address: %w", err)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := sdkCtx.KVStore(q.k.storeKey)
	positionStore := prefix.NewStore(store, types.PositionPrefix)

	positions := make([]*types.Position, 0)
	pageRes, err := query.Paginate(positionStore, req.Pagination, func(_ []byte, value []byte) error {
		pos := &types.Position{}
		if err := q.k.cdc.Unmarshal(value, pos); err != nil {
			return err
		}
		if pos.Address == req.Address {
			positions = append(positions, pos)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryPositionsByAddressResponse{Positions: positions, Pagination: pageRes}, nil
}

func parseMarketStatusFilter(raw string) (types.MarketStatus, error) {
	val := strings.TrimSpace(raw)
	if val == "" {
		return types.MarketStatus_MARKET_STATUS_UNSPECIFIED, nil
	}

	if iv, err := strconv.ParseInt(val, 10, 32); err == nil {
		status := types.MarketStatus(iv)
		switch status {
		case types.MarketStatus_MARKET_STATUS_UNSPECIFIED,
			types.MarketStatus_MARKET_STATUS_OPEN,
			types.MarketStatus_MARKET_STATUS_CLOSED,
			types.MarketStatus_MARKET_STATUS_RESOLVED,
			types.MarketStatus_MARKET_STATUS_VOIDED:
			return status, nil
		default:
			return types.MarketStatus_MARKET_STATUS_UNSPECIFIED, fmt.Errorf("invalid status value: %s", raw)
		}
	}

	switch strings.ToUpper(val) {
	case "MARKET_STATUS_UNSPECIFIED":
		return types.MarketStatus_MARKET_STATUS_UNSPECIFIED, nil
	case "MARKET_STATUS_OPEN":
		return types.MarketStatus_MARKET_STATUS_OPEN, nil
	case "MARKET_STATUS_CLOSED":
		return types.MarketStatus_MARKET_STATUS_CLOSED, nil
	case "MARKET_STATUS_RESOLVED":
		return types.MarketStatus_MARKET_STATUS_RESOLVED, nil
	case "MARKET_STATUS_VOIDED":
		return types.MarketStatus_MARKET_STATUS_VOIDED, nil
	default:
		return types.MarketStatus_MARKET_STATUS_UNSPECIFIED, fmt.Errorf("invalid status: %s", raw)
	}
}
