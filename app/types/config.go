package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func AddCustomGenesis(cdc codec.Codec, pGenesisData *map[string]json.RawMessage) {
	// Add custom genesis state here
	bankGenesis := banktypes.DefaultGenesisState()

	bankGenesis.DenomMetadata = append(bankGenesis.DenomMetadata, banktypes.Metadata{
		Description: "The native token of the Paxi network.",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "stake", // base denom
				Exponent: 0,
				Aliases:  []string{},
			},
			{
				Denom:    "paxi", // display denom
				Exponent: 6,
				Aliases:  []string{"PAXI"},
			},
		},
		Base:    "stake",
		Display: "paxi",
		Name:    "Paxi",
		Symbol:  "PAXI",
	})

	(*pGenesisData)[banktypes.ModuleName] = cdc.MustMarshalJSON(bankGenesis)
}
