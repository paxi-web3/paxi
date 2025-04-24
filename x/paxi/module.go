package paxi

import (
	"context"
	"encoding/json"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"

	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/paxi-web3/paxi/x/paxi/keeper"
	paxitypes "github.com/paxi-web3/paxi/x/paxi/types"
)

const (
	ConsensusVersion = 1
)

var (
	_ module.AppModuleBasic = AppModule{}
	//_ module.HasGenesis     = AppModule{}

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

func (AppModuleBasic) Name() string                                             { return paxitypes.ModuleName }
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino)          {}
func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {}

func (am AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *gwruntime.ServeMux) {
	// Register custom gRPC gateway routes here
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
	keeper keeper.Keeper
}

func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper,
) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc},
		keeper:         keeper,
	}
}

func (am AppModule) Name() string { return paxitypes.ModuleName }

func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	//paxitypes.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServer(am.keeper))
}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) {
	am.keeper.InitGenesis(ctx)
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	return nil
}

func (am AppModule) BeginBlock(ctx context.Context) error {
	return nil
}
