package keeper

import (
	"fmt"

	storetypes "cosmossdk.io/core/store"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	customminttypes "github.com/paxi-web3/paxi/x/custommint/types"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	bankKeeper   bankkeeper.Keeper
	storeService storetypes.KVStoreService
	authority    string
	mintDenom    string
}

func NewKeeper(cdc codec.BinaryCodec, bankKeeper bankkeeper.Keeper, storeService storetypes.KVStoreService, authority string, denom string) Keeper {
	return Keeper{
		cdc:          cdc,
		bankKeeper:   bankKeeper,
		storeService: storeService,
		authority:    authority,
		mintDenom:    denom,
	}
}

func (k Keeper) BlockProvision(ctx sdk.Context) error {
	// Get the current block height
	blockHeight := ctx.BlockHeight()

	// Query the current total supply of the mint denom (e.g., upaxi)
	currentSupply := k.bankKeeper.GetSupply(ctx, k.mintDenom).Amount

	// Get the inflation rate based on block height (different each year)
	inflationRate := k.GetInflationRateByHeight(blockHeight)

	// Calculate provision for this block: (total_supply * inflation_rate) / blocks_per_year
	provision := sdkmath.LegacyNewDecFromInt(currentSupply).Mul(inflationRate).Quo(sdkmath.LegacyNewDec(customminttypes.BlocksPerYear))

	// Retrieve the accumulated minting amount from store
	acc := k.GetAccumulator(ctx)

	// Add this block's provision to the accumulator
	acc = acc.Add(provision)

	// If accumulated provision exceeds minting threshold, mint actual coins
	if acc.GTE(sdkmath.LegacyNewDec(customminttypes.MintThreshold)) {
		// Truncate accumulated provision to integer amount for minting
		mintAmount := provision.TruncateInt()

		// Subtract the minted amount from accumulator
		acc = acc.Sub(sdkmath.LegacyNewDecFromInt(mintAmount))

		// Create the mint coin object
		mintCoin := sdk.NewCoin(k.mintDenom, mintAmount)

		// Mint the coins into this module account
		if err := k.bankKeeper.MintCoins(ctx, customminttypes.ModuleName, sdk.NewCoins(mintCoin)); err != nil {
			return fmt.Errorf("mint failed: %w", err)
		}

		// Check if new total supply would exceed expected maximum supply
		expectedSupply := k.ExpectedSupplyByHeight(blockHeight)
		if currentSupply.Add(mintAmount).GT(expectedSupply) {
			// Calculate excess coins to be burned
			excess := currentSupply.Add(mintAmount).Sub(expectedSupply)

			// Burn the excess tokens directly from the module account
			burnCoin := sdk.NewCoin(k.mintDenom, excess)
			if err := k.bankKeeper.BurnCoins(ctx, customminttypes.ModuleName, sdk.NewCoins(burnCoin)); err != nil {
				return fmt.Errorf("burn failed: %w", err)
			}

			// Adjust minted amount after burn
			mintAmount = mintAmount.Sub(excess)

			// If nothing remains after burn, exit early
			if mintAmount.IsZero() {
				k.SetAccumulator(ctx, acc)
				return nil
			}
		}

		// Send the remaining minted coins to the distribution module for staking rewards
		distributeCoin := sdk.NewCoin(k.mintDenom, mintAmount)
		if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, customminttypes.ModuleName, distrtypes.ModuleName, sdk.NewCoins(distributeCoin)); err != nil {
			return fmt.Errorf("send to distribution failed: %w", err)
		}
	}

	// Store the updated accumulator back into KVStore
	k.SetAccumulator(ctx, acc)
	return nil
}

func (k Keeper) GetAccumulator(ctx sdk.Context) sdkmath.LegacyDec {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get([]byte(customminttypes.AccumulatorKey))
	if err != nil {
		return sdkmath.LegacyZeroDec()
	}
	dec := sdkmath.LegacyDec{}
	if err := dec.Unmarshal(bz); err != nil {
		panic(fmt.Errorf("failed to unmarshal accumulator: %w", err))
	}
	return dec
}

func (k Keeper) SetAccumulator(ctx sdk.Context, dec sdkmath.LegacyDec) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := dec.Marshal()
	if err != nil {
		panic(fmt.Errorf("failed to marshal accumulator: %w", err))
	}
	store.Set([]byte(customminttypes.AccumulatorKey), bz)
}

func (k Keeper) GetInflationRateByHeight(height int64) sdkmath.LegacyDec {
	switch {
	case height < customminttypes.BlocksPerYear:
		return sdkmath.LegacyMustNewDecFromStr("0.10") // Year 1
	case height < 2*customminttypes.BlocksPerYear:
		return sdkmath.LegacyMustNewDecFromStr("0.05") // Year 2
	default:
		return sdkmath.LegacyMustNewDecFromStr("0.025") // Year 3+
	}
}

func (k Keeper) ExpectedSupplyByHeight(height int64) sdkmath.Int {
	// You can customize this further if desired. Currently it's a simple projection.
	baseSupply := sdkmath.NewInt(customminttypes.TotalSupply)
	var growth sdkmath.LegacyDec

	switch {
	case height < customminttypes.BlocksPerYear:
		growth = sdkmath.LegacyMustNewDecFromStr("1.10")
	case height < 2*customminttypes.BlocksPerYear:
		growth = sdkmath.LegacyMustNewDecFromStr("1.155") // 10% first year, 5% second year
	default:
		years := sdkmath.LegacyNewDec(height).QuoInt64(customminttypes.BlocksPerYear)
		// Compound 2.5% inflation per year starting from year 3
		growth = sdkmath.LegacyMustNewDecFromStr("1.155").Mul(sdkmath.LegacyMustNewDecFromStr("1.025").Power(uint64(years.TruncateInt64() - 2)))
	}

	return growth.MulInt(baseSupply).TruncateInt()
}

func (k Keeper) InitGenesis(ctx sdk.Context, data *customminttypes.GenesisState) {
	err := k.SetParams(ctx, data)
	if err != nil {
		panic(err)
	}

	params, err := k.GetParams(ctx)
	fmt.Println("key value:", params)
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *customminttypes.GenesisState {
	params, err := k.GetParams(ctx)
	if err != nil {
		panic(err)
	}
	return params
}
