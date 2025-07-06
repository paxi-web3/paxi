package types

import (
	sdkmath "cosmossdk.io/math"
)

type Pool struct {
	Prc20        string      `json:"prc20" yaml:"prc20"`                 // PRC20 contract address (Bech32)
	ReservePaxi  sdkmath.Int `json:"reserve_paxi" yaml:"reserve_paxi"`   // PAXI reserve amount
	ReservePRC20 sdkmath.Int `json:"reserve_prc20" yaml:"reserve_prc20"` // PRC20 reserve amount
}
