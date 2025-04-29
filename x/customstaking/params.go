package customstaking

var (
	MinBonedTokens int64 = 1_000_000_000 // 1,000 PAXI
	MaxCandidates  int   = 5000          // 5000 validators in the candidate list
	//BlocksPerUpdate int64 = 1000          // Blocks per update
	BlocksPerUpdate int64 = 2 // Blocks per update for testing
)
