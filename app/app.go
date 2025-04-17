// Full Paxi Cosmos SDK App with Module Manager (v0.47.17-compatible)
package app

import (
	"io"

	version "github.com/paxi-web3/paxi"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	wasm "github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

// ModuleBasics defines all basic Cosmos SDK modules
var ModuleBasics = module.NewBasicManager(
	auth.AppModuleBasic{},
	bank.AppModuleBasic{},
	staking.AppModuleBasic{},
	distribution.AppModuleBasic{},
	mint.AppModuleBasic{},
	slashing.AppModuleBasic{},
	gov.AppModuleBasic{},
	crisis.AppModuleBasic{},
	params.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	wasm.AppModuleBasic{},
)

// Keepers manage module state
var (
	AccountKeeper      authkeeper.AccountKeeper
	BankKeeper         bankkeeper.Keeper
	StakingKeeper      stakingkeeper.Keeper
	MintKeeper         mintkeeper.Keeper
	DistributionKeeper distributionkeeper.Keeper
	SlashingKeeper     slashingkeeper.Keeper
	CrisisKeeper       crisiskeeper.Keeper
	ParamsKeeper       paramskeeper.Keeper
	UpgradeKeeper      upgradekeeper.Keeper
	GovKeeper          govkeeper.Keeper
	WasmKeeper         wasmkeeper.Keeper
)

// EncodingConfig defines the codec and tx config
type EncodingConfig struct {
	InterfaceRegistry types.InterfaceRegistry
	Marshaler         codec.Codec
	TxConfig          auth.TxConfig
	Amino             *codec.LegacyAmino
}

// MakeEncodingConfig creates encoding and tx configuration
func MakeEncodingConfig() EncodingConfig {
	interfaceRegistry := types.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	return EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          authtypes.NewTxConfig(marshaler, authtypes.DefaultSignModes),
		Amino:             codec.NewLegacyAmino(),
	}
}

// GovRouter defines gov proposal handler routing
func GovRouter() govtypes.Router {
	r := govtypes.NewRouter()
	return *r
}

// NewPaxiApp constructs and returns a Cosmos SDK application instance
func NewPaxiApp(logger io.Writer, db store.CommitMultiStore, traceStore io.Writer, loadLatest bool, skipUpgradeHeights map[int64]bool, appOpts interface{}) *baseapp.BaseApp {
	encodingConfig := MakeEncodingConfig()

	// Define store keys for modules
	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		distributiontypes.StoreKey, slashingtypes.StoreKey, minttypes.StoreKey,
		govtypes.StoreKey, crisistypes.StoreKey, paramstypes.StoreKey, wasmtypes.StoreKey,
		upgradetypes.StoreKey,
	)

	tkeys := storetypes.NewTransientStoreKeys(paramstypes.TStoreKey)

	app := baseapp.NewBaseApp("paxi", logger, db, encodingConfig.TxConfig.TxDecoder(), nil)
	app.SetCommitMultiStoreTracer(traceStore)
	app.SetVersion(version.Version)

	ParamsKeeper = paramskeeper.NewKeeper(encodingConfig.Marshaler, encodingConfig.Amino, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])
	AccountKeeper = authkeeper.NewAccountKeeper(encodingConfig.Marshaler, keys[authtypes.StoreKey], ParamsKeeper.Subspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, map[string][]string{})
	BankKeeper = bankkeeper.NewBaseKeeper(encodingConfig.Marshaler, keys[banktypes.StoreKey], AccountKeeper, ParamsKeeper.Subspace(banktypes.ModuleName), map[string]bool{})
	StakingKeeper = stakingkeeper.NewKeeper(encodingConfig.Marshaler, keys[stakingtypes.StoreKey], AccountKeeper, BankKeeper, ParamsKeeper.Subspace(stakingtypes.ModuleName))
	MintKeeper = mintkeeper.NewKeeper(encodingConfig.Marshaler, keys[minttypes.StoreKey], ParamsKeeper.Subspace(minttypes.ModuleName), StakingKeeper, AccountKeeper, BankKeeper, authtypes.FeeCollectorName)
	DistributionKeeper = distributionkeeper.NewKeeper(encodingConfig.Marshaler, keys[distributiontypes.StoreKey], ParamsKeeper.Subspace(distributiontypes.ModuleName), AccountKeeper, BankKeeper, StakingKeeper, authtypes.FeeCollectorName, govtypes.ModuleName)
	SlashingKeeper = slashingkeeper.NewKeeper(encodingConfig.Marshaler, keys[slashingtypes.StoreKey], StakingKeeper, ParamsKeeper.Subspace(slashingtypes.ModuleName))
	GovKeeper = govkeeper.NewKeeper(encodingConfig.Marshaler, keys[govtypes.StoreKey], ParamsKeeper.Subspace(govtypes.ModuleName), AccountKeeper, BankKeeper, GovRouter())
	CrisisKeeper = crisiskeeper.NewKeeper(ParamsKeeper.Subspace(crisistypes.ModuleName), nil, BankKeeper, authtypes.FeeCollectorName)
	UpgradeKeeper = upgradekeeper.NewKeeper(skipUpgradeHeights, keys[upgradetypes.StoreKey], app, logger)
	WasmKeeper = wasmkeeper.NewKeeper(encodingConfig.Marshaler, keys[wasmtypes.StoreKey], AccountKeeper, BankKeeper, StakingKeeper, DistributionKeeper, GovKeeper, wasm.DefaultWasmConfig(), nil, nil)

	mm := module.NewManager(
		auth.NewAppModule(encodingConfig.Marshaler, AccountKeeper, nil),
		bank.NewAppModule(encodingConfig.Marshaler, BankKeeper, AccountKeeper),
		staking.NewAppModule(encodingConfig.Marshaler, StakingKeeper, AccountKeeper, BankKeeper),
		mint.NewAppModule(encodingConfig.Marshaler, MintKeeper, AccountKeeper),
		distribution.NewAppModule(encodingConfig.Marshaler, DistributionKeeper, AccountKeeper, BankKeeper, StakingKeeper),
		slashing.NewAppModule(encodingConfig.Marshaler, SlashingKeeper, AccountKeeper, BankKeeper, StakingKeeper),
		gov.NewAppModule(encodingConfig.Marshaler, GovKeeper, AccountKeeper, BankKeeper),
		crisis.NewAppModule(CrisisKeeper),
		wasm.NewAppModule(encodingConfig.Marshaler, WasmKeeper, AccountKeeper, BankKeeper),
		upgrade.NewAppModule(UpgradeKeeper),
	)

	mm.SetOrderBeginBlockers(upgradetypes.ModuleName, minttypes.ModuleName, distributiontypes.ModuleName, slashingtypes.ModuleName)
	mm.SetOrderEndBlockers(stakingtypes.ModuleName, govtypes.ModuleName)
	mm.SetOrderInitGenesis(
		authtypes.ModuleName, banktypes.ModuleName, stakingtypes.ModuleName,
		minttypes.ModuleName, distributiontypes.ModuleName, slashingtypes.ModuleName,
		govtypes.ModuleName, crisistypes.ModuleName, wasmtypes.ModuleName,
	)

	mm.RegisterInvariants(nil)
	mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)
	app.SetInitChainer(mm.InitGenesis)
	app.SetBeginBlocker(mm.BeginBlock)
	app.SetEndBlocker(mm.EndBlock)

	return app
}
