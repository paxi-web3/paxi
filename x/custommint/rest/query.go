package rest

import (
	"context"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	customminttypes "github.com/paxi-web3/paxi/x/custommint/types"
)

func TotalMintedHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryClient := customminttypes.NewQueryClient(clientCtx)

		// Call gRPC Query client
		res, err := queryClient.TotalMinted(context.Background(), &customminttypes.QueryTotalMintedRequest{})
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
