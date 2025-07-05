package swap

import (
	"context"
	"encoding/json"
	"fmt"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"

	"cosmossdk.io/core/appmodule"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/paxi-web3/paxi/x/swap/keeper"
	swaptypes "github.com/paxi-web3/paxi/x/swap/types"
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

// AppModuleBasic implements the AppModuleBasic interface for the swap module
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

func (AppModuleBasic) Name() string                                    { return swaptypes.ModuleName }
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
	swaptypes.RegisterMsg(registry) // Register the module-specific message types
}

func (am AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *gwruntime.ServeMux) {
	// Register custom gRPC gateway routes here
	if err := swaptypes.RegisterQueryHandlerClient(context.Background(), mux, swaptypes.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

func (am AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesis := swaptypes.DefaultGenesisState()
	bz, err := json.Marshal(genesis)
	if err != nil {
		panic(err)
	}
	return bz
}

func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var data swaptypes.GenesisState
	if err := json.Unmarshal(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", swaptypes.ModuleName, err)
	}
	return data.Validate()
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

func (am AppModule) Name() string { return swaptypes.ModuleName }

func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	swaptypes.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServer(am.keeper))
	swaptypes.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) {
	var genesisState swaptypes.GenesisState
	json.Unmarshal(data, &genesisState)
	am.keeper.InitGenesis(ctx, genesisState)
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := am.keeper.ExportGenesis(ctx)
	bz, err := json.Marshal(gs)
	if err != nil {
		panic(err)
	}
	return bz
}

func (am AppModule) BeginBlock(ctx context.Context) error {
	return nil
}

func (am AppModule) EndBlock(ctx context.Context) ([]abci.ValidatorUpdate, error) {
	return []abci.ValidatorUpdate{}, nil
}
