package customstaking

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

type MsgServer struct {
	msgServer types.MsgServer
}

func NewMsgServer(msgServer types.MsgServer) types.MsgServer {
	return &MsgServer{msgServer: msgServer}
}

func (m *MsgServer) Delegate(goCtx context.Context, msg *types.MsgDelegate) (*types.MsgDelegateResponse, error) {
	// Custom logic to enforce minimum delegation amount
	if msg.Amount.Amount.LT(sdkmath.NewInt(MinDelegation)) {
		return nil, fmt.Errorf("minimum delegation is %f PAXI", float64(MinDelegation)/float64(1_000_000))
	}
	return m.msgServer.Delegate(goCtx, msg)
}

func (m *MsgServer) Undelegate(goCtx context.Context, msg *types.MsgUndelegate) (*types.MsgUndelegateResponse, error) {
	// Custom logic to enforce minimum undelegation amount
	if msg.Amount.Amount.LT(sdkmath.NewInt(MinUndelegation)) {
		return nil, fmt.Errorf("minimum undelegation is %f PAXI", float64(MinUndelegation)/float64(1_000_000))
	}
	return m.msgServer.Undelegate(goCtx, msg)
}

func (m *MsgServer) BeginRedelegate(goCtx context.Context, msg *types.MsgBeginRedelegate) (*types.MsgBeginRedelegateResponse, error) {
	return m.msgServer.BeginRedelegate(goCtx, msg)
}

func (m *MsgServer) EditValidator(goCtx context.Context, msg *types.MsgEditValidator) (*types.MsgEditValidatorResponse, error) {
	return m.msgServer.EditValidator(goCtx, msg)
}

func (m *MsgServer) CreateValidator(goCtx context.Context, msg *types.MsgCreateValidator) (*types.MsgCreateValidatorResponse, error) {
	return m.msgServer.CreateValidator(goCtx, msg)
}

func (m *MsgServer) CancelUnbondingDelegation(ctx context.Context, msg *types.MsgCancelUnbondingDelegation) (*types.MsgCancelUnbondingDelegationResponse, error) {
	return m.msgServer.CancelUnbondingDelegation(ctx, msg)
}

func (m *MsgServer) UpdateParams(ctx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	return m.msgServer.UpdateParams(ctx, msg)
}
