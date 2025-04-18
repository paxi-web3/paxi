// Paxi Blockchain Main Entry
package main

import (
	"fmt"
	"os"

	"github.com/paxi-web3/paxi/app"
	paxicmd "github.com/paxi-web3/paxi/cmd"

	clientv2helpers "cosmossdk.io/client/v2/helpers"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {

	rootCmd := paxicmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, clientv2helpers.EnvPrefix, app.DefaultNodeHome); err != nil {
		fmt.Fprintln(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}
}
