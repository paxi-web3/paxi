package app

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BlockStatusDecorator struct {
	App *PaxiApp
}

const fileName = "block_status.json"

type BlockStatusData struct {
	TotalTxs uint64 `json:"total_txs"`
}

func (p BlockStatusDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	startGas := ctx.GasMeter().GasConsumed() // consumed at the begin

	ctx, err := next(ctx, tx, simulate)
	if err != nil || ctx.IsCheckTx() || simulate {
		return ctx, err
	}

	endGas := ctx.GasMeter().GasConsumed() // consumed at the end

	// accumulate gas
	gasUsed := endGas - startGas
	p.App.CurrentBlockGasUsed += gasUsed

	// accumulate txs
	p.App.TotalTxs += 1

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
	expectedMaxGasUsed := 5_000_000
	lastGasUsed := app.LastBlockGasUsed
	minGasPrice := 0.025 // upaxi
	estimatedGasPrice := minGasPrice + math.Log(1+float64(lastGasUsed)/float64(expectedMaxGasUsed))*minGasPrice
	return float32(estimatedGasPrice)
}

func (app *PaxiApp) GetTotalTxs() uint64 {
	return app.TotalTxs
}

func (app *PaxiApp) ReadBlockStatusFromFile() error {
	cleanPath := filepath.Clean("/" + fileName)
	fullPath := filepath.Join(DefaultNodeHome, cleanPath)

	if !strings.HasPrefix(fullPath, filepath.Clean(DefaultNodeHome)+string(os.PathSeparator)) {
		return os.ErrPermission
	}

	bz, err := os.ReadFile(fullPath)
	if bz == nil || err != nil {
		return nil
	}

	bsData := BlockStatusData{}
	if err = json.Unmarshal(bz, &bsData); err != nil {
		return nil
	}

	app.TotalTxs = bsData.TotalTxs

	return nil
}

func (app *PaxiApp) WriteBlockStatusToFile() error {
	cleanPath := filepath.Clean("/" + fileName)
	fullPath := filepath.Join(DefaultNodeHome, cleanPath)

	if !strings.HasPrefix(fullPath, filepath.Clean(DefaultNodeHome)+string(os.PathSeparator)) {
		return os.ErrPermission
	}

	tmpFile, err := os.CreateTemp(DefaultNodeHome, "tmpfile-*")
	if err != nil {
		return err
	}

	defer os.Remove(tmpFile.Name())

	bsData := BlockStatusData{
		TotalTxs: app.TotalTxs,
	}

	bz, err := json.Marshal(bsData)
	if err != nil {
		return err
	}

	if _, err := tmpFile.Write(bz); err != nil {
		return err
	}

	if err := tmpFile.Close(); err != nil {
		return err
	}

	return os.Rename(tmpFile.Name(), fullPath)
}
