package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func PositionKey(prc20 string, addr sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("position:%s:%s", prc20, addr.String()))
}
