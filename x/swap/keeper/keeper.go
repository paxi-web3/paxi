package keeper

import (
	"encoding/json"

	storetypes "cosmossdk.io/core/store"
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
	storeService    storetypes.KVStoreService
	authority       string
}

type BlockStatsKeeper interface {
}

func NewKeeper(cdc codec.BinaryCodec, bankKeeper bankkeeper.Keeper, accountKeeper authkeeper.AccountKeeper, wasmKeeper wasmtypes.ContractOpsKeeper, wasmQueryKeeper wasmkeeper.Keeper, storeService storetypes.KVStoreService, authority string) Keeper {
	return Keeper{
		cdc:             cdc,
		bankKeeper:      bankKeeper,
		accountKeeper:   accountKeeper,
		wasmKeeper:      wasmKeeper,
		wasmQueryKeeper: wasmQueryKeeper,
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

var poolPrefix = []byte{0x01}

func poolStoreKey(prc20 string) []byte {
	return append(poolPrefix, []byte(prc20)...) // key = 0x01 + prc20 address
}

// SetPool saves the pool to KVStore
func (k Keeper) SetPool(ctx sdk.Context, pool swaptypes.Pool) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := json.Marshal(&pool)
	if err != nil {
		panic(err)
	}
	err = store.Set(poolStoreKey(pool.Prc20), bz)
	if err != nil {
		panic(err)
	}
}

// GetPool fetches a pool by PRC20 contract address
func (k Keeper) GetPool(ctx sdk.Context, prc20 string) (swaptypes.Pool, bool) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(poolStoreKey(prc20))
	if err != nil || bz == nil {
		return swaptypes.Pool{}, false
	}
	var pool swaptypes.Pool
	err = json.Unmarshal(bz, &pool)
	if err != nil {
		return swaptypes.Pool{}, false
	}
	return pool, true
}
