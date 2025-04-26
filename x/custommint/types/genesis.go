package types

import (
	"fmt"
)

// DefaultGenesisState returns the default genesis state for the custommint module.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		MintDenom:           DefaultDenom,
		BlocksPerYear:       BlocksPerYear,
		FirstYearInflation:  FirstYearInflation,
		SecondYearInflation: SecondYearInflation,
		OtherYearInflation:  OtherYearInflation,
	}
}

// Validate validates the custommint genesis state.
func (gs *GenesisState) Validate() error {
	if gs.MintDenom == "" {
		return fmt.Errorf("MintDenom cannot be empty")
	}
	if gs.BlocksPerYear <= 0 {
		return fmt.Errorf("BlocksPerYear must be positive")
	}
	if gs.FirstYearInflation < 0.0 {
		return fmt.Errorf("FirstYearInflation must be positive")
	}
	if gs.SecondYearInflation < 0.0 {
		return fmt.Errorf("SecondYearInflation must be positive")
	}
	if gs.OtherYearInflation < 0.0 {
		return fmt.Errorf("OtherYearInflation must be positive")
	}
	return nil
}
