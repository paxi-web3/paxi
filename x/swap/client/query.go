package client

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/paxi-web3/paxi/x/swap/types"
	"github.com/spf13/cobra"
)

func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query params of swap module",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
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
		Use:   "swap",
		Short: "Querying commands for the swap module",
	}

	cmd.AddCommand(
		CmdQueryParams(),
	)

	return cmd
}
