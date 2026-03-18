package keeper

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/utils"
	"github.com/paxi-web3/paxi/x/prediction/types"
)

const marketCloseIndexBackfillScanLimit = uint64(5000)

type marketCloseIndexEntry struct {
	key       []byte
	closeTime int64
	marketID  uint64
}

func marketCloseIndexKeyForMarket(market *types.Market) ([]byte, bool) {
	if market.Status != types.MarketStatus_MARKET_STATUS_OPEN || market.CloseTime <= 0 {
		return nil, false
	}
	return types.MarketCloseIndexKey(market.CloseTime, market.Id), true
}

func parseMarketCloseIndexEntry(key []byte) (closeTime int64, marketID uint64, ok bool) {
	prefixLen := len(types.MarketCloseIdxPrefix)
	if len(key) != prefixLen+16 || !bytes.HasPrefix(key, types.MarketCloseIdxPrefix) {
		return 0, 0, false
	}
	closeTimeU := binary.BigEndian.Uint64(key[prefixLen : prefixLen+8])
	if closeTimeU > uint64(math.MaxInt64) {
		return 0, 0, false
	}
	closeTime = int64(closeTimeU)
	marketID = binary.BigEndian.Uint64(key[prefixLen+8 : prefixLen+16])
	return closeTime, marketID, true
}

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
	var oldIndexKey []byte
	oldIndexSet := false

	oldBz, err := store.Get(types.MarketStoreKey(market.Id))
	if err != nil {
		panic(err)
	}
	if oldBz != nil {
		oldMarket := &types.Market{}
		if err := k.cdc.Unmarshal(oldBz, oldMarket); err != nil {
			panic(err)
		}
		oldIndexKey, oldIndexSet = marketCloseIndexKeyForMarket(oldMarket)
	}

	newIndexKey, newIndexSet := marketCloseIndexKeyForMarket(market)
	if oldIndexSet && (!newIndexSet || !bytes.Equal(oldIndexKey, newIndexKey)) {
		k.mustDelete(store, oldIndexKey)
	}

	bz, err := k.cdc.Marshal(market)
	if err != nil {
		panic(err)
	}
	k.mustSet(store, types.MarketStoreKey(market.Id), bz)
	if newIndexSet && (!oldIndexSet || !bytes.Equal(oldIndexKey, newIndexKey)) {
		k.mustSet(store, newIndexKey, []byte{1})
	}
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
	k.backfillMarketCloseIndexStep(ctx)

	store := k.storeService.OpenKVStore(ctx)
	it, err := prefixIterator(store, types.MarketCloseIdxPrefix)
	if err != nil {
		panic(err)
	}
	defer it.Close()

	now := ctx.BlockTime().Unix()
	toClose := make([]marketCloseIndexEntry, 0)

	for ; it.Valid(); it.Next() {
		key := append([]byte(nil), it.Key()...)
		closeTime, marketID, ok := parseMarketCloseIndexEntry(key)
		if !ok {
			k.mustDelete(store, key)
			continue
		}
		if closeTime > now {
			break
		}
		toClose = append(toClose, marketCloseIndexEntry{
			key:       key,
			closeTime: closeTime,
			marketID:  marketID,
		})
	}

	for i := range toClose {
		entry := toClose[i]
		market, found := k.GetMarket(ctx, entry.marketID)
		if !found {
			k.mustDelete(store, entry.key)
			continue
		}
		if market.Status != types.MarketStatus_MARKET_STATUS_OPEN {
			k.mustDelete(store, entry.key)
			continue
		}
		if market.CloseTime != entry.closeTime {
			// Refresh stale index entry when close_time changed.
			k.mustDelete(store, entry.key)
			if updatedKey, ok := marketCloseIndexKeyForMarket(market); ok {
				k.mustSet(store, updatedKey, []byte{1})
			}
		}
		k.maybeCloseExpiredMarket(ctx, market)
	}
}

func (k Keeper) backfillMarketCloseIndexStep(ctx sdk.Context) {
	store := k.storeService.OpenKVStore(ctx)
	ready, err := store.Get(types.MarketCloseIdxReady)
	if err != nil {
		panic(err)
	}
	if len(ready) > 0 {
		return
	}

	cursor, err := store.Get(types.MarketCloseIdxCursor)
	if err != nil {
		panic(err)
	}

	start, end := utils.PrefixRange(types.MarketPrefix)
	if len(cursor) > 0 {
		start = cursor
	}
	it, err := store.Iterator(start, end)
	if err != nil {
		panic(err)
	}
	defer it.Close()

	var scanned uint64
	var lastScannedKey []byte

	for ; it.Valid() && scanned < marketCloseIndexBackfillScanLimit; it.Next() {
		key := append([]byte(nil), it.Key()...)
		if len(cursor) > 0 && bytes.Equal(key, cursor) {
			continue
		}

		lastScannedKey = key
		scanned++

		market := &types.Market{}
		if err := k.cdc.Unmarshal(it.Value(), market); err != nil {
			continue
		}
		if indexKey, ok := marketCloseIndexKeyForMarket(market); ok {
			k.mustSet(store, indexKey, []byte{1})
		}
	}

	if !it.Valid() {
		k.mustDelete(store, types.MarketCloseIdxCursor)
		k.mustSet(store, types.MarketCloseIdxReady, []byte{1})
		return
	}
	if len(lastScannedKey) > 0 {
		k.mustSet(store, types.MarketCloseIdxCursor, lastScannedKey)
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
