package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/x/paxi/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(k Keeper) types.MsgServer {
	return &msgServer{Keeper: k}
}

func (k msgServer) BurnToken(goCtx context.Context, msg *types.MsgBurnToken) (*types.MsgBurnTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	amount := sdk.NewCoins()
	for _, coin := range msg.Amount {
		amount = amount.Add(sdk.NewCoin(coin.Denom, coin.Amount))
	}

	err = k.Keeper.BurnFromUser(ctx, sender, amount)
	if err != nil {
		return nil, err
	}

	return &types.MsgBurnTokenResponse{}, nil
}
