package app

import (
	errorsmod "cosmossdk.io/errors"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	customwasmkeeper "github.com/paxi-web3/paxi/x/customwasm/keeper"
)

type WasmDecorator struct {
	gasRegister      wasmtypes.WasmGasRegisterConfig
	wasmKeeper       wasmkeeper.Keeper
	customWasmKeeper customwasmkeeper.Keeper
}

func (w WasmDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	for _, msg := range tx.GetMsgs() {
		switch m := msg.(type) {
		case *wasmtypes.MsgStoreCode:
			// Increase the cost of storing code of smart contract
			params := w.customWasmKeeper.GetParams(ctx)
			codeLen := len(m.WASMByteCode)
			baseGas := params.StoreCodeBaseGas
			extraGas := baseGas + uint64(codeLen)*params.StoreCodeMultiplier
			ctx.GasMeter().ConsumeGas(extraGas, "custom wasm store code penalty")
		case *wasmtypes.MsgInstantiateContract:
			// Increase the cost of instantiation of smart contract
			params := w.customWasmKeeper.GetParams(ctx)
			byteCode, err := w.wasmKeeper.GetByteCode(ctx, m.CodeID)
			if err != nil {
				return ctx, errorsmod.Wrapf(err, "failed to load bytecode for code id %d", m.CodeID)
			}
			codeSize := len(byteCode)
			baseGas := params.InstBaseGas
			extraGas := baseGas + uint64(codeSize)*params.InstMultiplier
			ctx.GasMeter().ConsumeGas(extraGas, "custom wasm instantial penalty")
		}
	}
	return next(ctx, tx, simulate)
}
