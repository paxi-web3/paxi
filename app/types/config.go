package types

import (
	"encoding/json"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	Denom    = "stake"
	Exponent = 6
)

func AddCustomGenesis(cdc codec.Codec, pGenesisData *map[string]json.RawMessage) {
	// Add custom genesis state here
	// Bank module
	bankGenesis := banktypes.DefaultGenesisState()

	bankGenesis.DenomMetadata = append(bankGenesis.DenomMetadata, banktypes.Metadata{
		Description: "The native token of the Paxi network.",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    Denom, // base denom
				Exponent: 0,
				Aliases:  []string{},
			},
			{
				Denom:    "paxi", // display denom
				Exponent: Exponent,
				Aliases:  []string{"PAXI"},
			},
		},
		Base:    Denom,
		Display: "paxi",
		Name:    "Paxi",
		Symbol:  "PAXI",
	})

	// Auth + Vesting module
	authGenesis := authtypes.DefaultGenesisState()
	vestingAddrList := []string{
		"paxi1hv7txq4kmc9spkk3nrq0xnhj0jkhk4ye53f50h", // devteam
		"paxi1drxrxujqf78zeej3nthgeegt7apax2j48wauqy", // foundation
		"paxi1yp4yc662nhvm9kx7nwexw4t97rlsm9q0kppvr5", // vc
	}
	vestingAmountList := [][]int64{
		{3_000_000, 3_000_000, 3_000_000, 3_000_000, 3_000_000}, // devteam
		{4_000_000, 1_500_000, 1_500_000, 1_500_000, 1_500_000}, // foundation
		{3_000_000, 3_000_000, 3_000_000, 3_000_000, 3_000_000}, // vc
	}

	now := time.Now().Unix()
	periodLength := int64(60 * 60 * 24 * 30 * 4)
	scaling := sdkmath.NewInt(10).ToLegacyDec().Power(Exponent).TruncateInt64()

	for i, addr := range vestingAddrList {
		// set vesting period
		vestingPeriods := []vestingtypes.Period{}
		var totalAmount int64 = 0
		for j, amount := range vestingAmountList[i] {
			stakeAmount := amount * scaling
			period := periodLength
			if j == 0 {
				period = 0
			}
			vestingPeriods = append(vestingPeriods, vestingtypes.Period{
				Length: period,
				Amount: sdk.Coins{sdk.NewInt64Coin(Denom, stakeAmount)},
			})
			totalAmount += stakeAmount
		}

		// add periodic vesting account to auth genesis
		vestingAcc := &vestingtypes.PeriodicVestingAccount{
			BaseVestingAccount: &vestingtypes.BaseVestingAccount{
				BaseAccount: &authtypes.BaseAccount{
					Address:       addr,
					AccountNumber: 0,
					Sequence:      0,
				},
				OriginalVesting: sdk.Coins{sdk.NewInt64Coin(Denom, totalAmount)},
			},
			StartTime:      now,
			VestingPeriods: vestingPeriods,
		}

		anyAcc, err := codectypes.NewAnyWithValue(vestingAcc)
		if err != nil {
			panic(err)
		}
		authGenesis.Accounts = append(authGenesis.Accounts, anyAcc)
	}

	// Grant 5% to DAO
	govGrantedStake := 5_000_000 * scaling
	govAddr := authtypes.NewModuleAddress(govtypes.ModuleName)
	bankGenesis.Balances = append(bankGenesis.Balances, banktypes.Balance{
		Address: govAddr.String(),
		Coins: sdk.Coins{
			sdk.NewInt64Coin(Denom, govGrantedStake),
		},
	})

	govAcc := &authtypes.ModuleAccount{
		BaseAccount: &authtypes.BaseAccount{
			Address: govAddr.String(),
		},
		Name:        govtypes.ModuleName,
		Permissions: []string{authtypes.Burner},
	}

	anyGov, err := codectypes.NewAnyWithValue(govAcc)
	if err != nil {
		panic(err)
	}
	authGenesis.Accounts = append(authGenesis.Accounts, anyGov)

	// Rewrite genesis data
	(*pGenesisData)[banktypes.ModuleName] = cdc.MustMarshalJSON(bankGenesis)
	(*pGenesisData)[authtypes.ModuleName] = cdc.MustMarshalJSON(authGenesis)
}
