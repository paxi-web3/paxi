package types

var (
	BlocksPerYear       = int64(6307200) // 365 days * 24 hours * 60 minutes * 60 seconds / 5 seconds per block
	FirstYearInflation  = float32(0.08)
	SecondYearInflation = float32(0.05)
	OtherYearInflation  = float32(0.025)
)

const (
	ModuleName   = "custommint"
	StoreKey     = ModuleName
	TotalMinted  = "total_minted"
	DefaultDenom = "upaxi"
)
