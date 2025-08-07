package types

import (
	"fmt"
	"math"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgSwap{}

// ValidateBasic performs stateless checks on MsgSwap fields.
// It uses fmt.Errorf exclusively for error reporting and adds overflow protection.
func (msg *MsgSwap) ValidateBasic() error {
	// Validate the creator address is a valid Bech32 address
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}
	// Validate the PRC20 contract address is a valid Bech32 address
	if _, err := sdk.AccAddressFromBech32(msg.Prc20); err != nil {
		return fmt.Errorf("invalid prc20 address: %w", err)
	}
	// OfferDenom must not be empty
	if msg.OfferDenom == "" {
		return fmt.Errorf("offer_denom cannot be empty")
	}

	// Parse and validate the offer amount string
	offerAmt, ok := sdkmath.NewIntFromString(msg.OfferAmount)
	if !ok {
		return fmt.Errorf("invalid offer amount: must be an integer string")
	}
	// Offer amount must be strictly positive
	if !offerAmt.IsPositive() {
		return fmt.Errorf("offer amount must be positive, got %s", offerAmt.String())
	}
	// Prevent panic when converting to uint64
	if offerAmt.GT(sdkmath.NewIntFromUint64(math.MaxUint64)) {
		return fmt.Errorf("offer amount too large: must fit in uint64")
	}

	// Parse and validate the minimum receive amount string
	minRecv, ok := sdkmath.NewIntFromString(msg.MinReceive)
	if !ok {
		return fmt.Errorf("invalid min_receive: must be an integer string")
	}
	// MinReceive must be non-negative
	if minRecv.IsNegative() {
		return fmt.Errorf("min_receive cannot be negative, got %s", minRecv.String())
	}
	// Prevent panic when converting to uint64
	if minRecv.GT(sdkmath.NewIntFromUint64(math.MaxUint64)) {
		return fmt.Errorf("min_receive too large: must fit in uint64")
	}

	return nil
}
