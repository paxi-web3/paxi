package prediction

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	_ "cosmossdk.io/api/cosmos/crypto/ed25519" // register to that it shows up in protoregistry.GlobalTypes
	"github.com/paxi-web3/paxi/x/prediction/types"
)

func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: types.Msg_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Use:       "update-params-proposal",
					Short:     "Submit a proposal to update prediction module params",
					Long:      "Submit a governance proposal to update prediction params.",
					Example:   "paxid tx prediction update-params-proposal '{\"max_batch_size\":500,\"create_market_bond\":\"20000000\",\"create_market_bond_denom\":\"uusdt\",\"market_fee_bps\":50,\"resolver_fee_share_percent\":80,\"max_order_lifetime_bh\":216000,\"max_open_orders_per_user\":1000,\"max_open_orders_per_market\":100,\"order_prune_interval_bh\":1000,\"order_prune_retain_bh\":21600,\"order_prune_scan_limit\":20000,\"order_prune_delete_limit\":10000}' --from key",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "params"},
					},
					GovProposal: true,
				},
			},
		},
	}
}
