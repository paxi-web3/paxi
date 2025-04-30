package customstaking

var (
	MinBondedTokens int64 = 1_000_000_000 // 1,000 PAXI
	MaxCandidates   int   = 2000          // 2000 validators in the candidate list
	BlocksPerUpdate int64 = 200           // Blocks per update
)
