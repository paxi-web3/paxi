package custommint

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	_ "cosmossdk.io/api/cosmos/crypto/ed25519" // register to that it shows up in protoregistry.GlobalTypes
	"github.com/paxi-web3/paxi/x/custommint/types"
)

func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	// Enables the 'params' field to appear in the proposal draft JSON,
	// so users can manually fill in module parameters during proposal creation.
	return &autocliv1.ModuleOptions{
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: types.Msg_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Use:       "update-params-proposal",
					Short:     "Submit a proposal to update custommint module params",
					Long:      "Submit a proposal to update inflation, burn threshold, etc.",
					Example:   "paxid tx custommint update-params-proposal '{\"first_year_inflation\":\"0.08\"}' --from key",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "params"},
					},
					GovProposal: true,
				},
			},
		},
	}
}
