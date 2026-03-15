package keeper

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/x/prediction/types"
)

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := json.Marshal(&params)
	if err != nil {
		panic(err)
	}
	if err := store.Set(types.KeyParams, bz); err != nil {
		panic(err)
	}
}

func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	params := types.DefaultParams()
	store := k.storeService.OpenKVStore(ctx)

	bz, err := store.Get(types.KeyParams)
	if err != nil || bz == nil {
		return params
	}

	if err := json.Unmarshal(bz, &params); err != nil {
		return types.DefaultParams()
	}

	return params
}
