package keeper

import (
	"fmt"

	storetypes "cosmossdk.io/core/store"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paxitypes "github.com/paxi-web3/paxi/x/paxi/types"
)

type Keeper struct {
	cdc           codec.BinaryCodec
	bankKeeper    bankkeeper.Keeper
	accountKeeper authkeeper.AccountKeeper
	storeService  storetypes.KVStoreService
}

func NewKeeper(cdc codec.BinaryCodec, bankKeeper bankkeeper.Keeper, accountKeeper authkeeper.AccountKeeper, storeService storetypes.KVStoreService) Keeper {
	return Keeper{
		cdc:           cdc,
		bankKeeper:    bankKeeper,
		accountKeeper: accountKeeper,
		storeService:  storeService,
	}
}

func (k Keeper) GetLockedVestingFromStore(ctx sdk.Context) (lockedVesting sdkmath.LegacyDec, exists bool) {
	store := k.storeService.OpenKVStore(ctx)

	lockedVesting = sdkmath.LegacyNewDecFromInt(sdkmath.ZeroInt())
	exists = false

	// Check if the block height is the same as the one stored
	bhBz, err := store.Get([]byte(paxitypes.BlockHeightKey))
	if err != nil || bhBz == nil {
		return lockedVesting, exists
	}

	bhInt := sdkmath.Int{}
	if err := bhInt.Unmarshal(bhBz); err != nil {
		panic(fmt.Errorf("failed to unmarshal block height: %w", err))
	}

	bh := ctx.BlockHeight()
	if !bhInt.Equal(sdkmath.NewInt(bh)) {
		return lockedVesting, exists
	}

	// Check if the locked vesting is stored
	lvBz, err := store.Get([]byte(paxitypes.LockedVestingKey))
	if err != nil || lvBz == nil {
		return lockedVesting, exists
	}

	if err := lockedVesting.Unmarshal(lvBz); err != nil {
		panic(fmt.Errorf("failed to unmarshal locked vesting: %w", err))
	}

	return lockedVesting, true
}

func (k Keeper) SetLockedVestingToStore(ctx sdk.Context, lvDec sdkmath.LegacyDec) {
	store := k.storeService.OpenKVStore(ctx)

	bh := ctx.BlockHeight()
	bhBz, err := sdkmath.NewInt(bh).Marshal()
	if err != nil {
		panic(fmt.Errorf("failed to marshal block height: %w", err))
	}

	lvBz, err := lvDec.Marshal()
	if err != nil {
		panic(fmt.Errorf("failed to marshal locked vesting: %w", err))
	}

	store.Set([]byte(paxitypes.BlockHeightKey), bhBz)
	store.Set([]byte(paxitypes.LockedVestingKey), lvBz)

}

func prefixRange(prefix []byte) (start []byte, end []byte) {
	if len(prefix) == 0 {
		return nil, nil
	}

	start = prefix

	end = make([]byte, len(prefix))
	copy(end, prefix)

	for i := len(end) - 1; i >= 0; i-- {
		if end[i] < 0xFF {
			end[i]++
			end = end[:i+1]
			return start, end
		}
	}

	return start, nil
}

func (k Keeper) GetLockedVesting(ctx sdk.Context) sdkmath.LegacyDec {
	lockedVesting, exists := k.GetLockedVestingFromStore(ctx)
	if exists {
		return lockedVesting
	}

	store := k.storeService.OpenKVStore(ctx)
	prefix := []byte(paxitypes.VestingAccountPrefix)
	start, end := prefixRange(prefix)
	iterator, err := store.Iterator(start, end)
	if err != nil {
		panic(err)
	}

	// Iterate over all the vesting accounts and sum the locked vesting
	lockedVesting = sdkmath.LegacyNewDecFromInt(sdkmath.ZeroInt())
	for ; iterator.Valid(); iterator.Next() {
		addrStr := string(iterator.Key()[len(paxitypes.VestingAccountPrefix):])
		addr, err := sdk.AccAddressFromBech32(addrStr)
		if err != nil {
			continue
		}

		acc := k.accountKeeper.GetAccount(ctx, addr)
		switch va := acc.(type) {
		case *vestingtypes.ContinuousVestingAccount:
			coins := va.LockedCoins(ctx.BlockTime())
			lockedVesting.Add(sdkmath.LegacyNewDecFromInt(coins.AmountOf(paxitypes.DefaultDenom)))
		case *vestingtypes.DelayedVestingAccount:
			coins := va.LockedCoins(ctx.BlockTime())
			lockedVesting.Add(sdkmath.LegacyNewDecFromInt(coins.AmountOf(paxitypes.DefaultDenom)))
		case *vestingtypes.PeriodicVestingAccount:
			coins := va.LockedCoins(ctx.BlockTime())
			lockedVesting.Add(sdkmath.LegacyNewDecFromInt(coins.AmountOf(paxitypes.DefaultDenom)))
		}
	}

	if lockedVesting.IsPositive() {
		k.SetLockedVestingToStore(ctx, lockedVesting)
	} else {
		lockedVesting = sdkmath.LegacyNewDecFromInt(sdkmath.ZeroInt())
	}

	return lockedVesting
}

func (k Keeper) InitGenesis(ctx sdk.Context) {
	// Interate all the accounts and store the vesting accounts
	accounts := k.accountKeeper.GetAllAccounts(ctx)
	store := k.storeService.OpenKVStore(ctx)
	for _, acc := range accounts {
		switch va := acc.(type) {
		case *vestingtypes.ContinuousVestingAccount,
			*vestingtypes.DelayedVestingAccount,
			*vestingtypes.PeriodicVestingAccount:
			// Store the vesting account
			addrKey := []byte(paxitypes.VestingAccountPrefix + va.String())
			store.Set(addrKey, []byte{1})
		}
	}
}
