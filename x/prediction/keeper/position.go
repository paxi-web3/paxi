package keeper

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/x/prediction/types"
)

func (k Keeper) GetPosition(ctx sdk.Context, marketID uint64, addr sdk.AccAddress) (*types.Position, bool) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.PositionStoreKey(marketID, addr))
	if err != nil || bz == nil {
		return nil, false
	}

	position := &types.Position{}
	if err := k.cdc.Unmarshal(bz, position); err != nil {
		return nil, false
	}

	return position, true
}

func (k Keeper) SetPosition(ctx sdk.Context, position *types.Position) {
	store := k.storeService.OpenKVStore(ctx)
	addr := sdk.MustAccAddressFromBech32(position.Address)
	bz, err := k.cdc.Marshal(position)
	if err != nil {
		panic(err)
	}
	k.mustSet(store, types.PositionStoreKey(position.MarketId, addr), bz)
}

func (k Keeper) GetAllPositions(ctx sdk.Context) []*types.Position {
	store := k.storeService.OpenKVStore(ctx)
	it, err := prefixIterator(store, types.PositionPrefix)
	if err != nil {
		panic(err)
	}
	defer it.Close()

	positions := make([]*types.Position, 0)
	for ; it.Valid(); it.Next() {
		pos := &types.Position{}
		if err := k.cdc.Unmarshal(it.Value(), pos); err != nil {
			continue
		}
		positions = append(positions, pos)
	}

	return positions
}

func (k Keeper) mustPositionInts(pos *types.Position) (
	yes sdkmath.Int,
	lockedYes sdkmath.Int,
	no sdkmath.Int,
	lockedNo sdkmath.Int,
	err error,
) {
	yes, err = parseNonNegativeInt(pos.YesShares, "yes_shares")
	if err != nil {
		return
	}
	lockedYes, err = parseNonNegativeInt(pos.LockedYesShares, "locked_yes_shares")
	if err != nil {
		return
	}
	no, err = parseNonNegativeInt(pos.NoShares, "no_shares")
	if err != nil {
		return
	}
	lockedNo, err = parseNonNegativeInt(pos.LockedNoShares, "locked_no_shares")
	if err != nil {
		return
	}

	if err = ensureNoNegative(yes, lockedYes, no, lockedNo); err != nil {
		return
	}
	if lockedYes.GT(yes) {
		err = fmt.Errorf("locked_yes_shares cannot exceed yes_shares")
		return
	}
	if lockedNo.GT(no) {
		err = fmt.Errorf("locked_no_shares cannot exceed no_shares")
		return
	}

	return
}

func (k Keeper) mustSetPositionInts(
	pos *types.Position,
	yes sdkmath.Int,
	lockedYes sdkmath.Int,
	no sdkmath.Int,
	lockedNo sdkmath.Int,
) {
	pos.YesShares = yes.String()
	pos.LockedYesShares = lockedYes.String()
	pos.NoShares = no.String()
	pos.LockedNoShares = lockedNo.String()
}

func (k Keeper) assertPositionInvariant(pos *types.Position) error {
	_, _, _, _, err := k.mustPositionInts(pos)
	return err
}
