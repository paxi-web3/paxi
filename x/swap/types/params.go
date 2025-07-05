package types

import (
	"fmt"
)

const (
	ModuleName   = "swap"
	StoreKey     = ModuleName
	DefaultDenom = "upaxi"
)

var KeyParams = []byte("swap_params")

// Params defines spam protection parameters for the Paxi blockchain.
type Params struct {
	CodeID     uint64 `json:"code_id" yaml:"code_id"`
	SwapFeeBPS uint64 `json:"swap_fee_bps" yaml:"swap_fee_bps"`
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

func DefaultParams() Params {
	return Params{
		CodeID:     1,  // Default code ID for prc-20 contracts
		SwapFeeBPS: 40, // 0.4%
	}
}

func (p Params) Validate() error {
	if p.SwapFeeBPS <= 0 || p.SwapFeeBPS > 10000 {
		return fmt.Errorf("commission rate must be between 0 and 10000")
	}
	if p.CodeID == 0 {
		return fmt.Errorf("code ID must be greater than 0")
	}
	return nil
}
