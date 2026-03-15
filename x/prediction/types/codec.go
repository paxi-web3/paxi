package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterMsg(cdc codectypes.InterfaceRegistry) {
	cdc.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgUpdateParams{},
		&MsgCreateMarket{},
		&MsgPlaceOrder{},
		&MsgCancelOrder{},
		&MsgSplitPosition{},
		&MsgMergePosition{},
		&MsgApplyTradeBatch{},
		&MsgResolveMarket{},
		&MsgVoidMarket{},
		&MsgClaimPayout{},
		&MsgClaimVoidRefund{},
	)

	msgservice.RegisterMsgServiceDesc(cdc, &_Msg_serviceDesc)
}
