package keeper

import (
	"context"
	"encoding/json"
	"fmt"

	store "cosmossdk.io/core/store"
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/utils"
	"github.com/paxi-web3/paxi/x/prediction/types"
)

type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}

type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
}

type PRC20Keeper interface {
	Execute(ctx sdk.Context, contractAddress, caller sdk.AccAddress, msg []byte, coins sdk.Coins) ([]byte, error)
}

type PRC20QueryKeeper interface {
	QuerySmart(ctx context.Context, contractAddress sdk.AccAddress, req []byte) ([]byte, error)
}

type Keeper struct {
	cdc           codec.BinaryCodec
	bankKeeper    BankKeeper
	accountKeeper AccountKeeper
	prc20Keeper   PRC20Keeper
	prc20Query    PRC20QueryKeeper
	storeKey      storetypes.StoreKey
	storeService  store.KVStoreService
	authority     string
}

func NewKeeper(
	cdc codec.BinaryCodec,
	bankKeeper BankKeeper,
	accountKeeper AccountKeeper,
	prc20Keeper PRC20Keeper,
	prc20QueryKeeper PRC20QueryKeeper,
	storeKey storetypes.StoreKey,
	storeService store.KVStoreService,
	authority string,
) Keeper {
	return Keeper{
		cdc:           cdc,
		bankKeeper:    bankKeeper,
		accountKeeper: accountKeeper,
		prc20Keeper:   prc20Keeper,
		prc20Query:    prc20QueryKeeper,
		storeKey:      storeKey,
		storeService:  storeService,
		authority:     authority,
	}
}

func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	k.SetParams(ctx, data.Params)
	k.SetNextMarketID(ctx, data.NextMarketID)
	k.SetNextOrderID(ctx, data.NextOrderID)

	for _, market := range data.Markets {
		k.SetMarket(ctx, market)
	}
	for _, order := range data.Orders {
		k.SetOrder(ctx, order)
	}
	for _, pos := range data.Positions {
		k.SetPosition(ctx, pos)
	}
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	gs := types.DefaultGenesisState()
	gs.Params = k.GetParams(ctx)
	gs.NextMarketID = k.GetNextMarketID(ctx)
	gs.NextOrderID = k.GetNextOrderID(ctx)
	gs.Markets = k.GetAllMarkets(ctx)
	gs.Orders = k.GetAllOrders(ctx)
	gs.Positions = k.GetAllPositions(ctx)
	return &gs
}

func parsePositiveInt(value, field string) (sdkmath.Int, error) {
	amount, ok := sdkmath.NewIntFromString(value)
	if !ok {
		return sdkmath.Int{}, fmt.Errorf("invalid %s", field)
	}
	if !amount.IsPositive() {
		return sdkmath.Int{}, fmt.Errorf("%s must be positive", field)
	}
	return amount, nil
}

func parseNonNegativeInt(value, field string) (sdkmath.Int, error) {
	amount, ok := sdkmath.NewIntFromString(value)
	if !ok {
		return sdkmath.Int{}, fmt.Errorf("invalid %s", field)
	}
	if amount.IsNegative() {
		return sdkmath.Int{}, fmt.Errorf("%s cannot be negative", field)
	}
	return amount, nil
}

func (k Keeper) mustSet(store store.KVStore, key []byte, value []byte) {
	if err := store.Set(key, value); err != nil {
		panic(err)
	}
}

func (k Keeper) mustDelete(store store.KVStore, key []byte) {
	if err := store.Delete(key); err != nil {
		panic(err)
	}
}

func (k Keeper) mustMarshalJSON(v any) []byte {
	bz, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return bz
}

func (k Keeper) getPositionOrDefault(ctx sdk.Context, marketID uint64, addr sdk.AccAddress) *types.Position {
	pos, found := k.GetPosition(ctx, marketID, addr)
	if found {
		return pos
	}
	return &types.Position{
		MarketId:          marketID,
		Address:           addr.String(),
		YesShares:         sdkmath.ZeroInt().String(),
		LockedYesShares:   sdkmath.ZeroInt().String(),
		NoShares:          sdkmath.ZeroInt().String(),
		LockedNoShares:    sdkmath.ZeroInt().String(),
		ClaimedPayout:     false,
		ClaimedVoidRefund: false,
	}
}

func ensureNoNegative(nums ...sdkmath.Int) error {
	for i := range nums {
		if nums[i].IsNegative() {
			return fmt.Errorf("negative balance")
		}
	}
	return nil
}

func prefixIterator(store store.KVStore, prefix []byte) (store.Iterator, error) {
	start, end := utils.PrefixRange(prefix)
	return store.Iterator(start, end)
}
