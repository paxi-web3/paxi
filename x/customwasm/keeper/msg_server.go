package keeper

import (
	"context"
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/paxi-web3/paxi/x/customwasm/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(k Keeper) types.MsgServer {
	return &msgServer{Keeper: k}
}

func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.Authority != k.authority {
		return nil, sdkerrors.ErrUnauthorized
	}

	parsedParams := types.Params{
		StoreCodeBaseGas:    msg.Params.StoreCodeBaseGas,
		StoreCodeMultiplier: msg.Params.StoreCodeMultiplier,
		InstBaseGas:         msg.Params.InstBaseGas,
		InstMultiplier:      msg.Params.InstMultiplier,
	}

	if err := parsedParams.Validate(); err != nil {
		return nil, err
	}

	store := k.storeService.OpenKVStore(ctx)
	bz, err := json.Marshal(&parsedParams)
	if err != nil {
		return nil, err
	}
	if err := store.Set(types.KeyParams, bz); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
