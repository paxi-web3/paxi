package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/paxi-web3/paxi/x/custommint/types"
	"github.com/spf13/cobra"
)

func CmdQueryTotalMinted() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-minted",
		Short: "Query total minted tokens",
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
		Use:   "custommint",
		Short: "Querying commands for the custommint module",
	}

	cmd.AddCommand(
		CmdQueryTotalMinted(),
	)

	return cmd
}
