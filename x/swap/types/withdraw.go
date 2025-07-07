package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgWithdrawLiquidity{}

func (msg *MsgWithdrawLiquidity) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}
	if _, err := sdk.AccAddressFromBech32(msg.Prc20); err != nil {
		return fmt.Errorf("invalid prc20 address: %w", err)
	}
	if amt, ok := sdkmath.NewIntFromString(msg.LpAmount); !ok || !amt.IsPositive() {
		return fmt.Errorf("invalid lp amount: must be positive integer string")
	}
	return nil
}
