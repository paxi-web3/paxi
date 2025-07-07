package keeper

import (
	store "cosmossdk.io/core/store"
	storetypes "cosmossdk.io/store/types"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	swaptypes "github.com/paxi-web3/paxi/x/swap/types"
)

type Keeper struct {
	cdc             codec.BinaryCodec
	bankKeeper      bankkeeper.Keeper
	accountKeeper   authkeeper.AccountKeeper
	wasmKeeper      wasmtypes.ContractOpsKeeper
	wasmQueryKeeper wasmkeeper.Keeper
	storeKey        storetypes.StoreKey
	storeService    store.KVStoreService
	authority       string
}

type BlockStatsKeeper interface {
}

func NewKeeper(cdc codec.BinaryCodec, bankKeeper bankkeeper.Keeper, accountKeeper authkeeper.AccountKeeper, wasmKeeper wasmtypes.ContractOpsKeeper, wasmQueryKeeper wasmkeeper.Keeper, storeKey storetypes.StoreKey, storeService store.KVStoreService, authority string) Keeper {
	return Keeper{
		cdc:             cdc,
		bankKeeper:      bankKeeper,
		accountKeeper:   accountKeeper,
		wasmKeeper:      wasmKeeper,
		wasmQueryKeeper: wasmQueryKeeper,
		storeKey:        storeKey,
		storeService:    storeService,
		authority:       authority,
	}
}

func (k Keeper) InitGenesis(ctx sdk.Context, data swaptypes.GenesisState) {
	// Set params
	k.SetParams(ctx, data.Params)
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *swaptypes.GenesisState {
	params := k.GetParams(ctx)
	return &swaptypes.GenesisState{
		Params: params,
	}
}
