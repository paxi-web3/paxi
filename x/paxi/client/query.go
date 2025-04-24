package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/paxi-web3/paxi/x/custommint/types"
	"github.com/spf13/cobra"
)

func CmdQueryLockedVesting() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locked-vesting",
		Short: "Query total locked tokens",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.TotalMinted(context.Background(), &types.QueryTotalMintedRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	return cmd
}

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "paxi",
		Short: "Querying commands for the paxi module",
	}

	cmd.AddCommand(
		CmdQueryLockedVesting(),
	)

	return cmd
}
