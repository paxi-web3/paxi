package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type SpamPreventDecorator struct {
	ak  authkeeper.AccountKeeper
	App *PaxiApp
}

func (s SpamPreventDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	const extraGasPerNewAccount = 300_000

	msgs := tx.GetMsgs()
	for _, msg := range msgs {
		switch m := msg.(type) {
		case *banktypes.MsgSend:
			// Add extra gas if the account does not exists
			addr := sdk.MustAccAddressFromBech32(m.ToAddress)
			if !s.ak.HasAccount(ctx, addr) {
				ctx.GasMeter().ConsumeGas(extraGasPerNewAccount, "new account creation penalty")
			}
		case *banktypes.MsgMultiSend:
			// Add extra gas if the account does not exists
			for _, output := range m.Outputs {
				addr := sdk.MustAccAddressFromBech32(output.Address)
				if !s.ak.HasAccount(ctx, addr) {
					ctx.GasMeter().ConsumeGas(extraGasPerNewAccount, "new account in MsgMultiSend")
				}
			}
		}
	}

	return next(ctx, tx, simulate)
}
