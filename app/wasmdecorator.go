package app

import (
	errorsmod "cosmossdk.io/errors"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
			baseGas := uint64(120_000_000)
			extraGas := baseGas + uint64(codeLen)*w.gasRegister.CompileCost
			ctx.GasMeter().ConsumeGas(extraGas, "custom wasm store code penalty")
		case *wasmtypes.MsgInstantiateContract:
			// Increase the cost of instantiation of smart contract
			byteCode, err := w.wasmKeeper.GetByteCode(ctx, m.CodeID)
			if err != nil {
				return ctx, errorsmod.Wrapf(err, "failed to load bytecode for code id %d", m.CodeID)
			}
			codeSize := len(byteCode)
			baseGas := uint64(30_000_000)
			gasMultiplier := uint64(100)
			extraGas := baseGas + uint64(codeSize)*gasMultiplier
			ctx.GasMeter().ConsumeGas(extraGas, "custom wasm instantial penalty")
		}
	}
	return next(ctx, tx, simulate)
}
