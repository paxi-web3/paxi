package prediction

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
	"github.com/paxi-web3/paxi/x/prediction/keeper"
	predictiontypes "github.com/paxi-web3/paxi/x/prediction/types"
)

const ConsensusVersion = 1

var (
	_ module.AppModuleBasic  = AppModule{}
	_ module.HasABCIEndBlock = AppModule{}

	_ appmodule.AppModule       = AppModule{}
	_ appmodule.HasBeginBlocker = AppModule{}
)

type AppModuleBasic struct {
	cdc codec.Codec
}

func (m AppModule) ConsensusVersion() uint64 { return ConsensusVersion }
func (am AppModule) IsOnePerModuleType()     {}
func (am AppModule) IsAppModule()            {}

func (AppModuleBasic) Name() string                                    { return predictiontypes.ModuleName }
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {}

func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	predictiontypes.RegisterMsg(registry)
}

func (am AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *gwruntime.ServeMux) {
	if err := predictiontypes.RegisterQueryHandlerClient(context.Background(), mux, predictiontypes.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

func (am AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesis := predictiontypes.DefaultGenesisState()
	bz, err := json.Marshal(genesis)
	if err != nil {
		panic(err)
	}
	return bz
}

func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var data predictiontypes.GenesisState
	if err := json.Unmarshal(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", predictiontypes.ModuleName, err)
	}
	return data.Validate()
}

type AppModule struct {
	AppModuleBasic
	keeper keeper.Keeper
}

func NewAppModule(cdc codec.Codec, keeper keeper.Keeper) AppModule {
	return AppModule{AppModuleBasic: AppModuleBasic{cdc: cdc}, keeper: keeper}
}

func (am AppModule) Name() string                                { return predictiontypes.ModuleName }
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	predictiontypes.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServer(am.keeper))
	predictiontypes.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) {
	var genesisState predictiontypes.GenesisState
	if err := json.Unmarshal(data, &genesisState); err != nil {
		panic(fmt.Sprintf("failed to unmarshal %s genesis state: %v", predictiontypes.ModuleName, err))
	}
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
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	am.keeper.AutoCloseExpiredMarkets(sdkCtx)
	am.keeper.AutoPruneOrders(sdkCtx)
	return nil
}

func (am AppModule) EndBlock(ctx context.Context) ([]abci.ValidatorUpdate, error) {
	return []abci.ValidatorUpdate{}, nil
}
