package types

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	wasmvmtypes "github.com/CosmWasm/wasmvm/v2/types"
)

func SmartContractGasRegisterConfig() wasmtypes.WasmGasRegisterConfig {
	return wasmtypes.WasmGasRegisterConfig{
		// Each time a wasm instance is created (cold load from disk)
		InstanceCost: 200_000,

		// Discounted cost when instance is reused from memory cache (e.g. pinned)
		InstanceCostDiscount: 100_000,

		// Cost per byte to compile wasm bytecode on store
		// NOTE: this code will only take effect in the WasmDecorator
		CompileCost: 300, // Increased from default (3).

		// Gas multiplier: how many CosmWasm gas = 1 Cosmos SDK gas
		GasMultiplier: 150_000, // Higher = more expensive per VM operation

		// Cost per attribute emitted by a wasm event (per count)
		EventPerAttributeCost: 300,

		// Cost per byte of key+value in event attributes
		EventAttributeDataCost: 150,

		// Free bytes allowed per event (before charging gas)
		EventAttributeDataFreeTier: 0, // No free tier

		// Cost per byte of the message passed to the contract
		ContractMessageDataCost: 150, // Penalize large messages

		// Cost per custom event (not per attribute)
		CustomEventCost: 300,

		// Cost to decompress `.wasm.gz` per byte during store
		UncompressCost: wasmvmtypes.UFraction{Numerator: 100, Denominator: 100}, // 1.0 gas per byte
	}
}
