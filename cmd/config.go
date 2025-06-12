package cmd

import (
	time "time"

	cmtcfg "github.com/cometbft/cometbft/config"
	cmttypes "github.com/cometbft/cometbft/types"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	apptypes "github.com/paxi-web3/paxi/app/types"
)

// initCometBFTConfig helps to override default CometBFT Config values.
// return cmtcfg.DefaultConfig if no custom configuration is required for the application.
func initCometBFTConfig() *cmtcfg.Config {
	cfg := cmtcfg.DefaultConfig()

	// P2P networking settings
	cfg.P2P.MaxNumInboundPeers = 80                      // Max inbound peers allowed to connect
	cfg.P2P.MaxNumOutboundPeers = 30                     // Max outbound peers this node will dial
	cfg.P2P.FlushThrottleTimeout = 50 * time.Millisecond // Delay between sending messages
	cfg.P2P.SendRate = 100 * 1024 * 1024                 // Max bytes per second to send (100MB/s)
	cfg.P2P.RecvRate = 100 * 1024 * 1024                 // Max bytes per second to receive (100MB/s)
	cfg.P2P.PexReactor = true                            // Enable peer exchange
	cfg.P2P.AllowDuplicateIP = false                     // Avoid duplicate IP connections

	// RPC server settings
	cfg.RPC.CORSAllowedOrigins = []string{}                                   // Restrict CORS origins
	cfg.RPC.CORSAllowedMethods = []string{"HEAD", "GET", "POST"}              // Allow these HTTP methods
	cfg.RPC.CORSAllowedHeaders = []string{"Origin", "Accept", "Content-Type"} // CORS headers
	cfg.RPC.MaxBodyBytes = 1_000_000                                          // Max request body size in bytes
	cfg.RPC.MaxHeaderBytes = 1_048_576                                        // Max request header size in bytes

	// Consensus configuration
	cfg.Consensus.TimeoutPropose = 3 * time.Second // Timeout for proposing a block
	cfg.Consensus.TimeoutProposeDelta = 400 * time.Millisecond
	cfg.Consensus.TimeoutPrevote = 1200 * time.Millisecond
	cfg.Consensus.TimeoutPrevoteDelta = 400 * time.Millisecond
	cfg.Consensus.TimeoutPrecommit = 1200 * time.Millisecond
	cfg.Consensus.TimeoutPrecommitDelta = 400 * time.Millisecond
	cfg.Consensus.TimeoutCommit = 4000 * time.Millisecond      // Timeout before getting into next block
	cfg.Consensus.CreateEmptyBlocks = true                     // Create blocks even with no transactions
	cfg.Consensus.CreateEmptyBlocksInterval = 10 * time.Second // Time between empty blocks
	cfg.Consensus.PeerGossipSleepDuration = 100 * time.Millisecond
	cfg.Consensus.PeerQueryMaj23SleepDuration = 1 * time.Second

	// Mempool configuration
	cfg.Mempool.Size = 10000                    // Max txs in mempool
	cfg.Mempool.MaxTxsBytes = 512 * 1024 * 1024 // 512 MB mempool capacity
	cfg.Mempool.MaxTxBytes = 3 * 1024 * 1024    // 3 MB per tx
	cfg.Mempool.CacheSize = 10000               // Tx cache size
	cfg.Mempool.Recheck = true                  // Recheck mempool on reorg
	cfg.Mempool.Broadcast = true                // Enable gossip broadcast

	// Database backend for storing blockchain data
	cfg.BaseConfig.DBBackend = "pebbledb" // Options: goleveldb, rocksdb, pebbledb, etc.

	return cfg
}

// initAppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func initAppConfig() (string, interface{}) {
	// Define a custom config structure to inject additional metadata into app.toml
	type CustomConfig struct {
		ChainTitle string `mapstructure:"chain-title"` // Name of the chain (custom field)
	}

	// Combine SDK server config with custom fields
	type CustomAppConfig struct {
		serverconfig.Config `mapstructure:",squash"` // Base Cosmos SDK config
		Custom              CustomConfig             `mapstructure:"custom"`
	}

	// Start from default SDK configuration
	srvCfg := serverconfig.DefaultConfig()

	// Set a minimum gas price (required for node startup)
	// This avoids the validator node halting due to missing gas price
	srvCfg.MinGasPrices = "0.05" + apptypes.DefaultDenom
	srvCfg.QueryGasLimit = 500000 // Set a reasonable gas limit for queries

	// Pruning
	srvCfg.Pruning = "custom"
	srvCfg.PruningKeepRecent = "40000"
	srvCfg.PruningInterval = "100"

	// Enable essential APIs and endpoints
	srvCfg.API.Enable = true
	srvCfg.API.Swagger = true
	srvCfg.GRPC.Enable = true
	srvCfg.GRPCWeb.Enable = true

	// Optional: disable telemetry unless needed
	srvCfg.Telemetry.Enabled = false
	srvCfg.Telemetry.PrometheusRetentionTime = 120

	// Enable snapshots to support state sync and archiving
	srvCfg.StateSync.SnapshotInterval = 0
	srvCfg.StateSync.SnapshotKeepRecent = 5

	// IavlCacheSize set the size of the iavl tree cache (in number of nodes).
	srvCfg.IAVLCacheSize = 2_000_000

	// Set default values for the custom configuration section
	customAppConfig := CustomAppConfig{
		Config: *srvCfg,
		Custom: CustomConfig{
			ChainTitle: "Paxi", // Set your own chain title here
		},
	}

	// Append custom config template to the default app.toml template
	customAppTemplate := serverconfig.DefaultConfigTemplate + `
[custom]
# The title of your blockchain, used for display or branding.
chain-title = "{{ .Custom.ChainTitle }}"`

	return customAppTemplate, customAppConfig
}

// Custom default consensus
func CustomDefaultConsensusParams() *cmttypes.ConsensusParams {
	cp := cmttypes.DefaultConsensusParams()
	cp.Block = cmttypes.BlockParams{
		MaxBytes: 5 * 1024 * 1024,
		MaxGas:   4_000_000_000,
	}
	return cp
}
