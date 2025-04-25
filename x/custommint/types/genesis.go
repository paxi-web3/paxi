package types

import (
	"fmt"
)

// DefaultGenesisState returns the default genesis state for the custommint module.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		MintDenom:     DefaultDenom,
		TotalSupply:   TotalSupply,
		BlocksPerYear: BlocksPerYear,
	}
}

// Validate validates the custommint genesis state.
func (gs *GenesisState) Validate() error {
	if gs.MintDenom == "" {
		return fmt.Errorf("mint denom cannot be empty")
	}
	if gs.TotalSupply <= 0 {
		return fmt.Errorf("total supply must be positive")
	}
	if gs.BlocksPerYear <= 0 {
		return fmt.Errorf("blocks per year must be positive")
	}
	return nil
}
