package app

import (
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PaxiDecorator struct {
	App *PaxiApp
}

func (p PaxiDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	startGas := ctx.GasMeter().GasConsumed() // consumed at the begin

	ctx, err := next(ctx, tx, simulate)
	if err != nil {
		return ctx, err
	}

	endGas := ctx.GasMeter().GasConsumed() // consumed at the end

	// accumulate gas
	gasUsed := endGas - startGas
	p.App.CurrentBlockGasUsed += gasUsed

	return ctx, nil
}

func (app *PaxiApp) GetLastBlockGasUsed() uint64 {
	return app.LastBlockGasUsed
}

func (app *PaxiApp) SetLastBlockGasUsed() {
	app.LastBlockGasUsed = app.CurrentBlockGasUsed
	app.CurrentBlockGasUsed = 0
}

func (app *PaxiApp) GetEstimatedGasPrice() float32 {
	// This is an estimation for reasonable gas price by gas used in last block
	expectedMaxGasUsed := 15_000_000
	lastGasUsed := app.LastBlockGasUsed
	minGasPrice := 0.1 // upaxi
	estimatedGasPrice := minGasPrice + math.Log(1+float64(lastGasUsed)/float64(expectedMaxGasUsed))*minGasPrice
	return float32(estimatedGasPrice)
}
