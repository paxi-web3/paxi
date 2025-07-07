package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func PositionKey(prc20 string, addr sdk.AccAddress) []byte {
	key := append([]byte(prc20), addr.Bytes()...)
	return append(PositionPrefix, key...)
}
