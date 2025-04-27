package paxi

import (
	"context"
	"encoding/json"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"

	"cosmossdk.io/core/appmodule"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/paxi-web3/paxi/x/paxi/keeper"
	rest "github.com/paxi-web3/paxi/x/paxi/rest"
	paxitypes "github.com/paxi-web3/paxi/x/paxi/types"
)

const (
	ConsensusVersion = 1
)

var (
	_ module.AppModuleBasic  = AppModule{}
	_ module.HasABCIEndBlock = AppModule{}

	_ appmodule.AppModule       = AppModule{}
	_ appmodule.HasBeginBlocker = AppModule{}
)

// AppModuleBasic implements the AppModuleBasic interface for the paxi module
// It handles codec registration and default genesis

type AppModuleBasic struct {
	cdc codec.Codec
}

// Return the version of the module
func (m AppModule) ConsensusVersion() uint64 {
	return ConsensusVersion
}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

func (AppModuleBasic) Name() string                                    { return paxitypes.ModuleName }
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {}

// RegisterInterfaces registers the module's message types and interfaces into the global interface registry.
//
// This allows the Cosmos SDK to correctly encode/decode the module's Msgs (e.g., MsgBurnToken)
// and recognize them during transaction processing.
//
// The registry is used to map concrete Go types (like MsgBurnToken) to protobuf type URLs,
// so that the Cosmos SDK's signing and transaction systems can work correctly.
//
// Called automatically during app initialization.
func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	paxitypes.RegisterMsg(registry) // Register the module-specific message types
}

func (am AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *gwruntime.ServeMux) {
	// Register custom gRPC gateway routes here
}

func (am AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr *mux.Router) {
	// Register custom REST gateway routes here
	rtr.HandleFunc("/paxi/paxi/locked_vesting", rest.LockedVestingHandler(clientCtx)).Methods("GET")
	rtr.HandleFunc("/paxi/paxi/circulating_supply", rest.CirculatingSupplyHandler(clientCtx)).Methods("GET")
	rtr.HandleFunc("/paxi/paxi/total_supply", rest.TotalSupplyHandler(clientCtx)).Methods("GET")
}

func (am AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return nil
}

func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	return nil
}

// AppModule implements the AppModule interface
// It provides the BeginBlock logic for minting

type AppModule struct {
	AppModuleBasic
	keeper            keeper.Keeper
	blockStatusKeeper keeper.BlockStatsKeeper
}

func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper,
	blockStatusKeeper keeper.BlockStatsKeeper,
) AppModule {
	return AppModule{
		AppModuleBasic:    AppModuleBasic{cdc: cdc},
		keeper:            keeper,
		blockStatusKeeper: blockStatusKeeper,
	}
}

func (am AppModule) Name() string { return paxitypes.ModuleName }

func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	paxitypes.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServer(am.keeper))
	paxitypes.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) {
	am.keeper.InitGenesis(ctx)
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	return nil
}

func (am AppModule) BeginBlock(ctx context.Context) error {
	// Update vesting records every 100 blocks
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	bh := sdkCtx.BlockHeight()
	if bh%100 == 0 {
		am.keeper.SetLockedVestingToStore(sdkCtx)
	}
	return nil
}

func (am AppModule) EndBlock(ctx context.Context) ([]abci.ValidatorUpdate, error) {
	am.blockStatusKeeper.SetLastBlockGasUsed()
	return []abci.ValidatorUpdate{}, nil
}
