package customstaking

var (
	MinBondedTokens int64 = 1_000_000_000 // 1,000 PAXI
	MaxCandidates   int   = 2000          // 2000 validators in the candidate list
	BlocksPerUpdate int64 = 200           // Blocks per update
	MinDelegation   int64 = 1_000_000     // 1 PAXI
	MinUndelegation int64 = 100_000       // 0.1 PAXI
)
