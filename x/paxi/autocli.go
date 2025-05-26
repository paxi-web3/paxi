package paxi

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	_ "cosmossdk.io/api/cosmos/crypto/ed25519" // register to that it shows up in protoregistry.GlobalTypes
	"github.com/paxi-web3/paxi/x/paxi/types"
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
					Short:     "Submit a proposal to update paxi module params",
					Long:      "Submit a proposal to update inflation, burn threshold, etc.",
					Example:   "paxid tx paxi update-params-proposal '{\"extra_gas_per_new_account\":\"400000\"}' --from key",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "params"},
					},
					GovProposal: true,
				},
			},
		},
	}
}
