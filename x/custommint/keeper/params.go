package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/paxi-web3/paxi/x/custommint/types"
)

var ParamsKey = []byte("custommint_params")

// GetParams returns the current x/slashing module parameters.
func (k Keeper) GetParams(ctx context.Context) (params types.GenesisState, err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := k.storeService.OpenKVStore(sdkCtx)

	bz, err := store.Get(ParamsKey)

	if err != nil {
		return params, err
	}
	if bz == nil {
		return params, nil
	}

	err = k.cdc.Unmarshal(bz, &params)
	return params, err
}

// SetParams sets the x/slashing module parameters.
// CONTRACT: This method performs no validation of the parameters.
func (k Keeper) SetParams(ctx context.Context, params *types.GenesisState) error {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := k.cdc.Marshal(params)
	if err != nil {
		return err
	}
	return store.Set(ParamsKey, bz)
}
