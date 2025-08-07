package cli

import (
	"fmt"
	"math"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/x/paxi/types"
	"github.com/spf13/cobra"
)

func CmdBurnToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn-token [amount]",
		Short: "Burn your own tokens",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender := clientCtx.GetFromAddress()

			// Parse the coin string (e.g. "1000000upaxi")
			amount, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			// Prevent panic when converting to uint64
			if amount.Amount.GT(sdkmath.NewIntFromUint64(math.MaxUint64)) {
				return fmt.Errorf("amount too large: must fit in uint64")
			}
			// Prevent exceeding the maximum total token supply
			maxSupply := sdkmath.NewInt(1_000_000_000_000_000) // example cap: 1e15 units
			if amount.Amount.GT(maxSupply) {
				return fmt.Errorf("amount exceeds maximum supply (%s)", maxSupply.String())
			}

			msg := &types.MsgBurnToken{
				Sender: sender.String(),
				Amount: []*sdk.Coin{&amount},
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "paxi",
		Short:              "Tx commands for the paxi module",
		DisableFlagParsing: true,
		RunE:               client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdBurnToken(),
	)

	return cmd
}
