package client

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/x/swap/types"
	"github.com/spf13/cobra"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "sawp",
		Short:              "Tx commands for the swap module",
		DisableFlagParsing: true,
		RunE:               client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdProvideLiquidity(),
	)

	return cmd
}

func CmdProvideLiquidity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "provide-liquidity [prc20] [paxi_amount] [prc20_amount]",
		Short: "Provide liquidity to a PAXI/PRC20 pool",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			prc20 := args[0]
			paxiAmountStr := args[1]
			prc20Amount := args[2]

			paxiAmount, err := sdk.ParseCoinNormalized(paxiAmountStr)
			if err != nil {
				return fmt.Errorf("invalid paxi amount: %w", err)
			}

			msg := &types.MsgProvideLiquidity{
				Creator:     clientCtx.GetFromAddress().String(),
				Prc20:       prc20,
				PaxiAmount:  &paxiAmount,
				Prc20Amount: prc20Amount,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
