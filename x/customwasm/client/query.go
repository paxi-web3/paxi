package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/paxi-web3/paxi/x/customwasm/types"
	"github.com/spf13/cobra"
)

func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query params of wasm module",
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
		Use:   "custommint",
		Short: "Querying commands for the custommint module",
	}

	cmd.AddCommand(
		CmdQueryParams(),
	)

	return cmd
}
