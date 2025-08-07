package types

import (
	"fmt"
)

const (
	ModuleName   = "swap"
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	DefaultDenom = "upaxi"
)

var (
	PoolPrefix     = []byte{0x01} // Pools
	PositionPrefix = []byte{0x02} // Provider Positions
	KeyParams      = []byte("swap_params")
)

// Params defines spam protection parameters for the Paxi blockchain.
type Params struct {
	CodeID       uint64 `json:"code_id" yaml:"code_id"`
	SwapFeeBPS   uint64 `json:"swap_fee_bps" yaml:"swap_fee_bps"`
	MinLiquidity uint64 `json:"min_liquidity" yaml:"min_liquidity"`
}

type GenesisState struct {
	Params Params `json:"params" yaml:"params"`
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

func (gs GenesisState) Validate() error {
	return gs.Params.Validate()
}

const (
	// The maximum allowed MinLiquidity can be adjusted according to business needs to prevent overflow or logical anomalies.
	MaxMinLiquidity = ^uint64(0) / 2
)

func DefaultParams() Params {
	return Params{
		CodeID:       1,         // Default code ID for prc-20 contracts
		SwapFeeBPS:   4,         // 0.4% note: 1000 BPS = 100%
		MinLiquidity: 1_000_000, // 1 Paxi
	}
}

// Validate checks all Params fields for acceptable values.
// Returns an error if any parameter is out-of-bounds, using fmt.Errorf.
func (p Params) Validate() error {
	// Ensure swap fee is within [1, 10000] BPS
	if p.SwapFeeBPS < 1 || p.SwapFeeBPS > 10_000 {
		return fmt.Errorf("swap fee BPS must be between 1 and 10000 (got %d)", p.SwapFeeBPS)
	}

	// Ensure CodeID is non-zero
	if p.CodeID == 0 {
		return fmt.Errorf("wasm CodeID must be greater than 0")
	}

	// Ensure MinLiquidity is positive
	if p.MinLiquidity == 0 {
		return fmt.Errorf("min liquidity must be greater than 0")
	}

	// Prevent unreasonably large MinLiquidity values
	if p.MinLiquidity > MaxMinLiquidity {
		return fmt.Errorf("min liquidity too large (got %d, max %d)", p.MinLiquidity, MaxMinLiquidity)
	}

	return nil
}

func LPTokenDenom(prc20 string) string {
	return fmt.Sprintf("lp/%s", prc20)
}
