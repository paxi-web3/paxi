package main

import (
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/paxi-web3/paxi/app"
	"github.com/spf13/cobra"

	sdkruntime "github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	sdkservertypes "github.com/cosmos/cosmos-sdk/server/types"
)

func main() {
	// Get the directory of the executable (not working dir)
	_, exePath, _, ok := runtime.Caller(0)
	if !ok {
		panic("failed to get executable path")
	}
	exeDir := filepath.Dir(exePath)
	homeDir := filepath.Join(exeDir, ".paxi")

	// Define the root command for the Paxi blockchain CLI
	rootCmd := &cobra.Command{
		Use:               "paxid",
		Short:             "Paxi Blockchain Node",
		PersistentPreRunE: server.PersistentPreRunEFn(serverconfig.DefaultConfig()),
	}

	// Register common Cosmos SDK CLI server commands (init, start, etc.)
	server.RegisterServerCommands(
		rootCmd,
		homeDir,      // Use executable path/.paxi as home
		newApp,       // Function to create the app
		exportApp,    // Function to export app state
		addInitFlags, // Additional init flags
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// newApp creates a new instance of the Paxi application
func newApp(
	logger sdkservertypes.Logger,
	db sdkservertypes.DB,
	traceStore io.Writer,
	appOpts sdkruntime.AppOptions,
) sdkservertypes.Application {
	return app.NewPaxiApp(logger, db, traceStore, true, map[int64]bool{}, appOpts)
}

// exportApp exports the application state and validator set
func exportApp(
	logger sdkservertypes.Logger,
	db sdkservertypes.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	appOpts sdkruntime.AppOptions,
) (sdkservertypes.ExportedApp, error) {
	paxiApp := app.NewPaxiApp(logger, db, traceStore, true, map[int64]bool{}, appOpts)
	return paxiApp.ExportAppStateAndValidators(forZeroHeight, nil)
}

// addInitFlags adds custom CLI flags like chain-id or home directory
func addInitFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("chain-id", "paxi-dev", "The ID of the chain to connect to")
	cmd.PersistentFlags().String("home", "", "The application home directory (default: same folder as executable/.paxi)")
}
