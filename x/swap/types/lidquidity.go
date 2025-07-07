package types

import (
	"fmt"

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

func (msg *MsgProvideLiquidity) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}

	paxiAmount, err := sdk.ParseCoinNormalized(msg.PaxiAmount)
	if err != nil {
		return fmt.Errorf("invalid paxi amount: %w", err)
	}
	if paxiAmount.Denom != DefaultDenom {
		return fmt.Errorf("invalid denom: expected %s", DefaultDenom)
	}
	if !paxiAmount.Amount.IsPositive() {
		return fmt.Errorf("paxi amount must be positive")
	}

	if amt, ok := sdkmath.NewIntFromString(msg.Prc20Amount); !ok || !amt.IsPositive() {
		return fmt.Errorf("invalid prc20 amount: must be positive integer string")
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
