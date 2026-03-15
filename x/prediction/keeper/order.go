package keeper

import (
	"encoding/binary"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/x/prediction/types"
)

func (k Keeper) GetOrder(ctx sdk.Context, marketID uint64, orderID uint64) (*types.Order, bool) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.OrderStoreKey(marketID, orderID))
	if err != nil || bz == nil {
		return nil, false
	}

	order := &types.Order{}
	if err := k.cdc.Unmarshal(bz, order); err != nil {
		return nil, false
	}

	return order, true
}

func (k Keeper) GetOrderByID(ctx sdk.Context, orderID uint64) (*types.Order, bool) {
	store := k.storeService.OpenKVStore(ctx)
	indexBz, err := store.Get(types.OrderIDIndexKey(orderID))
	if err == nil && len(indexBz) == 8 {
		marketID := binary.BigEndian.Uint64(indexBz)
		order, found := k.GetOrder(ctx, marketID, orderID)
		if found {
			return order, true
		}
		// Clean stale index entry.
		k.mustDelete(store, types.OrderIDIndexKey(orderID))
	}

	// Fallback scan for legacy data without index.
	it, err := prefixIterator(store, types.OrderPrefix)
	if err != nil {
		panic(err)
	}
	defer it.Close()
	for ; it.Valid(); it.Next() {
		order := &types.Order{}
		if err := k.cdc.Unmarshal(it.Value(), order); err != nil {
			continue
		}
		if order.Id != orderID {
			continue
		}
		marketBz := make([]byte, 8)
		binary.BigEndian.PutUint64(marketBz, order.MarketId)
		k.mustSet(store, types.OrderIDIndexKey(orderID), marketBz)
		return order, true
	}

	return nil, false
}

func (k Keeper) SetOrder(ctx sdk.Context, order *types.Order) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := k.cdc.Marshal(order)
	if err != nil {
		panic(err)
	}
	k.mustSet(store, types.OrderStoreKey(order.MarketId, order.Id), bz)

	marketBz := make([]byte, 8)
	binary.BigEndian.PutUint64(marketBz, order.MarketId)
	k.mustSet(store, types.OrderIDIndexKey(order.Id), marketBz)
}

func (k Keeper) DeleteOrder(ctx sdk.Context, marketID uint64, orderID uint64) {
	store := k.storeService.OpenKVStore(ctx)
	k.mustDelete(store, types.OrderStoreKey(marketID, orderID))
	k.mustDelete(store, types.OrderIDIndexKey(orderID))
}

func (k Keeper) GetAllOrders(ctx sdk.Context) []*types.Order {
	store := k.storeService.OpenKVStore(ctx)
	it, err := prefixIterator(store, types.OrderPrefix)
	if err != nil {
		panic(err)
	}
	defer it.Close()

	orders := make([]*types.Order, 0)
	for ; it.Valid(); it.Next() {
		order := &types.Order{}
		if err := k.cdc.Unmarshal(it.Value(), order); err != nil {
			continue
		}
		orders = append(orders, order)
	}
	return orders
}

func (k Keeper) SetNextOrderID(ctx sdk.Context, nextID uint64) {
	if nextID == 0 {
		nextID = 1
	}
	store := k.storeService.OpenKVStore(ctx)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, nextID)
	k.mustSet(store, types.NextOrderIDKey, bz)
}

func (k Keeper) GetNextOrderID(ctx sdk.Context) uint64 {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.NextOrderIDKey)
	if err != nil || len(bz) != 8 {
		return 1
	}
	id := binary.BigEndian.Uint64(bz)
	if id == 0 {
		return 1
	}
	return id
}

func (k Keeper) NextOrderID(ctx sdk.Context) uint64 {
	id := k.GetNextOrderID(ctx)
	k.SetNextOrderID(ctx, id+1)
	return id
}

func (k Keeper) assertOrderInvariant(order *types.Order) error {
	amount, err := parsePositiveInt(order.Amount, "amount")
	if err != nil {
		return err
	}
	filled, err := parseNonNegativeInt(order.FilledAmount, "filled_amount")
	if err != nil {
		return err
	}
	remaining, err := parseNonNegativeInt(order.RemainingAmount, "remaining_amount")
	if err != nil {
		return err
	}
	if !filled.Add(remaining).Equal(amount) {
		return fmt.Errorf("filled_amount + remaining_amount must equal amount")
	}
	return nil
}

func (k Keeper) expireOrderIfNeeded(ctx sdk.Context, order *types.Order) error {
	if order.ExpireBh <= 0 {
		return nil
	}
	if order.Status != types.OrderStatus_ORDER_STATUS_OPEN && order.Status != types.OrderStatus_ORDER_STATUS_PARTIALLY_FILLED {
		return nil
	}
	if ctx.BlockHeight() >= order.ExpireBh {
		prevStatus := order.Status
		order.Status = types.OrderStatus_ORDER_STATUS_EXPIRED
		if err := k.onOrderStatusTransition(ctx, order, prevStatus); err != nil {
			return err
		}
		k.SetOrder(ctx, order)
	}
	return nil
}

func (k Keeper) fillOrder(order *types.Order, matchAmount sdkmath.Int) error {
	filled, err := parseNonNegativeInt(order.FilledAmount, "filled_amount")
	if err != nil {
		return err
	}
	remaining, err := parseNonNegativeInt(order.RemainingAmount, "remaining_amount")
	if err != nil {
		return err
	}

	if remaining.LT(matchAmount) {
		return fmt.Errorf("match_amount exceeds order remaining_amount")
	}

	filled = filled.Add(matchAmount)
	remaining = remaining.Sub(matchAmount)

	order.FilledAmount = filled.String()
	order.RemainingAmount = remaining.String()
	if remaining.IsZero() {
		order.Status = types.OrderStatus_ORDER_STATUS_FILLED
	} else {
		order.Status = types.OrderStatus_ORDER_STATUS_PARTIALLY_FILLED
	}

	return k.assertOrderInvariant(order)
}

func isBuySide(side types.OrderSide) bool {
	return side == types.OrderSide_ORDER_SIDE_BUY_YES || side == types.OrderSide_ORDER_SIDE_BUY_NO
}

func isSellSide(side types.OrderSide) bool {
	return side == types.OrderSide_ORDER_SIDE_SELL_YES || side == types.OrderSide_ORDER_SIDE_SELL_NO
}
