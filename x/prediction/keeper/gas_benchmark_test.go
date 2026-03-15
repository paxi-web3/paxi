package keeper

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/x/prediction/types"
)

const (
	benchmarkTradeCount       = 500
	benchmarkMinGasPriceUpaxi = 0.05
	benchmarkUpaxiPerPaxi     = 1_000_000.0
)

type mockPRC20ExecKeeper struct {
	gasPerExecute uint64
}

func (m mockPRC20ExecKeeper) Execute(
	ctx sdk.Context,
	_ sdk.AccAddress,
	_ sdk.AccAddress,
	_ []byte,
	_ sdk.Coins,
) ([]byte, error) {
	if m.gasPerExecute > 0 {
		ctx.GasMeter().ConsumeGas(m.gasPerExecute, "mock prc20 execute")
	}
	return nil, nil
}

func BenchmarkApplyTradeBatch500Native(b *testing.B) {
	var totalGas uint64
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		k, ctx, bank := setupKeeper(b)
		creator := testAddress(160)
		resolver := testAddress(161)
		buyer := testAddress(162)
		seller := testAddress(163)

		params := k.GetParams(ctx)
		params.MaxBatchSize = benchmarkTradeCount
		params.MarketFeeBps = 100
		params.ResolverFeeSharePercent = 100
		k.SetParams(ctx, params)

		marketID := mustCreateMarket(b, k, ctx, bank, creator, resolver)
		mustFund(bank, buyer, 10_000_000)
		sellerPos := k.getPositionOrDefault(ctx, marketID, sdk.MustAccAddressFromBech32(seller))
		k.mustSetPositionInts(
			sellerPos,
			sdkmath.NewInt(benchmarkTradeCount),
			sdkmath.ZeroInt(),
			sdkmath.ZeroInt(),
			sdkmath.ZeroInt(),
		)
		k.SetPosition(ctx, sellerPos)

		buyOrderID := mustPlaceOrder(b, k, ctx, &types.MsgPlaceOrder{
			Trader:     buyer,
			MarketId:   marketID,
			Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
			OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
			Amount:     "500",
			LimitPrice: "10000",
		})
		sellOrderID := mustPlaceOrder(b, k, ctx, &types.MsgPlaceOrder{
			Trader:     seller,
			MarketId:   marketID,
			Side:       types.OrderSide_ORDER_SIDE_SELL_YES,
			OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
			Amount:     "500",
			LimitPrice: "10000",
		})

		trades := make([]types.TradeMatch, benchmarkTradeCount)
		for j := 0; j < benchmarkTradeCount; j++ {
			trades[j] = types.TradeMatch{
				TradeId:        "native-bench-" + intToStr(uint64(j)),
				OrderAId:       buyOrderID,
				OrderBId:       sellOrderID,
				MatchAmount:    "1",
				ExecutionPrice: "10000",
			}
		}

		startGas := ctx.GasMeter().GasConsumed()
		b.StartTimer()
		_, err := k.ApplyTradeBatch(ctx, &types.MsgApplyTradeBatch{
			Sender:   resolver,
			MarketId: marketID,
			BatchId:  "native-bench-batch",
			Trades:   trades,
		})
		b.StopTimer()
		if err != nil {
			b.Fatalf("apply batch failed: %v", err)
		}
		totalGas += ctx.GasMeter().GasConsumed() - startGas
	}

	reportGasCostMetrics(b, totalGas, benchmarkTradeCount)
}

func BenchmarkApplyTradeBatch500PRC20(b *testing.B) {
	var totalGas uint64
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		prc20Exec := mockPRC20ExecKeeper{gasPerExecute: 250_000}
		prc20Query := newMockPRC20QueryKeeper()
		k, ctx, bank := setupBenchmarkKeeperWithPRC20(b, prc20Exec, prc20Query)

		creator := testAddress(170)
		resolver := testAddress(171)
		yesBuyer := testAddress(172)
		noBuyer := testAddress(173)
		contract := testAddress(174)

		params := k.GetParams(ctx)
		params.MaxBatchSize = benchmarkTradeCount
		params.MarketFeeBps = 100
		params.ResolverFeeSharePercent = 100
		k.SetParams(ctx, params)

		mustFund(bank, creator, 1_000_000)
		now := ctx.BlockTime().Unix()
		marketID, err := k.CreateMarket(ctx, &types.MsgCreateMarket{
			Creator:                creator,
			Resolver:               resolver,
			Title:                  "PRC20 benchmark market",
			Description:            "benchmark",
			Rule:                   "benchmark",
			OutcomeType:            "BINARY",
			Outcomes:               []string{"YES", "NO"},
			CollateralType:         types.CollateralType_COLLATERAL_TYPE_PRC20,
			CollateralContractAddr: contract,
			OpenTime:               now - 10,
			CloseTime:              now + 3600,
			ResolveTime:            now + 7200,
		})
		if err != nil {
			b.Fatalf("create market failed: %v", err)
		}

		moduleAddr := mockAccountKeeper{}.GetModuleAddress(types.ModuleName).String()
		required := sdkmath.NewInt(benchmarkTradeCount).MulRaw(10_000)
		prc20Query.setAllowance(contract, yesBuyer, moduleAddr, required)
		prc20Query.setAllowance(contract, noBuyer, moduleAddr, required)

		buyYesID := mustPlaceOrder(b, k, ctx, &types.MsgPlaceOrder{
			Trader:     yesBuyer,
			MarketId:   marketID,
			Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
			OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
			Amount:     "500",
			LimitPrice: "10000",
		})
		buyNoID := mustPlaceOrder(b, k, ctx, &types.MsgPlaceOrder{
			Trader:     noBuyer,
			MarketId:   marketID,
			Side:       types.OrderSide_ORDER_SIDE_BUY_NO,
			OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
			Amount:     "500",
			LimitPrice: "10000",
		})

		trades := make([]types.TradeMatch, benchmarkTradeCount)
		for j := 0; j < benchmarkTradeCount; j++ {
			trades[j] = types.TradeMatch{
				TradeId:        "prc20-bench-" + intToStr(uint64(j)),
				OrderAId:       buyYesID,
				OrderBId:       buyNoID,
				MatchAmount:    "1",
				ExecutionPrice: "10000",
			}
		}

		startGas := ctx.GasMeter().GasConsumed()
		b.StartTimer()
		_, err = k.ApplyTradeBatch(ctx, &types.MsgApplyTradeBatch{
			Sender:   resolver,
			MarketId: marketID,
			BatchId:  "prc20-bench-batch",
			Trades:   trades,
		})
		b.StopTimer()
		if err != nil {
			b.Fatalf("apply batch failed: %v", err)
		}
		totalGas += ctx.GasMeter().GasConsumed() - startGas
	}

	reportGasCostMetrics(b, totalGas, benchmarkTradeCount)
}

func setupBenchmarkKeeperWithPRC20(
	b testing.TB,
	prc20Exec PRC20Keeper,
	prc20Query PRC20QueryKeeper,
) (Keeper, sdk.Context, *mockBankKeeper) {
	b.Helper()

	k, ctx, bank, _ := setupKeeperWithPRC20Query(b)
	k = NewKeeper(
		k.cdc,
		bank,
		mockAccountKeeper{},
		prc20Exec,
		prc20Query,
		k.storeKey,
		k.storeService,
		k.authority,
	)
	k.InitGenesis(ctx, types.DefaultGenesisState())
	params := k.GetParams(ctx)
	params.CreateMarketBond = "100"
	params.CreateMarketBondDenom = "upaxi"
	params.MarketFeeBps = 100
	params.ResolverFeeSharePercent = 100
	params.MaxBatchSize = benchmarkTradeCount
	k.SetParams(ctx, params)

	return k, ctx, bank
}

func reportGasCostMetrics(b *testing.B, totalGas uint64, tradeCount uint64) {
	if b.N == 0 {
		return
	}
	avgGasPerBatch := float64(totalGas) / float64(b.N)
	avgGasPerTrade := avgGasPerBatch / float64(tradeCount)
	avgPaxiPerBatch := (avgGasPerBatch * benchmarkMinGasPriceUpaxi) / benchmarkUpaxiPerPaxi
	avgPaxiPerTrade := (avgGasPerTrade * benchmarkMinGasPriceUpaxi) / benchmarkUpaxiPerPaxi

	b.ReportMetric(avgGasPerBatch, "gas/batch")
	b.ReportMetric(avgGasPerTrade, "gas/trade")
	b.ReportMetric(avgPaxiPerBatch, "paxi/batch@0.05")
	b.ReportMetric(avgPaxiPerTrade, "paxi/trade@0.05")
}
