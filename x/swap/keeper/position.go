package keeper

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/x/swap/types"
)

func (k Keeper) GetPosition(ctx sdk.Context, prc20 string, creator sdk.AccAddress) (types.ProviderPosition, bool) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.PositionKey(prc20, creator)
	bz, err := store.Get(key)
	if err != nil {
		return types.ProviderPosition{}, false
	}
	var pos types.ProviderPosition
	_ = json.Unmarshal(bz, &pos)
	return pos, true
}

func (k Keeper) SetPosition(ctx sdk.Context, pos types.ProviderPosition) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.PositionKey(pos.Prc20, sdk.MustAccAddressFromBech32(pos.Creator))
	bz, _ := json.Marshal(pos)
	_ = store.Set(key, bz)
}

func (k Keeper) DeletePosition(ctx sdk.Context, prc20 string, creator sdk.AccAddress) {
	store := k.storeService.OpenKVStore(ctx)
	key := []byte(fmt.Sprintf("position:%s:%s", prc20, creator.String()))
	_ = store.Delete(key)
}
