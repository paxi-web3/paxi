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
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	gogoany "github.com/cosmos/gogoproto/types/any"
	paxitypes "github.com/paxi-web3/paxi/x/paxi/types"
)

const (
	DefaultDenom = "upaxi"
	Exponent     = 6
)

func AddCustomGenesis(cdc codec.Codec, pGenesisData *map[string]json.RawMessage) {
	// Add custom genesis state here
	// Bank module
	bankGenesis := banktypes.DefaultGenesisState()

	bankGenesis.DenomMetadata = append(bankGenesis.DenomMetadata, banktypes.Metadata{
		Description: "The native token of the Paxi network.",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    DefaultDenom, // base denom
				Exponent: 0,
				Aliases:  []string{},
			},
			{
				Denom:    "paxi", // display denom
				Exponent: Exponent,
				Aliases:  []string{"PAXI"},
			},
		},
		Base:    DefaultDenom,
		Display: "paxi",
		Name:    "Paxi",
		Symbol:  "PAXI",
	})

	// Auth + Vesting module
	totalSupply := int64(0)
	authGenesis := authtypes.DefaultGenesisState()
	vestingAddrList := []string{
		"paxi1hv7txq4kmc9spkk3nrq0xnhj0jkhk4ye53f50h", // devteam
		"paxi1drxrxujqf78zeej3nthgeegt7apax2j48wauqy", // foundation
		"paxi1yp4yc662nhvm9kx7nwexw4t97rlsm9q0kppvr5", // vc
	}
	vestingAmountList := [][]int64{
		{2_000_000, 3_000_000, 3_000_000, 3_000_000, 3_000_000}, // devteam
		{4_000_000, 1_500_000, 1_500_000, 1_500_000, 1_500_000}, // foundation
		{3_000_000, 3_000_000, 3_000_000, 3_000_000, 3_000_000}, // vc
	}

	now := time.Now().Unix()
	periodLength := int64(60 * 60 * 24 * 30 * 4)
	scaling := sdkmath.NewInt(10).ToLegacyDec().Power(Exponent).TruncateInt64()

	for i, addr := range vestingAddrList {
		// set vesting period
		vestingPeriods := []vestingtypes.Period{}
		var subTotalStake int64 = 0
		var endtime int64 = 0

		for j, amount := range vestingAmountList[i] {
			stakeAmount := amount * scaling
			period := periodLength
			if j == 0 {
				period = 0
			}

			vestingPeriods = append(vestingPeriods, vestingtypes.Period{
				Length: period,
				Amount: sdk.Coins{sdk.NewInt64Coin(DefaultDenom, stakeAmount)},
			})
			endtime += period
			subTotalStake += stakeAmount
		}

		// add periodic vesting account to auth genesis
		vestingAcc := &vestingtypes.PeriodicVestingAccount{
			BaseVestingAccount: &vestingtypes.BaseVestingAccount{
				BaseAccount: &authtypes.BaseAccount{
					Address:       addr,
					AccountNumber: 0,
					Sequence:      0,
				},
				OriginalVesting: sdk.Coins{sdk.NewInt64Coin(DefaultDenom, subTotalStake)},
				EndTime:         endtime + now,
			},
			StartTime:      now,
			VestingPeriods: vestingPeriods,
		}

		anyAcc, err := codectypes.NewAnyWithValue(vestingAcc)
		if err != nil {
			panic(err)
		}
		authGenesis.Accounts = append(authGenesis.Accounts, anyAcc)

		// add the account to bank balances
		bankGenesis.Balances = append(bankGenesis.Balances, banktypes.Balance{
			Address: addr,
			Coins:   sdk.Coins{sdk.NewInt64Coin(DefaultDenom, subTotalStake)},
		})

		// add total supply
		totalSupply += subTotalStake
	}

	// Grant 5% to DAO (community pool) / 45% to ICO / 10 % to Incentives
	distrAddr := authtypes.NewModuleAddress(distrtypes.ModuleName)
	icoAddrString := "paxi1rnwv6u92cc0v55ed8vk7zy3cjd32j5qq4487wg"
	incentivesAddrString := "paxi1ag0jvu2uswq2355m573a9exx2ew29gfnuxh5px"

	grantAddrList := []string{
		distrAddr.String(),   // DAO
		icoAddrString,        // ICO
		incentivesAddrString, // Incentives
	}
	grantStakeList := []int64{
		5_000_000 * scaling,  // DAO
		45_000_000 * scaling, // ICO
		10_000_000 * scaling, // Incentives
	}
	grantAccountList := []*gogoany.Any{
		// DAO
		(func() *gogoany.Any {
			distrAcc := &authtypes.ModuleAccount{
				BaseAccount: authtypes.NewBaseAccount(distrAddr, nil, 0, 0),
				Name:        distrtypes.ModuleName,
				Permissions: []string{},
			}
			anyDistr, err := codectypes.NewAnyWithValue(distrAcc)
			if err != nil {
				panic(err)
			}
			return anyDistr
		})(),
		// ICO
		(func() *gogoany.Any {
			addr, err := sdk.AccAddressFromBech32(icoAddrString)
			if err != nil {
				panic(err)
			}
			acc := authtypes.NewBaseAccountWithAddress(addr)
			anyAcc, err := codectypes.NewAnyWithValue(acc)
			if err != nil {
				panic(err)
			}
			return anyAcc
		})(),
		// Incentives
		(func() *gogoany.Any {
			addr, err := sdk.AccAddressFromBech32(incentivesAddrString)
			if err != nil {
				panic(err)
			}
			acc := authtypes.NewBaseAccountWithAddress(addr)
			anyAcc, err := codectypes.NewAnyWithValue(acc)
			if err != nil {
				panic(err)
			}
			return anyAcc
		})(),
	}

	for i, addr := range grantAddrList {
		// Add balance
		bankGenesis.Balances = append(bankGenesis.Balances, banktypes.Balance{
			Address: addr,
			Coins: sdk.Coins{
				sdk.NewInt64Coin(DefaultDenom, grantStakeList[i]),
			},
		})

		// Add account
		authGenesis.Accounts = append(authGenesis.Accounts, grantAccountList[i])

		// Add to total supply
		totalSupply += grantStakeList[i]
	}

	// Grant 5% from distribution account to community pool
	distrGenesis := distrtypes.DefaultGenesisState()
	distrGenesis.FeePool.CommunityPool = []sdk.DecCoin{sdk.NewDecCoin(DefaultDenom, sdkmath.NewInt(grantStakeList[0]))}

	// Modify staking genesis
	stakingGenesis := stakingtypes.DefaultGenesisState()
	stakingGenesis.Params = stakingtypes.Params{
		UnbondingTime:     time.Hour * 24 * 7 * 1, // 1 week
		MaxValidators:     100,
		MaxEntries:        7,
		HistoricalEntries: 10000,
		BondDenom:         DefaultDenom,
		MinCommissionRate: stakingGenesis.Params.MinCommissionRate, // 0.05
	}

	// Add total supply
	bankGenesis.Supply = bankGenesis.Supply.Add(sdk.NewInt64Coin(DefaultDenom, totalSupply))

	// Set denom from stake to upaxi
	govGenesis := govv1.DefaultGenesisState()
	govGenesis.Params.ExpeditedMinDeposit[0].Denom = DefaultDenom
	govGenesis.Params.MinDeposit[0].Denom = DefaultDenom

	// Add burn token module account
	burnTokenAcc := &authtypes.ModuleAccount{
		BaseAccount: authtypes.NewBaseAccount(authtypes.NewModuleAddress(paxitypes.BurnTokenAccountName), nil, 0, 0),
		Name:        paxitypes.BurnTokenAccountName,
		Permissions: []string{authtypes.Burner},
	}
	anyBurnToken, err := codectypes.NewAnyWithValue(burnTokenAcc)
	if err != nil {
		panic(err)
	}
	authGenesis.Accounts = append(authGenesis.Accounts, anyBurnToken)

	// custom slasing
	slasingGenesis := slashingtypes.DefaultGenesisState()
	slasingParams := slashingtypes.DefaultParams()
	slasingParams.DowntimeJailDuration = 10 * time.Minute
	slasingParams.MinSignedPerWindow = sdkmath.LegacyMustNewDecFromStr("0.1")
	slasingParams.SignedBlocksWindow = int64(100)
	slasingParams.SlashFractionDoubleSign = sdkmath.LegacyMustNewDecFromStr("0.1")
	slasingParams.SlashFractionDowntime = sdkmath.LegacyMustNewDecFromStr("0.01")
	slasingGenesis.Params = slasingParams

	// Rewrite genesis data
	(*pGenesisData)[govtypes.ModuleName] = cdc.MustMarshalJSON(govGenesis)
	(*pGenesisData)[banktypes.ModuleName] = cdc.MustMarshalJSON(bankGenesis)
	(*pGenesisData)[authtypes.ModuleName] = cdc.MustMarshalJSON(authGenesis)
	(*pGenesisData)[stakingtypes.ModuleName] = cdc.MustMarshalJSON(stakingGenesis)
	(*pGenesisData)[distrtypes.ModuleName] = cdc.MustMarshalJSON(distrGenesis)
	(*pGenesisData)[slashingtypes.ModuleName] = cdc.MustMarshalJSON(slasingGenesis)
}
