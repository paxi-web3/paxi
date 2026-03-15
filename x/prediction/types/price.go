package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
)

const (
	// CollateralUnit is the base precision for upaxi/prc20 (6 decimals).
	CollateralUnit int64 = 1_000_000
	// PriceTickSize represents 0.01 in 6-decimal precision.
	PriceTickSize int64 = 10_000
	MinPriceTicks int64 = PriceTickSize
	MaxPriceTicks int64 = CollateralUnit
)

func ParsePriceTicks(value string, field string) (sdkmath.Int, error) {
	price, ok := sdkmath.NewIntFromString(value)
	if !ok {
		return sdkmath.Int{}, fmt.Errorf("invalid %s", field)
	}
	if !price.IsPositive() {
		return sdkmath.Int{}, fmt.Errorf("%s must be positive", field)
	}

	min := sdkmath.NewInt(MinPriceTicks)
	max := sdkmath.NewInt(MaxPriceTicks)
	if price.LT(min) || price.GT(max) {
		return sdkmath.Int{}, fmt.Errorf("%s must be between %d and %d", field, MinPriceTicks, MaxPriceTicks)
	}
	if !price.ModRaw(PriceTickSize).IsZero() {
		return sdkmath.Int{}, fmt.Errorf("%s must be a multiple of %d", field, PriceTickSize)
	}

	return price, nil
}
