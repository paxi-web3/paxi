package keeper

import (
	"context"
	"encoding/json"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/paxi-web3/paxi/x/custommint/types"
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

	// Parse and validate BurnThreshold
	burnThreshold, ok := sdkmath.NewIntFromString(msg.Params.BurnThreshold)
	if !ok {
		return nil, fmt.Errorf("invalid BurnThreshold: %q", msg.Params.BurnThreshold)
	}
	if burnThreshold.IsNegative() {
		return nil, fmt.Errorf("BurnThreshold must be non-negative")
	}

	parsedParams := types.Params{
		BurnThreshold:       burnThreshold,
		BurnRatio:           sdkmath.LegacyMustNewDecFromStr(msg.Params.BurnRatio),
		BlocksPerYear:       msg.Params.BlocksPerYear,
		FirstYearInflation:  sdkmath.LegacyMustNewDecFromStr(msg.Params.FirstYearInflation),
		SecondYearInflation: sdkmath.LegacyMustNewDecFromStr(msg.Params.SecondYearInflation),
		OtherYearInflation:  sdkmath.LegacyMustNewDecFromStr(msg.Params.OtherYearInflation),
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
