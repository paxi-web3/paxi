package utils

import (
	"math/big"
	"math/bits"

	sdkmath "cosmossdk.io/math"
)

// UintSqrt returns floor(sqrt(x)) deterministically with pure integer math.
func UintSqrt(x uint64) uint64 {
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

// IntSqrt returns floor(sqrt(x)) using pure integer math on sdkmath.Int.
// Panics if x < 0.
func IntSqrt(x sdkmath.Int) sdkmath.Int {
	if x.IsNegative() {
		panic("IntSqrt: negative input")
	}
	if x.IsZero() || x.Equal(sdkmath.NewInt(1)) {
		return x
	}

	// Work on a copy of x's big.Int
	xb := new(big.Int).Set(x.BigInt())

	// Initial guess r ≈ 2^ceil(bitlen(x)/2)
	n := xb.BitLen()
	r := new(big.Int).Lsh(big.NewInt(1), uint((n+1)/2))

	// Newton iteration: r_{k+1} = floor((r_k + x / r_k) / 2)
	// This sequence is monotonically non-increasing and converges to floor(sqrt(x)).
	t := new(big.Int)
	q := new(big.Int)
	for {
		q.Quo(xb, r) // q = x / r
		t.Add(r, q)  // t = r + x/r
		t.Rsh(t, 1)  // t = (r + x/r) / 2

		if t.Cmp(r) >= 0 {
			// Converged: r is floor(sqrt(x))
			break
		}
		r.Set(t)
	}

	return sdkmath.NewIntFromBigInt(r)
}

type Rng64 uint64

func (r *Rng64) Next() uint64 {
	// splitmix64
	x := uint64(*r) + 0x9e3779b97f4a7c15
	*r = Rng64(x)
	z := x
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}
func (r *Rng64) Int63n(n int64) int64 { return int64(r.Next() % uint64(n)) }
