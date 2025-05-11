package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/paxi-web3/paxi/x/paxi/types"
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

			res, err := queryClient.LockedVesting(context.Background(), &types.QueryLockedVestingRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	return cmd
}

func CmdQueryCirculatingSupply() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "circulating-supply",
		Short: "Query circulating supply tokens",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.CirculatingSupply(context.Background(), &types.QueryCirculatingSupplyRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	return cmd
}

func CmdQueryTotalSupply() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-supply",
		Short: "Query total supply tokens",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.TotalSupply(context.Background(), &types.QueryTotalSupplyRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	return cmd
}

func CmdQueryEstimatedGasPrice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "estimated-gas-price",
		Short: "Query estimated gas price",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.EstimatedGasPrice(context.Background(), &types.QueryEstimatedGasPriceRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	return cmd
}

func CmdQueryLastBlockGasUsed() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "last-block-gas-used",
		Short: "Query gas used of the last block",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.LastBlockGasUsed(context.Background(), &types.QueryLastBlockGasUsedRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	return cmd
}

func CmdQueryTotalTxs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-txs",
		Short: "Query total txs that recorded",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.TotalTxs(context.Background(), &types.QueryTotalTxsRequest{})
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
		CmdQueryLastBlockGasUsed(),
		CmdQueryEstimatedGasPrice(),
		CmdQueryTotalSupply(),
		CmdQueryCirculatingSupply(),
	)

	return cmd
}
