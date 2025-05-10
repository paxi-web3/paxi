package types

import "fmt"

const (
	ModuleName = "customwasm"
	StoreKey   = ModuleName
)

var KeyParams = []byte("customwasm_params")

// Params defines gas cost configuration for storing and instantiating WASM contracts.
type Params struct {
	StoreCodeBaseGas    uint64 `json:"store_code_base_gas" yaml:"store_code_base_gas"`
	StoreCodeMultiplier uint64 `json:"store_code_multiplier" yaml:"store_code_multiplier"`
	InstBaseGas         uint64 `json:"inst_base_gas" yaml:"inst_base_gas"`
	InstMultiplier      uint64 `json:"inst_multiplier" yaml:"inst_multiplier"`
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
		StoreCodeBaseGas:    uint64(240_000_000),
		StoreCodeMultiplier: uint64(300),
		InstBaseGas:         uint64(30_000_000),
		InstMultiplier:      uint64(100),
	}
}

func (p Params) Validate() error {
	if p.StoreCodeBaseGas <= 0 || p.StoreCodeMultiplier <= 0 || p.InstBaseGas <= 0 || p.InstMultiplier <= 0 {
		return fmt.Errorf("all params from customwasm can't be less or equal than 0")
	}
	return nil
}
