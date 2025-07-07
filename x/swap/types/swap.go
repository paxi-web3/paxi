package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgSwap{}

func (msg *MsgSwap) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}
	if _, err := sdk.AccAddressFromBech32(msg.Prc20); err != nil {
		return fmt.Errorf("invalid prc20 address: %w", err)
	}
	if msg.OfferDenom == "" {
		return fmt.Errorf("offer_denom cannot be empty")
	}
	if amt, ok := sdkmath.NewIntFromString(msg.OfferAmount); !ok || !amt.IsPositive() {
		return fmt.Errorf("invalid offer amount")
	}
	if minRecv, ok := sdkmath.NewIntFromString(msg.MinReceive); !ok || minRecv.IsNegative() {
		return fmt.Errorf("invalid min_receive")
	}
	return nil
}
