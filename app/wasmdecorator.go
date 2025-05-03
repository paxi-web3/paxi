package app

import (
	errorsmod "cosmossdk.io/errors"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type WasmDecorator struct {
	gasRegister wasmtypes.WasmGasRegisterConfig
	wasmKeeper  wasmkeeper.Keeper
}

func (w WasmDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	for _, msg := range tx.GetMsgs() {
		switch m := msg.(type) {
		case *wasmtypes.MsgStoreCode:
			// Increase the cost of storing code of smart contract
			codeLen := len(m.WASMByteCode)
			extraGas := uint64(codeLen) * w.gasRegister.CompileCost
			ctx.GasMeter().ConsumeGas(extraGas, "custom wasm store code penalty")
		case *wasmtypes.MsgInstantiateContract:
			// Only allow the uploader to deploy
			codeInfo := w.wasmKeeper.GetCodeInfo(ctx, m.CodeID)
			if codeInfo == nil {
				return ctx, errorsmod.Wrapf(sdkerrors.ErrNotFound, "code id %d not found", m.CodeID)
			}

			creatorAddr, err := sdk.AccAddressFromBech32(codeInfo.Creator)
			if err != nil {
				return ctx, errorsmod.Wrap(err, "invalid creator address")
			}

			sender := sdk.MustAccAddressFromBech32(m.Sender)
			if !sender.Equals(creatorAddr) {
				return ctx, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only uploader can instantiate this code")
			}
		}
	}
	return next(ctx, tx, simulate)
}
