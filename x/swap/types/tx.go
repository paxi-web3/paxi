package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgProvideLiquidity{}

func (msg *MsgProvideLiquidity) Route() string {
	return RouterKey
}

func (msg *MsgProvideLiquidity) Type() string {
	return "ProvideLiquidity"
}

func (msg *MsgProvideLiquidity) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(fmt.Sprintf("invalid creator address: %v", err))
	}
	return []sdk.AccAddress{addr}
}

func (msg *MsgProvideLiquidity) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}
	if msg.PaxiAmount.Denom != DefaultDenom {
		return fmt.Errorf("invalid denom: expected %s", DefaultDenom)
	}
	if !msg.PaxiAmount.Amount.IsPositive() {
		return fmt.Errorf("paxi amount must be positive")
	}
	if amt, ok := sdkmath.NewIntFromString(msg.Prc20Amount); !ok || !amt.IsPositive() {
		return fmt.Errorf("invalid prc20 amount: must be positive integer string")
	}
	return nil
}
