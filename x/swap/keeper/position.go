package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/x/swap/types"
)

func (k Keeper) GetPosition(ctx sdk.Context, prc20 string, creator sdk.AccAddress) (types.ProviderPosition, bool) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.PositionKey(prc20, creator)

	bz, err := store.Get(key)
	if err != nil || bz == nil {
		return types.ProviderPosition{}, false
	}

	var pos types.ProviderPosition
	if err := k.cdc.Unmarshal(bz, &pos); err != nil {
		return types.ProviderPosition{}, false
	}
	return pos, true
}

func (k Keeper) SetPosition(ctx sdk.Context, pos types.ProviderPosition) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.PositionKey(pos.Prc20, sdk.MustAccAddressFromBech32(pos.Creator))

	bz, err := k.cdc.Marshal(&pos)
	if err != nil {
		panic(fmt.Errorf("failed to marshal position: %w", err))
	}

	if err := store.Set(key, bz); err != nil {
		panic(fmt.Errorf("failed to store position: %w", err))
	}
}

func (k Keeper) DeletePosition(ctx sdk.Context, prc20 string, creator sdk.AccAddress) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.PositionKey(prc20, creator)
	if err := store.Delete(key); err != nil {
		panic(fmt.Errorf("failed to delete position: %w", err))
	}
}
