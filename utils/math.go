package utils

import "math/bits"

// IntSqrt returns floor(sqrt(x)) deterministically with pure integer math.
func IntSqrt(x uint64) uint64 {
	if x == 0 || x == 1 {
		return x
	}
	// Newton's method for integer sqrt (monotone decreasing to floor sqrt)
	r := uint64(1) << ((bits.Len64(x) + 1) / 2) // initial guess ~ 2^(ceil(log2(sqrt(x))))
	for {
		nr := (r + x/r) >> 1
		if nr >= r {
			return r
		}
		r = nr
	}
}
