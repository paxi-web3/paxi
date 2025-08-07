package types

import (
	"fmt"
	"math"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgWithdrawLiquidity{}

// ValidateBasic performs stateless checks on MsgWithdrawLiquidity fields.
// It uses fmt.Errorf exclusively for error reporting and adds overflow protection.
func (msg *MsgWithdrawLiquidity) ValidateBasic() error {
	// Validate that Creator is a valid Bech32 address
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}

	// Validate that Prc20 is a valid Bech32 address
	if _, err := sdk.AccAddressFromBech32(msg.Prc20); err != nil {
		return fmt.Errorf("invalid prc20 address: %w", err)
	}

	// Parse LP amount string into sdk.Int
	lpAmt, ok := sdkmath.NewIntFromString(msg.LpAmount)
	if !ok {
		return fmt.Errorf("invalid lp amount: must be an integer string")
	}
	// Ensure LP amount is strictly positive
	if !lpAmt.IsPositive() {
		return fmt.Errorf("lp amount must be positive, got %s", lpAmt.String())
	}

	// Prevent panic when converting to uint64
	if lpAmt.GT(sdkmath.NewIntFromUint64(math.MaxUint64)) {
		return fmt.Errorf("lp amount too large: must fit in uint64")
	}

	return nil
}
