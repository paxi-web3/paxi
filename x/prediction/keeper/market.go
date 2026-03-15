package keeper

import (
	"encoding/binary"
	"fmt"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/x/prediction/types"
)

func (k Keeper) GetMarket(ctx sdk.Context, marketID uint64) (*types.Market, bool) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.MarketStoreKey(marketID))
	if err != nil || bz == nil {
		return nil, false
	}

	market := &types.Market{}
	if err := k.cdc.Unmarshal(bz, market); err != nil {
		return nil, false
	}

	return market, true
}

func (k Keeper) SetMarket(ctx sdk.Context, market *types.Market) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := k.cdc.Marshal(market)
	if err != nil {
		panic(err)
	}
	k.mustSet(store, types.MarketStoreKey(market.Id), bz)
}

func (k Keeper) GetAllMarkets(ctx sdk.Context) []*types.Market {
	store := k.storeService.OpenKVStore(ctx)
	it, err := prefixIterator(store, types.MarketPrefix)
	if err != nil {
		panic(err)
	}
	defer it.Close()

	markets := make([]*types.Market, 0)
	for ; it.Valid(); it.Next() {
		market := &types.Market{}
		if err := k.cdc.Unmarshal(it.Value(), market); err != nil {
			continue
		}
		markets = append(markets, market)
	}
	return markets
}

func (k Keeper) SetNextMarketID(ctx sdk.Context, nextID uint64) {
	if nextID == 0 {
		nextID = 1
	}
	store := k.storeService.OpenKVStore(ctx)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, nextID)
	k.mustSet(store, types.NextMarketIDKey, bz)
}

func (k Keeper) GetNextMarketID(ctx sdk.Context) uint64 {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.NextMarketIDKey)
	if err != nil || len(bz) != 8 {
		return 1
	}
	id := binary.BigEndian.Uint64(bz)
	if id == 0 {
		return 1
	}
	return id
}

func (k Keeper) NextMarketID(ctx sdk.Context) uint64 {
	id := k.GetNextMarketID(ctx)
	k.SetNextMarketID(ctx, id+1)
	return id
}

func (k Keeper) ValidateMarketInvariants(m *types.Market) error {
	yes, err := parseNonNegativeInt(m.TotalYesShares, "total_yes_shares")
	if err != nil {
		return err
	}
	no, err := parseNonNegativeInt(m.TotalNoShares, "total_no_shares")
	if err != nil {
		return err
	}
	if yes.IsNegative() || no.IsNegative() {
		return fmt.Errorf("market share totals cannot be negative")
	}

	return nil
}

func (k Keeper) maybeCloseExpiredMarket(ctx sdk.Context, market *types.Market) bool {
	if market.Status != types.MarketStatus_MARKET_STATUS_OPEN {
		return false
	}
	if ctx.BlockTime().Unix() < market.CloseTime {
		return false
	}
	market.Status = types.MarketStatus_MARKET_STATUS_CLOSED
	k.SetMarket(ctx, market)
	return true
}

func (k Keeper) AutoCloseExpiredMarkets(ctx sdk.Context) {
	store := k.storeService.OpenKVStore(ctx)
	it, err := prefixIterator(store, types.MarketPrefix)
	if err != nil {
		panic(err)
	}
	defer it.Close()

	toClose := make([]uint64, 0)

	for ; it.Valid(); it.Next() {
		market := &types.Market{}
		if err := k.cdc.Unmarshal(it.Value(), market); err != nil {
			continue
		}
		if market.Status == types.MarketStatus_MARKET_STATUS_OPEN && ctx.BlockTime().Unix() >= market.CloseTime {
			toClose = append(toClose, market.Id)
		}
	}

	for _, id := range toClose {
		market, found := k.GetMarket(ctx, id)
		if !found {
			continue
		}
		k.maybeCloseExpiredMarket(ctx, market)
	}
}

func normalizeOutcome(outcome string) string {
	return strings.ToUpper(strings.TrimSpace(outcome))
}

func parseOutcome(outcome types.Outcome) (string, error) {
	switch outcome {
	case types.Outcome_OUTCOME_YES:
		return "YES", nil
	case types.Outcome_OUTCOME_NO:
		return "NO", nil
	default:
		return "", types.ErrInvalidOutcome
	}
}

func intToStr(i uint64) string {
	return fmt.Sprintf("%d", i)
}

func (k Keeper) getMarketShareInts(m *types.Market) (yes sdkmath.Int, no sdkmath.Int, err error) {
	yes, err = parseNonNegativeInt(m.TotalYesShares, "total_yes_shares")
	if err != nil {
		return sdkmath.Int{}, sdkmath.Int{}, err
	}
	no, err = parseNonNegativeInt(m.TotalNoShares, "total_no_shares")
	if err != nil {
		return sdkmath.Int{}, sdkmath.Int{}, err
	}
	return yes, no, nil
}
