package types

import fmt "fmt"

const (
	ModuleName           = "paxi"
	StoreKey             = ModuleName
	LockedVestingKey     = "locked_vesting"
	VestingAccountPrefix = "vesting_account/"
	DefaultDenom         = "upaxi"
	BurnTokenAccountName = "burn_token_account"
)

var KeyParams = []byte("paxi_params")

// Params defines spam protection parameters for the Paxi blockchain.
type Params struct {
	ExtraGasPerNewAccount uint64 `json:"extra_gas_per_new_account" yaml:"extra_gas_per_new_account"`
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
		ExtraGasPerNewAccount: uint64(1_200_000),
	}
}

func (p Params) Validate() error {
	if p.ExtraGasPerNewAccount <= 0 {
		return fmt.Errorf("all params from paxi can't be less or equal than 0")
	}
	return nil
}
