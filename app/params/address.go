package params

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	Bech32PrefixAccAddr  = "paxi"
	Bech32PrefixAccPub   = "paxipub"
	Bech32PrefixValAddr  = "paxivaloper"
	Bech32PrefixValPub   = "paxivaloperpub"
	Bech32PrefixConsAddr = "paxivalcons"
	Bech32PrefixConsPub  = "paxivalconspub"
)

func InitAddressRules() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)

	cfg := sdk.GetConfig()
	cfg.SetCoinType(118)
	cfg.SetPurpose(44)

	cfg.Seal()
}
