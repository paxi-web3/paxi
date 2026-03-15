package types

import (
	"crypto/sha256"
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	MarketPrefix        = []byte{0x01}
	PositionPrefix      = []byte{0x02}
	NextMarketIDKey     = []byte{0x03}
	AppliedTradePrefix  = []byte{0x04}
	OrderPrefix         = []byte{0x05}
	NextOrderIDKey      = []byte{0x06}
	OpenOrderUserKey    = []byte{0x07}
	OpenOrderUserMktKey = []byte{0x08}
	OrderPruneCursorKey = []byte{0x09}
	OrderIDIndexPrefix  = []byte{0x0A}
)

func MarketStoreKey(marketID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, marketID)
	return append(MarketPrefix, bz...)
}

func PositionStoreKey(marketID uint64, addr sdk.AccAddress) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, marketID)

	key := make([]byte, 0, len(PositionPrefix)+len(bz)+len(addr.Bytes()))
	key = append(key, PositionPrefix...)
	key = append(key, bz...)
	key = append(key, addr.Bytes()...)
	return key
}

func OrderStoreKey(marketID uint64, orderID uint64) []byte {
	marketBz := make([]byte, 8)
	orderBz := make([]byte, 8)
	binary.BigEndian.PutUint64(marketBz, marketID)
	binary.BigEndian.PutUint64(orderBz, orderID)

	key := make([]byte, 0, len(OrderPrefix)+len(marketBz)+len(orderBz))
	key = append(key, OrderPrefix...)
	key = append(key, marketBz...)
	key = append(key, orderBz...)
	return key
}

func AppliedTradeStoreKey(marketID uint64, tradeID string) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, marketID)
	h := sha256.Sum256([]byte(tradeID))

	key := make([]byte, 0, len(AppliedTradePrefix)+len(bz)+len(h))
	key = append(key, AppliedTradePrefix...)
	key = append(key, bz...)
	key = append(key, h[:]...)
	return key
}

func OpenOrderCountByUserKey(addr sdk.AccAddress) []byte {
	key := make([]byte, 0, len(OpenOrderUserKey)+len(addr.Bytes()))
	key = append(key, OpenOrderUserKey...)
	key = append(key, addr.Bytes()...)
	return key
}

func OpenOrderCountByUserMarketKey(addr sdk.AccAddress, marketID uint64) []byte {
	marketBz := make([]byte, 8)
	binary.BigEndian.PutUint64(marketBz, marketID)

	key := make([]byte, 0, len(OpenOrderUserMktKey)+len(addr.Bytes())+len(marketBz))
	key = append(key, OpenOrderUserMktKey...)
	key = append(key, addr.Bytes()...)
	key = append(key, marketBz...)
	return key
}

func OrderIDIndexKey(orderID uint64) []byte {
	orderBz := make([]byte, 8)
	binary.BigEndian.PutUint64(orderBz, orderID)

	key := make([]byte, 0, len(OrderIDIndexPrefix)+len(orderBz))
	key = append(key, OrderIDIndexPrefix...)
	key = append(key, orderBz...)
	return key
}
