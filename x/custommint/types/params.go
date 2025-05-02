package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
)

const (
	ModuleName   = "custommint"
	StoreKey     = ModuleName
	TotalMinted  = "total_minted"
	DefaultDenom = "upaxi"
)

var KeyParams = []byte("custommint_params")

type Params struct {
	BurnThreshold       sdkmath.Int       `json:"burn_threshold" yaml:"burn_threshold"`
	BurnRatio           sdkmath.LegacyDec `json:"burn_ratio" yaml:"burn_ratio"`
	BlocksPerYear       int64             `json:"blocks_per_year" yaml:"blocks_per_year"`
	FirstYearInflation  sdkmath.LegacyDec `json:"first_year_inflation" yaml:"first_year_inflation"`
	SecondYearInflation sdkmath.LegacyDec `json:"second_year_inflation" yaml:"second_year_inflation"`
	OtherYearInflation  sdkmath.LegacyDec `json:"other_year_inflation" yaml:"other_year_inflation"`
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
		BurnThreshold:       sdkmath.NewInt(100_000),
		BurnRatio:           sdkmath.LegacyNewDecWithPrec(50, 2), // 0.50 (50%)
		BlocksPerYear:       7884000,                             // 365 * 24 * 60 * 60 / 4s block time
		FirstYearInflation:  sdkmath.LegacyMustNewDecFromStr("0.08"),
		SecondYearInflation: sdkmath.LegacyMustNewDecFromStr("0.04"),
		OtherYearInflation:  sdkmath.LegacyMustNewDecFromStr("0.02"),
	}
}

func (p Params) Validate() error {
	if p.BurnThreshold.IsNegative() {
		return fmt.Errorf("burn threshold cannot be negative")
	}
	if p.BurnRatio.IsNegative() || p.BurnRatio.GT(sdkmath.LegacyOneDec()) {
		return fmt.Errorf("burn ratio must be between 0 and 1")
	}
	if p.BlocksPerYear <= 0 {
		return fmt.Errorf("blocks per year must be positive")
	}
	if p.FirstYearInflation.IsNegative() || p.SecondYearInflation.IsNegative() || p.OtherYearInflation.IsNegative() {
		return fmt.Errorf("inflation values must be non-negative")
	}
	return nil
}
