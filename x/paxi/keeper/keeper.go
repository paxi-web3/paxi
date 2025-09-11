package keeper

import (
	"fmt"
	"sort"
	"time"

	store "cosmossdk.io/core/store"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/paxi-web3/paxi/utils"
	paxitypes "github.com/paxi-web3/paxi/x/paxi/types"
)

type Keeper struct {
	cdc               codec.BinaryCodec
	bankKeeper        bankkeeper.Keeper
	accountKeeper     authkeeper.AccountKeeper
	blockStatusKeeper BlockStatsKeeper
	storeService      store.KVStoreService
	authority         string
}

type BlockStatsKeeper interface {
	GetLastBlockGasUsed() uint64
	SetLastBlockGasUsed()
	GetTotalTxs() uint64
	WriteBlockStatusToFile() error
}

func NewKeeper(cdc codec.BinaryCodec, bankKeeper bankkeeper.Keeper, accountKeeper authkeeper.AccountKeeper, storeService store.KVStoreService, blockStatusKeeper BlockStatsKeeper, authority string) Keeper {
	return Keeper{
		cdc:               cdc,
		bankKeeper:        bankKeeper,
		accountKeeper:     accountKeeper,
		storeService:      storeService,
		blockStatusKeeper: blockStatusKeeper,
		authority:         authority,
	}
}

func (k Keeper) GetLockedVestingFromStore(ctx sdk.Context) (lockedVesting sdkmath.LegacyDec, exists bool) {
	store := k.storeService.OpenKVStore(ctx)

	lockedVesting = sdkmath.LegacyNewDecFromInt(sdkmath.ZeroInt())
	exists = false

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

func (k Keeper) SetLockedVestingToStore(ctx sdk.Context) {
	store := k.storeService.OpenKVStore(ctx)
	prefix := []byte(paxitypes.VestingAccountPrefix)
	start, end := utils.PrefixRange(prefix)
	iterator, err := store.Iterator(start, end)
	if err != nil {
		panic(err)
	}

	defer iterator.Close()

	// Iterate over all the vesting accounts and sum the locked vesting
	lockedVesting := sdkmath.LegacyNewDecFromInt(sdkmath.ZeroInt())
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
			lockedVesting = lockedVesting.Add(sdkmath.LegacyNewDecFromInt(coins.AmountOf(paxitypes.DefaultDenom)))
		case *vestingtypes.DelayedVestingAccount:
			coins := va.LockedCoins(ctx.BlockTime())
			lockedVesting = lockedVesting.Add(sdkmath.LegacyNewDecFromInt(coins.AmountOf(paxitypes.DefaultDenom)))
		case *vestingtypes.PeriodicVestingAccount:
			coins := va.LockedCoins(ctx.BlockTime())
			lockedVesting = lockedVesting.Add(sdkmath.LegacyNewDecFromInt(coins.AmountOf(paxitypes.DefaultDenom)))
		}
	}

	if lockedVesting.IsNegative() {
		lockedVesting = sdkmath.LegacyNewDec(0)
	}

	lvBz, err := lockedVesting.Marshal()
	if err != nil {
		panic(fmt.Errorf("failed to marshal locked vesting: %w", err))
	}

	if erro := store.Set([]byte(paxitypes.LockedVestingKey), lvBz); erro != nil {
		panic(fmt.Errorf("failed to store locked vesting: %w", erro))
	}
}

func (k Keeper) GetLockedVesting(ctx sdk.Context) sdkmath.LegacyDec {
	lockedVesting, exists := k.GetLockedVestingFromStore(ctx)
	if exists {
		return lockedVesting
	} else {
		return sdkmath.LegacyNewDec(0)
	}
}

func (k Keeper) InitGenesis(ctx sdk.Context, data paxitypes.GenesisState) {
	// Set params
	k.SetParams(ctx, data.Params)

	// Interate all the accounts and store the vesting accounts
	accounts := k.accountKeeper.GetAllAccounts(ctx)
	store := k.storeService.OpenKVStore(ctx)
	for _, acc := range accounts {
		switch va := acc.(type) {
		case *vestingtypes.ContinuousVestingAccount,
			*vestingtypes.DelayedVestingAccount,
			*vestingtypes.PeriodicVestingAccount:
			// Store the vesting account
			addrKey := []byte(paxitypes.VestingAccountPrefix + va.GetAddress().String())
			if err := store.Set(addrKey, []byte{1}); err != nil {
				panic(fmt.Errorf("failed to store vesting account: %w", err))
			}
		}
	}
}

func (k Keeper) BurnFromUser(ctx sdk.Context, sender sdk.AccAddress, amount sdk.Coins) error {
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, paxitypes.BurnTokenAccountName, amount)
	if err != nil {
		return err
	}

	err = k.bankKeeper.BurnCoins(ctx, paxitypes.BurnTokenAccountName, amount)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetUnlockSchedules(ctx sdk.Context) []*paxitypes.UnlockSchedule {
	store := k.storeService.OpenKVStore(ctx)
	prefix := []byte(paxitypes.VestingAccountPrefix)
	start, end := utils.PrefixRange(prefix)
	iterator, err := store.Iterator(start, end)
	if err != nil {
		panic(err)
	}

	defer iterator.Close()

	var unlockSchedules []*paxitypes.UnlockSchedule
	for ; iterator.Valid(); iterator.Next() {
		addrStr := string(iterator.Key()[len(paxitypes.VestingAccountPrefix):])
		addr, err := sdk.AccAddressFromBech32(addrStr)
		if err != nil {
			continue
		}

		if acc := k.accountKeeper.GetAccount(ctx, addr); acc != nil {
			if pva, ok := acc.(*vestingtypes.PeriodicVestingAccount); ok {
				startTime := pva.StartTime
				currentTime := startTime
				for _, period := range pva.VestingPeriods {
					currentTime += period.Length
					unlockSchedules = append(unlockSchedules, &paxitypes.UnlockSchedule{
						TimeUnix: time.Unix(currentTime, 0).Unix(),
						TimeStr:  time.Unix(currentTime, 0).Format("2006-01-02"),
						Address:  pva.Address,
						Amount:   period.Amount[0].Amount.Abs().Int64(),
						Denom:    period.Amount[0].Denom,
					})
				}
			}
		}
	}

	sort.Slice(unlockSchedules, func(i, j int) bool {
		return unlockSchedules[i].TimeUnix < unlockSchedules[j].TimeUnix
	})

	return unlockSchedules
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *paxitypes.GenesisState {
	params := k.GetParams(ctx)
	return &paxitypes.GenesisState{
		Params: params,
	}
}
