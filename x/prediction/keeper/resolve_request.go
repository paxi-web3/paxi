package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/x/prediction/types"
)

func (k Keeper) GetResolveRequest(ctx sdk.Context, marketID uint64) (*types.ResolveRequest, bool) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.ResolveRequestStoreKey(marketID))
	if err != nil || bz == nil {
		return nil, false
	}

	req := &types.ResolveRequest{}
	if err := k.cdc.Unmarshal(bz, req); err != nil {
		return nil, false
	}
	return req, true
}

func (k Keeper) SetResolveRequest(ctx sdk.Context, req *types.ResolveRequest) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := k.cdc.Marshal(req)
	if err != nil {
		panic(err)
	}
	k.mustSet(store, types.ResolveRequestStoreKey(req.MarketId), bz)
}
