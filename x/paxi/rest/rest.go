package rest

import (
	"context"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	paxitypes "github.com/paxi-web3/paxi/x/paxi/types"
)

func LockedVestingHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryClient := paxitypes.NewQueryClient(clientCtx)

		// Call gRPC Query client
		res, err := queryClient.LockedVesting(context.Background(), &paxitypes.QueryLockedVestingRequest{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Output JSON
		w.Header().Set("Content-Type", "application/json")
		bz, err := clientCtx.Codec.MarshalJSON(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(bz)
	}
}

func CirculatingSupplyHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryClient := paxitypes.NewQueryClient(clientCtx)

		// Call gRPC Query client
		res, err := queryClient.CirculatingSupply(context.Background(), &paxitypes.QueryCirculatingSupplyRequest{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Output JSON
		w.Header().Set("Content-Type", "application/json")
		bz, err := clientCtx.Codec.MarshalJSON(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(bz)
	}
}

func TotalSupplyHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryClient := paxitypes.NewQueryClient(clientCtx)

		// Call gRPC Query client
		res, err := queryClient.TotalSupply(context.Background(), &paxitypes.QueryTotalSupplyRequest{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Output JSON
		w.Header().Set("Content-Type", "application/json")
		bz, err := clientCtx.Codec.MarshalJSON(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(bz)
	}
}

func EstimatedGasPriceHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryClient := paxitypes.NewQueryClient(clientCtx)

		// Call gRPC Query client
		res, err := queryClient.EstimatedGasPrice(context.Background(), &paxitypes.QueryEstimatedGasPriceRequest{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Output JSON
		w.Header().Set("Content-Type", "application/json")
		bz, err := clientCtx.Codec.MarshalJSON(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(bz)
	}
}

func LastBlockGasUsedHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryClient := paxitypes.NewQueryClient(clientCtx)

		// Call gRPC Query client
		res, err := queryClient.LastBlockGasUsed(context.Background(), &paxitypes.QueryLastBlockGasUsedRequest{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Output JSON
		w.Header().Set("Content-Type", "application/json")
		bz, err := clientCtx.Codec.MarshalJSON(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(bz)
	}
}
