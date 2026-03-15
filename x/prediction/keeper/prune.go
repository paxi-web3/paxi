package keeper

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/utils"
	"github.com/paxi-web3/paxi/x/prediction/types"
)

func isTerminalOrderStatus(status types.OrderStatus) bool {
	return status == types.OrderStatus_ORDER_STATUS_FILLED ||
		status == types.OrderStatus_ORDER_STATUS_CANCELLED ||
		status == types.OrderStatus_ORDER_STATUS_EXPIRED
}

func pruneAnchorHeight(order *types.Order) int64 {
	if order.ClosedBh > 0 {
		return order.ClosedBh
	}
	if order.Status == types.OrderStatus_ORDER_STATUS_EXPIRED && order.ExpireBh > 0 {
		return order.ExpireBh
	}
	return order.CreatedBh
}

func (k Keeper) shouldPruneOrder(order *types.Order, thresholdBh int64) bool {
	if !isTerminalOrderStatus(order.Status) {
		return false
	}
	filled, err := parseNonNegativeInt(order.FilledAmount, "filled_amount")
	if err != nil {
		return false
	}
	// Keep all orders that have ever been traded (filled_amount > 0) for auditability.
	if !filled.IsZero() {
		return false
	}
	anchor := pruneAnchorHeight(order)
	if anchor <= 0 {
		return false
	}
	return anchor <= thresholdBh
}

func (k Keeper) AutoPruneOrders(ctx sdk.Context) {
	params := k.GetParams(ctx)
	if ctx.BlockHeight() <= 0 {
		return
	}
	currentBh := uint64(ctx.BlockHeight())
	if currentBh%params.OrderPruneIntervalBh != 0 {
		return
	}
	if currentBh <= params.OrderPruneRetainBh {
		return
	}
	thresholdBh := int64(currentBh - params.OrderPruneRetainBh)

	store := k.storeService.OpenKVStore(ctx)
	cursor, err := store.Get(types.OrderPruneCursorKey)
	if err != nil {
		panic(err)
	}

	start, end := utils.PrefixRange(types.OrderPrefix)
	if len(cursor) > 0 {
		start = cursor
	}
	it, err := store.Iterator(start, end)
	if err != nil {
		panic(err)
	}
	defer it.Close()

	var scanned uint64
	var pruned uint64
	var lastScannedKey []byte
	orderIDsToDelete := make([][2]uint64, 0)

	for ; it.Valid() && scanned < params.OrderPruneScanLimit && pruned < params.OrderPruneDeleteLimit; it.Next() {
		key := append([]byte(nil), it.Key()...)
		if len(cursor) > 0 && bytes.Equal(key, cursor) {
			continue
		}

		lastScannedKey = key
		scanned++

		order := &types.Order{}
		if err := k.cdc.Unmarshal(it.Value(), order); err != nil {
			continue
		}
		if !k.shouldPruneOrder(order, thresholdBh) {
			continue
		}

		orderIDsToDelete = append(orderIDsToDelete, [2]uint64{order.MarketId, order.Id})
		pruned++
	}

	for i := range orderIDsToDelete {
		marketID := orderIDsToDelete[i][0]
		orderID := orderIDsToDelete[i][1]
		k.DeleteOrder(ctx, marketID, orderID)
	}

	if !it.Valid() {
		k.mustDelete(store, types.OrderPruneCursorKey)
	} else if len(lastScannedKey) > 0 {
		k.mustSet(store, types.OrderPruneCursorKey, lastScannedKey)
	}

	if scanned > 0 || pruned > 0 {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventOrdersPruned,
				sdk.NewAttribute(types.AttributeKeyPrunedOrders, intToStr(pruned)),
				sdk.NewAttribute(types.AttributeKeyScannedOrders, intToStr(scanned)),
			),
		)
	}
}
