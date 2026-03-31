package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/x/prediction/types"
)

func isOpenOrderStatus(status types.OrderStatus) bool {
	return status == types.OrderStatus_ORDER_STATUS_OPEN || status == types.OrderStatus_ORDER_STATUS_PARTIALLY_FILLED
}

func (k Keeper) onOrderStatusTransition(ctx sdk.Context, order *types.Order, prevStatus types.OrderStatus) error {
	prevOpen := isOpenOrderStatus(prevStatus)
	nextOpen := isOpenOrderStatus(order.Status)

	if prevOpen && !nextOpen {
		order.ClosedBh = ctx.BlockHeight()
	} else if !prevOpen && nextOpen {
		order.ClosedBh = 0
	}

	return nil
}
