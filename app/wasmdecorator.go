package app

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type WasmDecorator struct {
	gasRegister wasmtypes.WasmGasRegisterConfig
}

func (w WasmDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// Increase the cost of storing code of smart contract
	if msg, ok := tx.GetMsgs()[0].(*wasmtypes.MsgStoreCode); ok {
		codeLen := len(msg.WASMByteCode)
		extraGas := uint64(codeLen) * w.gasRegister.CompileCost
		ctx.GasMeter().ConsumeGas(extraGas, "custom wasm store code penalty")
	}

	return next(ctx, tx, simulate)
}
