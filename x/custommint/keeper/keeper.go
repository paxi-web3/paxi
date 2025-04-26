package keeper

import (
	"fmt"
	"math/big"
	"strconv"

	storetypes "cosmossdk.io/core/store"
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
	// Get genesis state
	genesisState := k.ExportGenesis(ctx)

	// Get the current block height
	blockHeight := ctx.BlockHeight()
	const mintThreshold int64 = 1 // 1 blocks per mint

	if blockHeight%mintThreshold != 0 {
		return nil
	}

	// Calculate provision for this block: (total_supply * inflation_rate) / (blocks_per_year / mint_threshold)
	currentSupply := k.bankKeeper.GetSupply(ctx, k.mintDenom).Amount
	inflationRate := k.GetInflationRateByHeight(genesisState, blockHeight)
	provision := sdkmath.LegacyNewDecFromInt(currentSupply).
		Mul(inflationRate).
		Quo(sdkmath.LegacyNewDec(genesisState.BlocksPerYear).Quo(sdkmath.LegacyNewDec(mintThreshold)))

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

func (k Keeper) GetInflationRateByHeight(genesisState *customminttypes.GenesisState, height int64) sdkmath.LegacyDec {

	switch {
	case height < genesisState.BlocksPerYear:
		return sdkmath.LegacyMustNewDecFromStr(strconv.FormatFloat(float64(genesisState.FirstYearInflation), 'f', -1, 32)) // Year 1
	case height < 2*genesisState.BlocksPerYear:
		return sdkmath.LegacyMustNewDecFromStr(strconv.FormatFloat(float64(genesisState.SecondYearInflation), 'f', -1, 32)) // Year 2
	default:
		return sdkmath.LegacyMustNewDecFromStr(strconv.FormatFloat(float64(genesisState.OtherYearInflation), 'f', -1, 32)) // Year 3+
	}
}

func (k Keeper) InitGenesis(ctx sdk.Context, data customminttypes.GenesisState) {
	err := k.SetParams(ctx, &data)
	if err != nil {
		panic(err)
	}
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *customminttypes.GenesisState {
	params, err := k.GetParams(ctx)
	if err != nil {
		panic(err)
	}
	return &params
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
