package paxi

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
	"github.com/paxi-web3/paxi/x/paxi/keeper"
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
	if err := paxitypes.RegisterQueryHandlerClient(context.Background(), mux, paxitypes.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

func (am AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesis := paxitypes.DefaultGenesisState()
	bz, err := json.Marshal(genesis)
	if err != nil {
		panic(err)
	}
	return bz
}

func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var data paxitypes.GenesisState
	if err := json.Unmarshal(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", paxitypes.ModuleName, err)
	}
	return data.Validate()
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
	var genesisState paxitypes.GenesisState
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
	// Update vesting records every 5000 blocks
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	bh := sdkCtx.BlockHeight()
	if bh%5000 == 0 || bh == 1 {
		am.keeper.SetLockedVestingToStore(sdkCtx)
	}
	return nil
}

func (am AppModule) EndBlock(ctx context.Context) ([]abci.ValidatorUpdate, error) {
	am.blockStatusKeeper.SetLastBlockGasUsed()
	am.blockStatusKeeper.WriteBlockStatusToFile()
	return []abci.ValidatorUpdate{}, nil
}
