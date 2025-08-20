package keeper

import (
	"fmt"
	"math/big"

	store "cosmossdk.io/core/store"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	customminttypes "github.com/paxi-web3/paxi/x/custommint/types"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	bankKeeper   bankkeeper.Keeper
	storeService store.KVStoreService
	authority    string
	mintDenom    string
}

func NewKeeper(cdc codec.BinaryCodec, bankKeeper bankkeeper.Keeper, storeService store.KVStoreService, authority string, denom string) Keeper {
	return Keeper{
		cdc:          cdc,
		bankKeeper:   bankKeeper,
		storeService: storeService,
		authority:    authority,
		mintDenom:    denom,
	}
}

func (k Keeper) BlockProvision(ctx sdk.Context) error {
	// Get params
	params := k.GetParams(ctx)

	// Get the current block height
	blockHeight := ctx.BlockHeight()
	const mintThreshold int64 = 1 // 1 blocks per mint

	if blockHeight%mintThreshold != 0 {
		return nil
	}

	// Calculate provision for this block: (total_supply * inflation_rate) / (blocks_per_year / mint_threshold)
	currentSupply := k.bankKeeper.GetSupply(ctx, k.mintDenom).Amount
	inflationRate := k.GetInflationRateByHeight(params, blockHeight)
	provision := sdkmath.LegacyNewDecFromInt(currentSupply).
		Mul(inflationRate).
		Quo(sdkmath.LegacyNewDec(params.BlocksPerYear).Quo(sdkmath.LegacyNewDec(mintThreshold)))

	// Mint
	mintAmount := provision.TruncateInt()
	mintCoin := sdk.NewCoin(k.mintDenom, mintAmount)
	if err := k.bankKeeper.MintCoins(ctx, customminttypes.ModuleName, sdk.NewCoins(mintCoin)); err != nil {
		return fmt.Errorf("mint failed: %w", err)
	}

	// Get total minted amount from store
	totalMinted := k.GetTotalMinted(ctx)
	k.SetTotalMinted(ctx, totalMinted.Add(mintAmount))

	// Send the remaining minted coins to the distribution module for staking rewards
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, customminttypes.ModuleName, authtypes.FeeCollectorName, sdk.NewCoins(mintCoin)); err != nil {
		return fmt.Errorf("send to distribution failed: %w", err)
	}

	return nil
}

func (k Keeper) GetInflationRateByHeight(params customminttypes.Params, height int64) sdkmath.LegacyDec {
	switch {
	case height < params.BlocksPerYear:
		return params.FirstYearInflation // Year 1
	case height < 2*params.BlocksPerYear:
		return params.SecondYearInflation // Year 2
	default:
		return params.OtherYearInflation // Year 3+
	}
}

func (k Keeper) InitGenesis(ctx sdk.Context, data customminttypes.GenesisState) {
	k.SetParams(ctx, data.Params)
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *customminttypes.GenesisState {
	params := k.GetParams(ctx)
	return &customminttypes.GenesisState{
		Params: params,
	}
}

func (k Keeper) SetTotalMinted(ctx sdk.Context, amount sdkmath.Int) {
	store := k.storeService.OpenKVStore(ctx)
	bz := amount.BigInt().Bytes()
	store.Set([]byte(customminttypes.TotalMinted), bz)
}

func (k Keeper) GetTotalMinted(ctx sdk.Context) sdkmath.Int {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get([]byte(customminttypes.TotalMinted))
	if err != nil || bz == nil {
		return sdkmath.ZeroInt()
	}
	return sdkmath.NewIntFromBigInt(new(big.Int).SetBytes(bz))
}

func (k Keeper) SetTotalBurned(ctx sdk.Context, amount sdkmath.Int) {
	store := k.storeService.OpenKVStore(ctx)
	bz := amount.BigInt().Bytes()
	store.Set([]byte(customminttypes.TotalBurned), bz)
}

func (k Keeper) GetTotalBurned(ctx sdk.Context) sdkmath.Int {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get([]byte(customminttypes.TotalBurned))
	if err != nil || bz == nil {
		return sdkmath.ZeroInt()
	}
	return sdkmath.NewIntFromBigInt(new(big.Int).SetBytes(bz))
}

func (k Keeper) BurnExcessTokens(ctx sdk.Context) {
	// Get params
	params := k.GetParams(ctx)
	threshold := params.BurnThreshold // upaxi
	feeCollectorAddr := authtypes.NewModuleAddress(authtypes.FeeCollectorName)
	fees := k.bankKeeper.GetBalance(ctx, feeCollectorAddr, k.mintDenom)

	if fees.Amount.LTE(threshold) {
		return
	}

	// There is 50% chance to burn all balance from the fee pool
	if ctx.BlockHeight()%2 == 0 {
		fees.Amount = fees.Amount.ToLegacyDec().Mul(params.BurnRatio).RoundInt()
		err := k.bankKeeper.BurnCoins(ctx, authtypes.FeeCollectorName, sdk.NewCoins(fees))
		if err != nil {
			panic(err)
		}
		totalBurned := k.GetTotalBurned(ctx).Add(fees.Amount)
		k.SetTotalBurned(ctx, totalBurned)
		ctx.Logger().Info("Token burned", "amount", fees.Amount.String())
	}
}
