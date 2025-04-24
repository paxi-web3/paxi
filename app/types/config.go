package types

import (
	"encoding/json"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func AddCustomGenesis(cdc codec.Codec, pGenesisData *map[string]json.RawMessage) {
	// Add custom genesis state here
	// Bank module
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

	// Auth + Vesting module
	authGenesis := authtypes.DefaultGenesisState()
	addrStr := "paxi1hv7txq4kmc9spkk3nrq0xnhj0jkhk4ye53f50h"
	addr, err := sdk.AccAddressFromBech32(addrStr)
	if err != nil {
		panic(err)
	}

	now := time.Now().Unix()
	periodLength := int64(60 * 60 * 24 * 30 * 4)

	vestingAcc := &vestingtypes.PeriodicVestingAccount{
		BaseVestingAccount: &vestingtypes.BaseVestingAccount{
			BaseAccount: &authtypes.BaseAccount{
				Address:       addr.String(),
				AccountNumber: 0,
				Sequence:      0,
			},
			OriginalVesting: sdk.Coins{sdk.NewInt64Coin(denom, totalAmount)},
		},
		StartTime:      now,
		VestingPeriods: periods,
	}

	authGenesis.Accounts = append(authGenesis.Accounts, vestingAcc)

	(*pGenesisData)[authtypes.ModuleName] = cdc.MustMarshalJSON(authGenesis)
}
