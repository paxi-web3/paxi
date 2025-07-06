package utils

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
)

func ValidatePositiveDecString(s string) (sdkmath.LegacyDec, error) {
	val, err := sdkmath.LegacyNewDecFromStr(s)
	if err != nil {
		return sdkmath.LegacyDec{}, fmt.Errorf("invalid decimal string: %w", err)
	}
	if val.IsNegative() || val.IsZero() {
		return sdkmath.LegacyDec{}, fmt.Errorf("decimal must be positive")
	}
	return val, nil
}

func ValidatePositiveIntString(s string) (sdkmath.Int, error) {
	val, ok := sdkmath.NewIntFromString(s)
	if !ok {
		return sdkmath.Int{}, fmt.Errorf("invalid decimal string")
	}
	if val.IsNegative() || val.IsZero() {
		return sdkmath.Int{}, fmt.Errorf("decimal must be positive")
	}
	return val, nil
}
