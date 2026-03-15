package keeper

import (
	"encoding/binary"
	"fmt"
	"math"

	store "cosmossdk.io/core/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/x/prediction/types"
)

func isOpenOrderStatus(status types.OrderStatus) bool {
	return status == types.OrderStatus_ORDER_STATUS_OPEN || status == types.OrderStatus_ORDER_STATUS_PARTIALLY_FILLED
}

func (k Keeper) getOpenOrderCountByUserRaw(store store.KVStore, addr sdk.AccAddress) (uint64, bool) {
	bz, err := store.Get(types.OpenOrderCountByUserKey(addr))
	if err != nil || len(bz) != 8 {
		return 0, false
	}
	return binary.BigEndian.Uint64(bz), true
}

func (k Keeper) getOpenOrderCountByUserMarketRaw(store store.KVStore, addr sdk.AccAddress, marketID uint64) (uint64, bool) {
	bz, err := store.Get(types.OpenOrderCountByUserMarketKey(addr, marketID))
	if err != nil || len(bz) != 8 {
		return 0, false
	}
	return binary.BigEndian.Uint64(bz), true
}

func (k Keeper) setOpenOrderCountByUser(store store.KVStore, addr sdk.AccAddress, count uint64) {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	k.mustSet(store, types.OpenOrderCountByUserKey(addr), bz)
}

func (k Keeper) setOpenOrderCountByUserMarket(store store.KVStore, addr sdk.AccAddress, marketID uint64, count uint64) {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	k.mustSet(store, types.OpenOrderCountByUserMarketKey(addr, marketID), bz)
}

func (k Keeper) isOrderEffectivelyOpen(ctx sdk.Context, order *types.Order) bool {
	if !isOpenOrderStatus(order.Status) {
		return false
	}
	if order.ExpireBh > 0 && ctx.BlockHeight() >= order.ExpireBh {
		return false
	}
	return true
}

func (k Keeper) scanOpenOrderCountsForUser(ctx sdk.Context, addr sdk.AccAddress, marketID uint64) (uint64, uint64) {
	store := k.storeService.OpenKVStore(ctx)
	it, err := prefixIterator(store, types.OrderPrefix)
	if err != nil {
		panic(err)
	}
	defer it.Close()

	var userCount uint64
	var marketCount uint64
	for ; it.Valid(); it.Next() {
		order := &types.Order{}
		if err := k.cdc.Unmarshal(it.Value(), order); err != nil {
			continue
		}
		if order.Trader != addr.String() {
			continue
		}
		if !k.isOrderEffectivelyOpen(ctx, order) {
			continue
		}
		userCount++
		if order.MarketId == marketID {
			marketCount++
		}
	}

	return userCount, marketCount
}

func (k Keeper) getOrRebuildOpenOrderCounts(ctx sdk.Context, addr sdk.AccAddress, marketID uint64) (uint64, uint64) {
	store := k.storeService.OpenKVStore(ctx)
	userCount, userFound := k.getOpenOrderCountByUserRaw(store, addr)
	marketCount, marketFound := k.getOpenOrderCountByUserMarketRaw(store, addr, marketID)
	if userFound && marketFound {
		return userCount, marketCount
	}

	userCount, marketCount = k.scanOpenOrderCountsForUser(ctx, addr, marketID)
	k.setOpenOrderCountByUser(store, addr, userCount)
	k.setOpenOrderCountByUserMarket(store, addr, marketID, marketCount)
	return userCount, marketCount
}

func (k Keeper) incrementOpenOrderCounts(ctx sdk.Context, addr sdk.AccAddress, marketID uint64) error {
	store := k.storeService.OpenKVStore(ctx)
	userCount, marketCount := k.getOrRebuildOpenOrderCounts(ctx, addr, marketID)
	if userCount == math.MaxUint64 || marketCount == math.MaxUint64 {
		return fmt.Errorf("open order counter overflow")
	}
	k.setOpenOrderCountByUser(store, addr, userCount+1)
	k.setOpenOrderCountByUserMarket(store, addr, marketID, marketCount+1)
	return nil
}

func (k Keeper) decrementOpenOrderCounts(ctx sdk.Context, addr sdk.AccAddress, marketID uint64) error {
	store := k.storeService.OpenKVStore(ctx)
	userCount, marketCount := k.getOrRebuildOpenOrderCounts(ctx, addr, marketID)
	if userCount == 0 || marketCount == 0 {
		return fmt.Errorf("open order counter underflow")
	}
	k.setOpenOrderCountByUser(store, addr, userCount-1)
	k.setOpenOrderCountByUserMarket(store, addr, marketID, marketCount-1)
	return nil
}

func (k Keeper) enforceOpenOrderLimit(ctx sdk.Context, addr sdk.AccAddress, marketID uint64, params types.Params) error {
	userCount, marketCount := k.getOrRebuildOpenOrderCounts(ctx, addr, marketID)
	if userCount >= params.MaxOpenOrdersPerUser || marketCount >= params.MaxOpenOrdersPerMarket {
		// Re-scan to avoid false positives from stale counters (e.g. orders expired by height).
		store := k.storeService.OpenKVStore(ctx)
		userCount, marketCount = k.scanOpenOrderCountsForUser(ctx, addr, marketID)
		k.setOpenOrderCountByUser(store, addr, userCount)
		k.setOpenOrderCountByUserMarket(store, addr, marketID, marketCount)
	}
	if userCount >= params.MaxOpenOrdersPerUser {
		return fmt.Errorf("max_open_orders_per_user exceeded")
	}
	if marketCount >= params.MaxOpenOrdersPerMarket {
		return fmt.Errorf("max_open_orders_per_market exceeded")
	}
	return nil
}

func (k Keeper) onOrderStatusTransition(ctx sdk.Context, order *types.Order, prevStatus types.OrderStatus) error {
	prevOpen := isOpenOrderStatus(prevStatus)
	nextOpen := isOpenOrderStatus(order.Status)

	if prevOpen && !nextOpen {
		order.ClosedBh = ctx.BlockHeight()
	} else if !prevOpen && nextOpen {
		order.ClosedBh = 0
	}

	if prevOpen == nextOpen {
		return nil
	}

	addr, err := sdk.AccAddressFromBech32(order.Trader)
	if err != nil {
		return err
	}

	if nextOpen {
		return k.incrementOpenOrderCounts(ctx, addr, order.MarketId)
	}
	return k.decrementOpenOrderCounts(ctx, addr, order.MarketId)
}
