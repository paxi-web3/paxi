package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterMsg(cdc codectypes.InterfaceRegistry) {
	cdc.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgBurnToken{},
		&MsgUpdateParams{},
	)
}
