package types

var (
	TotalSupply   = int64(100_000_000_000_000) // Genesis supply
	BlocksPerYear = int64(6307200)             // 365 days * 24 hours * 60 minutes * 60 seconds / 5 seconds per block
)

const (
	ModuleName   = "custommint"
	StoreKey     = ModuleName
	TotalMinted  = "total_minted"
	DefaultDenom = "upaxi"
)
