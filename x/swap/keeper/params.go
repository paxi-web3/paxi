package keeper

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/paxi-web3/paxi/x/swap/types"
)

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	if err := store.Set(types.KeyParams, bz); err != nil {
		panic(err)
	}
}

// GetParams retrieves the current params from store
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	params := types.Params{}

	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.KeyParams)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bz, &params)
	if err != nil {
		panic(err)
	}

	return params

}
