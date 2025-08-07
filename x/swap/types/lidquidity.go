package types

import (
	"fmt"
	"math"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Pool struct {
	Prc20        string      `json:"prc20" yaml:"prc20"`                 // PRC20 contract address (Bech32)
	ReservePaxi  sdkmath.Int `json:"reserve_paxi" yaml:"reserve_paxi"`   // PAXI reserve amount
	ReservePRC20 sdkmath.Int `json:"reserve_prc20" yaml:"reserve_prc20"` // PRC20 reserve amount
	TotalShares  sdkmath.Int `json:"total_shares" yaml:"total_shares"`   // Total LP shares issued
}

var _ sdk.Msg = &MsgProvideLiquidity{}

// ValidateBasic performs stateless checks on MsgProvideLiquidity fields.
// It uses fmt.Errorf exclusively for error reporting.
func (msg *MsgProvideLiquidity) ValidateBasic() error {
	// Validate the creator address is a valid Bech32 address
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}

	// Parse and validate the Paxi coin string
	paxiAmount, err := sdk.ParseCoinNormalized(msg.PaxiAmount)
	if err != nil {
		return fmt.Errorf("invalid paxi amount: %w", err)
	}
	// Ensure the denom matches the default Paxi denom
	if paxiAmount.Denom != DefaultDenom {
		return fmt.Errorf("invalid denom: expected %s, got %s", DefaultDenom, paxiAmount.Denom)
	}
	// Amount must be strictly positive
	if !paxiAmount.Amount.IsPositive() {
		return fmt.Errorf("paxi amount must be positive, got %s", paxiAmount.Amount.String())
	}

	// Parse and validate the PRC20 amount string
	prc20Amt, ok := sdkmath.NewIntFromString(msg.Prc20Amount)
	if !ok {
		return fmt.Errorf("invalid prc20 amount: must be an integer string")
	}
	// Amount must be strictly positive
	if !prc20Amt.IsPositive() {
		return fmt.Errorf("prc20 amount must be positive, got %s", prc20Amt.String())
	}

	// Prevent panic when converting to uint64
	if paxiAmount.Amount.GT(sdkmath.NewIntFromUint64(math.MaxUint64)) {
		return fmt.Errorf("paxi amount too large: must fit in uint64")
	}
	// Prevent exceeding the maximum total supply
	maxSupply := sdkmath.NewInt(1_000_000_000_000_000) // e.g., 1e15 units
	if paxiAmount.Amount.GT(maxSupply) {
		return fmt.Errorf("paxi amount exceeds maximum supply (%s)", maxSupply.String())
	}

	// Similarly, enforce upper bound check for prc20Amt
	if prc20Amt.GT(sdkmath.NewIntFromUint64(math.MaxUint64)) {
		return fmt.Errorf("prc20 amount too large: must fit in uint64")
	}

	return nil
}

func PoolStoreKey(prc20 string) []byte {
	return append(PoolPrefix, []byte(prc20)...)
}

// ToProto converts internal Pool to proto.Pool
func (p Pool) ToProto() PoolProto {
	return PoolProto{
		Prc20:        p.Prc20,
		ReservePaxi:  p.ReservePaxi.String(),
		ReservePrc20: p.ReservePRC20.String(),
		TotalShares:  p.TotalShares.String(),
	}
}

// PoolFromProto converts proto.Pool to internal Pool
func PoolFromProto(pp *PoolProto) (Pool, error) {
	reservePaxi, ok := sdkmath.NewIntFromString(pp.ReservePaxi)
	if !ok {
		return Pool{}, fmt.Errorf("invalid ReservePaxi: %s", pp.ReservePaxi)
	}

	reservePRC20, ok := sdkmath.NewIntFromString(pp.ReservePrc20)
	if !ok {
		return Pool{}, fmt.Errorf("invalid ReservePRC20: %s", pp.ReservePrc20)
	}

	totalShares, ok := sdkmath.NewIntFromString(pp.TotalShares)
	if !ok {
		return Pool{}, fmt.Errorf("invalid TotalShares: %s", pp.TotalShares)
	}

	return Pool{
		Prc20:        pp.Prc20,
		ReservePaxi:  reservePaxi,
		ReservePRC20: reservePRC20,
		TotalShares:  totalShares,
	}, nil
}
