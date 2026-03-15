package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/paxi-web3/paxi/x/prediction/types"
	"github.com/stretchr/testify/require"
)

type mockBankKeeper struct {
	accountBalances map[string]sdkmath.Int
	moduleBalances  map[string]sdkmath.Int
	transferCount   int
}

func newMockBankKeeper() *mockBankKeeper {
	return &mockBankKeeper{
		accountBalances: make(map[string]sdkmath.Int),
		moduleBalances:  make(map[string]sdkmath.Int),
	}
}

func balanceKey(addr, denom string) string { return addr + "|" + denom }

func (m *mockBankKeeper) addAccountBalance(addr, denom string, amount sdkmath.Int) {
	key := balanceKey(addr, denom)
	cur, ok := m.accountBalances[key]
	if !ok {
		cur = sdkmath.ZeroInt()
	}
	m.accountBalances[key] = cur.Add(amount)
}

func (m *mockBankKeeper) addModuleBalance(denom string, amount sdkmath.Int) {
	cur, ok := m.moduleBalances[denom]
	if !ok {
		cur = sdkmath.ZeroInt()
	}
	m.moduleBalances[denom] = cur.Add(amount)
}

func (m *mockBankKeeper) AccountBalance(addr, denom string) sdkmath.Int {
	cur, ok := m.accountBalances[balanceKey(addr, denom)]
	if !ok {
		return sdkmath.ZeroInt()
	}
	return cur
}

func (m *mockBankKeeper) ModuleBalance(denom string) sdkmath.Int {
	cur, ok := m.moduleBalances[denom]
	if !ok {
		return sdkmath.ZeroInt()
	}
	return cur
}

func (m *mockBankKeeper) moveFromAccount(addr, denom string, amount sdkmath.Int) error {
	key := balanceKey(addr, denom)
	cur, ok := m.accountBalances[key]
	if !ok {
		cur = sdkmath.ZeroInt()
	}
	if cur.LT(amount) {
		return fmt.Errorf("insufficient funds")
	}
	m.accountBalances[key] = cur.Sub(amount)
	return nil
}

func (m *mockBankKeeper) moveFromModule(denom string, amount sdkmath.Int) error {
	cur, ok := m.moduleBalances[denom]
	if !ok {
		cur = sdkmath.ZeroInt()
	}
	if cur.LT(amount) {
		return fmt.Errorf("insufficient module funds")
	}
	m.moduleBalances[denom] = cur.Sub(amount)
	return nil
}

func (m *mockBankKeeper) SendCoinsFromAccountToModule(_ context.Context, sender sdk.AccAddress, _ string, amt sdk.Coins) error {
	for i := range amt {
		coin := amt[i]
		if err := m.moveFromAccount(sender.String(), coin.Denom, coin.Amount); err != nil {
			return err
		}
		m.addModuleBalance(coin.Denom, coin.Amount)
	}
	m.transferCount++
	return nil
}

func (m *mockBankKeeper) SendCoinsFromModuleToAccount(_ context.Context, _ string, recipient sdk.AccAddress, amt sdk.Coins) error {
	for i := range amt {
		coin := amt[i]
		if err := m.moveFromModule(coin.Denom, coin.Amount); err != nil {
			return err
		}
		m.addAccountBalance(recipient.String(), coin.Denom, coin.Amount)
	}
	m.transferCount++
	return nil
}

func (m *mockBankKeeper) SendCoins(_ context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	moduleAddr := mockAccountKeeper{}.GetModuleAddress(types.ModuleName).String()
	fromIsModule := fromAddr.String() == moduleAddr
	toIsModule := toAddr.String() == moduleAddr

	for i := range amt {
		coin := amt[i]
		switch {
		case fromIsModule:
			if err := m.moveFromModule(coin.Denom, coin.Amount); err != nil {
				return err
			}
		default:
			if err := m.moveFromAccount(fromAddr.String(), coin.Denom, coin.Amount); err != nil {
				return err
			}
		}

		switch {
		case toIsModule:
			m.addModuleBalance(coin.Denom, coin.Amount)
		default:
			m.addAccountBalance(toAddr.String(), coin.Denom, coin.Amount)
		}
	}
	m.transferCount++
	return nil
}

type mockAccountKeeper struct{}

func (m mockAccountKeeper) GetModuleAddress(moduleName string) sdk.AccAddress {
	addr := make([]byte, 20)
	copy(addr, []byte(moduleName))
	return sdk.AccAddress(addr)
}

type mockPRC20QueryKeeper struct {
	allowances map[string]sdkmath.Int
}

func newMockPRC20QueryKeeper() *mockPRC20QueryKeeper {
	return &mockPRC20QueryKeeper{allowances: make(map[string]sdkmath.Int)}
}

func allowanceKey(contract, owner, spender string) string {
	return contract + "|" + owner + "|" + spender
}

func (m *mockPRC20QueryKeeper) setAllowance(contract, owner, spender string, allowance sdkmath.Int) {
	m.allowances[allowanceKey(contract, owner, spender)] = allowance
}

func (m *mockPRC20QueryKeeper) QuerySmart(_ context.Context, contractAddress sdk.AccAddress, req []byte) ([]byte, error) {
	var query prc20AllowanceQuery
	if err := json.Unmarshal(req, &query); err != nil {
		return nil, err
	}

	allowance, ok := m.allowances[allowanceKey(contractAddress.String(), query.Allowance.Owner, query.Allowance.Spender)]
	if !ok {
		allowance = sdkmath.ZeroInt()
	}
	return json.Marshal(prc20AllowanceResponse{Allowance: allowance.String()})
}

func setupKeeper(tb testing.TB) (Keeper, sdk.Context, *mockBankKeeper) {
	tb.Helper()

	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	testCtx := sdktestutil.DefaultContextWithDB(tb, storeKey, storetypes.NewTransientStoreKey("transient_test"))
	ctx := testCtx.Ctx

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	bank := newMockBankKeeper()
	k := NewKeeper(
		cdc,
		bank,
		mockAccountKeeper{},
		nil,
		nil,
		storeKey,
		runtime.NewKVStoreService(storeKey),
		"",
	)
	k.InitGenesis(ctx, types.DefaultGenesisState())

	params := k.GetParams(ctx)
	params.CreateMarketBond = "100"
	params.CreateMarketBondDenom = "upaxi"
	params.MarketFeeBps = 100
	params.ResolverFeeSharePercent = 100
	params.MaxBatchSize = 100
	k.SetParams(ctx, params)

	return k, ctx, bank
}

func setupKeeperWithAuthority(tb testing.TB, authority string) (Keeper, sdk.Context, *mockBankKeeper) {
	tb.Helper()

	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	testCtx := sdktestutil.DefaultContextWithDB(tb, storeKey, storetypes.NewTransientStoreKey("transient_test"))
	ctx := testCtx.Ctx

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	bank := newMockBankKeeper()
	k := NewKeeper(
		cdc,
		bank,
		mockAccountKeeper{},
		nil,
		nil,
		storeKey,
		runtime.NewKVStoreService(storeKey),
		authority,
	)
	k.InitGenesis(ctx, types.DefaultGenesisState())

	return k, ctx, bank
}

func setupKeeperWithPRC20Query(tb testing.TB) (Keeper, sdk.Context, *mockBankKeeper, *mockPRC20QueryKeeper) {
	tb.Helper()

	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	testCtx := sdktestutil.DefaultContextWithDB(tb, storeKey, storetypes.NewTransientStoreKey("transient_test"))
	ctx := testCtx.Ctx

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	bank := newMockBankKeeper()
	prc20Query := newMockPRC20QueryKeeper()
	k := NewKeeper(
		cdc,
		bank,
		mockAccountKeeper{},
		nil,
		prc20Query,
		storeKey,
		runtime.NewKVStoreService(storeKey),
		"",
	)
	k.InitGenesis(ctx, types.DefaultGenesisState())

	params := k.GetParams(ctx)
	params.CreateMarketBond = "100"
	params.CreateMarketBondDenom = "upaxi"
	params.MarketFeeBps = 100
	params.ResolverFeeSharePercent = 100
	params.MaxBatchSize = 100
	k.SetParams(ctx, params)

	return k, ctx, bank, prc20Query
}

func testAddress(seed byte) string {
	addr := make([]byte, 20)
	for i := range addr {
		addr[i] = seed
	}
	return sdk.AccAddress(addr).String()
}

func mustFund(bank *mockBankKeeper, addr string, amount int64) {
	bank.addAccountBalance(addr, "upaxi", sdkmath.NewInt(amount))
}

func mustCreateMarket(tb testing.TB, k Keeper, ctx sdk.Context, bank *mockBankKeeper, creator, resolver string) uint64 {
	tb.Helper()
	mustFund(bank, creator, 1_000_000)
	now := ctx.BlockTime().Unix()
	msg := &types.MsgCreateMarket{
		Creator:         creator,
		Resolver:        resolver,
		Title:           "BTC > 100k?",
		Description:     "test market",
		Rule:            "simple",
		OutcomeType:     "BINARY",
		Outcomes:        []string{"YES", "NO"},
		CollateralType:  types.CollateralType_COLLATERAL_TYPE_NATIVE,
		CollateralDenom: "upaxi",
		OpenTime:        now - 10,
		CloseTime:       now + 3600,
		ResolveTime:     now + 7200,
	}
	marketID, err := k.CreateMarket(ctx, msg)
	require.NoError(tb, err)
	return marketID
}

func mustPlaceOrder(tb testing.TB, k Keeper, ctx sdk.Context, msg *types.MsgPlaceOrder) uint64 {
	tb.Helper()
	if msg.ExpireBh == 0 {
		msg.ExpireBh = ctx.BlockHeight() + 100
	}
	orderID, err := k.PlaceOrder(ctx, msg)
	require.NoError(tb, err)
	return orderID
}

func TestPlaceOrderWithoutEscrow(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(1)
	resolver := testAddress(2)
	trader := testAddress(3)

	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)
	beforeTransfers := bank.transferCount

	orderID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     trader,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "10",
		LimitPrice: "10000",
	})

	order, found := k.GetOrder(ctx, marketID, orderID)
	require.True(t, found)
	require.Equal(t, "10", order.Amount)
	require.Equal(t, "0", order.FilledAmount)
	require.Equal(t, "10", order.RemainingAmount)
	require.Equal(t, types.OrderStatus_ORDER_STATUS_OPEN, order.Status)

	require.Equal(t, beforeTransfers, bank.transferCount, "place order must not transfer funds")

	byID, found := k.GetOrderByID(ctx, orderID)
	require.True(t, found)
	require.Equal(t, marketID, byID.MarketId)
	require.Equal(t, trader, byID.Trader)
}

func TestPlaceOrderValidation(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(11)
	resolver := testAddress(12)
	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)
	trader := testAddress(13)

	tests := []struct {
		name    string
		msg     types.MsgPlaceOrder
		wantErr string
	}{
		{
			name: "limit requires limit_price",
			msg: types.MsgPlaceOrder{
				Trader:    trader,
				MarketId:  marketID,
				Side:      types.OrderSide_ORDER_SIDE_BUY_YES,
				OrderType: types.OrderType_ORDER_TYPE_LIMIT,
				Amount:    "10",
			},
			wantErr: "invalid limit_price",
		},
		{
			name: "market requires worst_price",
			msg: types.MsgPlaceOrder{
				Trader:    trader,
				MarketId:  marketID,
				Side:      types.OrderSide_ORDER_SIDE_BUY_NO,
				OrderType: types.OrderType_ORDER_TYPE_MARKET,
				Amount:    "10",
			},
			wantErr: "invalid worst_price",
		},
		{
			name: "limit order rejects worst_price",
			msg: types.MsgPlaceOrder{
				Trader:     trader,
				MarketId:   marketID,
				Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
				OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
				Amount:     "10",
				LimitPrice: "10000",
				WorstPrice: "20000",
			},
			wantErr: "worst_price must be empty for limit order",
		},
		{
			name: "market order rejects limit_price",
			msg: types.MsgPlaceOrder{
				Trader:     trader,
				MarketId:   marketID,
				Side:       types.OrderSide_ORDER_SIDE_BUY_NO,
				OrderType:  types.OrderType_ORDER_TYPE_MARKET,
				Amount:     "10",
				LimitPrice: "10000",
				WorstPrice: "20000",
			},
			wantErr: "limit_price must be empty for market order",
		},
		{
			name: "valid limit order",
			msg: types.MsgPlaceOrder{
				Trader:     trader,
				MarketId:   marketID,
				Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
				OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
				Amount:     "10",
				LimitPrice: "10000",
			},
		},
		{
			name: "valid market order",
			msg: types.MsgPlaceOrder{
				Trader:     trader,
				MarketId:   marketID,
				Side:       types.OrderSide_ORDER_SIDE_BUY_NO,
				OrderType:  types.OrderType_ORDER_TYPE_MARKET,
				Amount:     "10",
				WorstPrice: "20000",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			msg := tc.msg
			if msg.ExpireBh == 0 {
				msg.ExpireBh = ctx.BlockHeight() + 100
			}
			_, err := k.PlaceOrder(ctx, &msg)
			if tc.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.wantErr)
			}
		})
	}
}

func TestPlaceOrderBuyPRC20AllowanceCheck(t *testing.T) {
	k, ctx, bank, prc20Query := setupKeeperWithPRC20Query(t)
	creator := testAddress(14)
	resolver := testAddress(15)
	buyer := testAddress(16)
	contract := testAddress(17)

	mustFund(bank, creator, 1_000_000)
	now := ctx.BlockTime().Unix()
	marketID, err := k.CreateMarket(ctx, &types.MsgCreateMarket{
		Creator:                creator,
		Resolver:               resolver,
		Title:                  "PRC20 market",
		Description:            "test market",
		Rule:                   "simple",
		OutcomeType:            "BINARY",
		Outcomes:               []string{"YES", "NO"},
		CollateralType:         types.CollateralType_COLLATERAL_TYPE_PRC20,
		CollateralContractAddr: contract,
		OpenTime:               now - 10,
		CloseTime:              now + 3600,
		ResolveTime:            now + 7200,
	})
	require.NoError(t, err)

	moduleAddr := mockAccountKeeper{}.GetModuleAddress(types.ModuleName)
	prc20Query.setAllowance(contract, buyer, moduleAddr.String(), sdkmath.NewInt(99_999))

	_, err = k.PlaceOrder(ctx, &types.MsgPlaceOrder{
		Trader:     buyer,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "10",
		LimitPrice: "10000",
		ExpireBh:   ctx.BlockHeight() + 100,
	})
	require.ErrorContains(t, err, "insufficient prc20 allowance")

	prc20Query.setAllowance(contract, buyer, moduleAddr.String(), sdkmath.NewInt(100_000))

	_, err = k.PlaceOrder(ctx, &types.MsgPlaceOrder{
		Trader:     buyer,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "10",
		LimitPrice: "10000",
		ExpireBh:   ctx.BlockHeight() + 100,
	})
	require.NoError(t, err)
}

func TestPlaceOrderExpireBhBounds(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	ctx = ctx.WithBlockHeight(10)
	creator := testAddress(91)
	resolver := testAddress(92)
	trader := testAddress(93)
	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)

	_, err := k.PlaceOrder(ctx, &types.MsgPlaceOrder{
		Trader:     trader,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "1",
		LimitPrice: "10000",
		ExpireBh:   ctx.BlockHeight(),
	})
	require.ErrorContains(t, err, "expire_bh must be greater than current block height")

	params := k.GetParams(ctx)
	_, err = k.PlaceOrder(ctx, &types.MsgPlaceOrder{
		Trader:     trader,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "1",
		LimitPrice: "10000",
		ExpireBh:   ctx.BlockHeight() + int64(params.MaxOrderLifetimeBh) + 1,
	})
	require.ErrorContains(t, err, "expire_bh exceeds max_order_lifetime_bh")
}

func TestPlaceOrderOpenOrderLimitAndCancelRelease(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(94)
	resolver := testAddress(95)
	trader := testAddress(96)
	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)

	params := k.GetParams(ctx)
	params.MaxOpenOrdersPerUser = 2
	params.MaxOpenOrdersPerMarket = 2
	k.SetParams(ctx, params)

	order1 := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     trader,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "1",
		LimitPrice: "10000",
		ExpireBh:   ctx.BlockHeight() + 100,
	})
	mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     trader,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "1",
		LimitPrice: "10000",
		ExpireBh:   ctx.BlockHeight() + 101,
	})

	_, err := k.PlaceOrder(ctx, &types.MsgPlaceOrder{
		Trader:     trader,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "1",
		LimitPrice: "10000",
		ExpireBh:   ctx.BlockHeight() + 102,
	})
	require.ErrorContains(t, err, "max_open_orders_per_user exceeded")

	err = k.CancelOrder(ctx, &types.MsgCancelOrder{
		Trader:   trader,
		MarketId: marketID,
		OrderId:  order1,
	})
	require.NoError(t, err)

	_, err = k.PlaceOrder(ctx, &types.MsgPlaceOrder{
		Trader:     trader,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "1",
		LimitPrice: "10000",
		ExpireBh:   ctx.BlockHeight() + 103,
	})
	require.NoError(t, err)
}

func TestAutoPruneOrders(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	ctx = ctx.WithBlockHeight(100)

	creator := testAddress(101)
	resolver := testAddress(102)
	trader := testAddress(103)
	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)

	params := k.GetParams(ctx)
	params.OrderPruneIntervalBh = 1
	params.OrderPruneRetainBh = 2
	params.OrderPruneScanLimit = 100
	params.OrderPruneDeleteLimit = 100
	k.SetParams(ctx, params)

	expiredOrderID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     trader,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "1",
		LimitPrice: "10000",
		ExpireBh:   ctx.BlockHeight() + 1,
	})
	ctxExpire := ctx.WithBlockHeight(101)
	expiredOrder, found := k.GetOrder(ctxExpire, marketID, expiredOrderID)
	require.True(t, found)
	require.NoError(t, k.expireOrderIfNeeded(ctxExpire, expiredOrder))
	expiredOrder, found = k.GetOrder(ctxExpire, marketID, expiredOrderID)
	require.True(t, found)
	require.Equal(t, int64(101), expiredOrder.ClosedBh)

	openOrderID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     trader,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "1",
		LimitPrice: "10000",
		ExpireBh:   ctx.BlockHeight() + 100,
	})

	ctx = ctx.WithBlockHeight(103) // threshold = 101, expired order should be pruned
	k.AutoPruneOrders(ctx)

	_, found = k.GetOrder(ctx, marketID, expiredOrderID)
	require.False(t, found)
	_, found = k.GetOrderByID(ctx, expiredOrderID)
	require.False(t, found)
	_, found = k.GetOrder(ctx, marketID, openOrderID)
	require.True(t, found)
	_, found = k.GetOrderByID(ctx, openOrderID)
	require.True(t, found)
}

func TestAutoPruneOrdersWithCursorAndDeleteLimit(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	ctx = ctx.WithBlockHeight(200)

	creator := testAddress(111)
	resolver := testAddress(112)
	trader := testAddress(113)
	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)

	params := k.GetParams(ctx)
	params.OrderPruneIntervalBh = 1
	params.OrderPruneRetainBh = 1
	params.OrderPruneScanLimit = 100
	params.OrderPruneDeleteLimit = 1
	k.SetParams(ctx, params)

	orderIDs := make([]uint64, 0, 3)
	for i := 0; i < 3; i++ {
		orderID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
			Trader:     trader,
			MarketId:   marketID,
			Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
			OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
			Amount:     "1",
			LimitPrice: "10000",
			ExpireBh:   ctx.BlockHeight() + 1,
		})
		orderIDs = append(orderIDs, orderID)
	}
	ctxExpire := ctx.WithBlockHeight(201)
	for _, id := range orderIDs {
		order, found := k.GetOrder(ctxExpire, marketID, id)
		require.True(t, found)
		require.NoError(t, k.expireOrderIfNeeded(ctxExpire, order))
	}

	ctx = ctx.WithBlockHeight(202)
	k.AutoPruneOrders(ctx)
	_, found := k.GetOrder(ctx, marketID, orderIDs[0])
	require.False(t, found)
	_, found = k.GetOrder(ctx, marketID, orderIDs[1])
	require.True(t, found)

	ctx = ctx.WithBlockHeight(203)
	k.AutoPruneOrders(ctx)
	_, found = k.GetOrder(ctx, marketID, orderIDs[1])
	require.False(t, found)
	_, found = k.GetOrder(ctx, marketID, orderIDs[2])
	require.True(t, found)

	ctx = ctx.WithBlockHeight(204)
	k.AutoPruneOrders(ctx)
	_, found = k.GetOrder(ctx, marketID, orderIDs[2])
	require.False(t, found)
}

func TestApplyTradeBatchBuyYesBuyNoFeeNotExtra(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(18)
	resolver := testAddress(19)
	yesBuyer := testAddress(20)
	noBuyer := testAddress(21)

	params := k.GetParams(ctx)
	params.MarketFeeBps = 100
	params.ResolverFeeSharePercent = 100
	k.SetParams(ctx, params)

	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)
	mustFund(bank, yesBuyer, 20_000)
	mustFund(bank, noBuyer, 20_000)

	buyYesID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     yesBuyer,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "1",
		LimitPrice: "10000",
	})
	buyNoID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     noBuyer,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_NO,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "1",
		LimitPrice: "10000",
	})

	_, err := k.ApplyTradeBatch(ctx, &types.MsgApplyTradeBatch{
		Sender:   resolver,
		MarketId: marketID,
		BatchId:  "batch-buybuy-fee",
		Trades: []types.TradeMatch{{
			TradeId:        "t-buybuy-fee-1",
			OrderAId:       buyYesID,
			OrderBId:       buyNoID,
			MatchAmount:    "1",
			ExecutionPrice: "10000",
		}},
	})
	require.NoError(t, err)

	// Fee should be deducted from collected collateral in module, not extra charged to buyers.
	require.Equal(t, sdkmath.NewInt(10_000), bank.AccountBalance(yesBuyer, "upaxi"))
	require.Equal(t, sdkmath.NewInt(10_000), bank.AccountBalance(noBuyer, "upaxi"))
	require.Equal(t, sdkmath.NewInt(200), bank.AccountBalance(resolver, "upaxi"))
}

func TestSplitMergeNoFee(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(26)
	resolver := testAddress(27)
	trader := testAddress(28)

	params := k.GetParams(ctx)
	params.MarketFeeBps = 100
	params.ResolverFeeSharePercent = 100
	k.SetParams(ctx, params)

	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)
	mustFund(bank, trader, 20_000)

	moduleBefore := bank.ModuleBalance("upaxi")
	traderBefore := bank.AccountBalance(trader, "upaxi")

	err := k.SplitPosition(ctx, &types.MsgSplitPosition{
		Trader:   trader,
		MarketId: marketID,
		Amount:   "1000",
	})
	require.NoError(t, err)

	pos, found := k.GetPosition(ctx, marketID, sdk.MustAccAddressFromBech32(trader))
	require.True(t, found)
	require.Equal(t, "1000", pos.YesShares)
	require.Equal(t, "1000", pos.NoShares)
	require.Equal(t, traderBefore.SubRaw(1000), bank.AccountBalance(trader, "upaxi"))
	require.Equal(t, moduleBefore.AddRaw(1000), bank.ModuleBalance("upaxi"))

	market, found := k.GetMarket(ctx, marketID)
	require.True(t, found)
	require.Equal(t, "1000", market.TotalYesShares)
	require.Equal(t, "1000", market.TotalNoShares)

	resolverBefore := bank.AccountBalance(resolver, "upaxi")
	err = k.MergePosition(ctx, &types.MsgMergePosition{
		Trader:   trader,
		MarketId: marketID,
		Amount:   "400",
	})
	require.NoError(t, err)

	pos, _ = k.GetPosition(ctx, marketID, sdk.MustAccAddressFromBech32(trader))
	require.Equal(t, "600", pos.YesShares)
	require.Equal(t, "600", pos.NoShares)
	require.Equal(t, traderBefore.SubRaw(600), bank.AccountBalance(trader, "upaxi"))
	require.Equal(t, moduleBefore.AddRaw(600), bank.ModuleBalance("upaxi"))
	require.Equal(t, resolverBefore, bank.AccountBalance(resolver, "upaxi"), "split/merge must not charge fee")

	market, _ = k.GetMarket(ctx, marketID)
	require.Equal(t, "600", market.TotalYesShares)
	require.Equal(t, "600", market.TotalNoShares)
}

func TestMergeInsufficientSharesRejected(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(29)
	resolver := testAddress(30)
	trader := testAddress(31)

	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)
	mustFund(bank, trader, 20_000)

	err := k.SplitPosition(ctx, &types.MsgSplitPosition{
		Trader:   trader,
		MarketId: marketID,
		Amount:   "100",
	})
	require.NoError(t, err)

	err = k.MergePosition(ctx, &types.MsgMergePosition{
		Trader:   trader,
		MarketId: marketID,
		Amount:   "101",
	})
	require.ErrorContains(t, err, "insufficient YES shares")
}

func TestApplyTradeBatchPartialFill(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(21)
	resolver := testAddress(22)
	buyer := testAddress(24)
	seller := testAddress(25)

	params := k.GetParams(ctx)
	params.MarketFeeBps = 100
	params.ResolverFeeSharePercent = 100
	k.SetParams(ctx, params)

	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)
	mustFund(bank, buyer, 50_000)
	mustFund(bank, seller, 10)

	sellerPos := k.getPositionOrDefault(ctx, marketID, sdk.MustAccAddressFromBech32(seller))
	k.mustSetPositionInts(sellerPos, sdkmath.NewInt(20), sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt())
	k.SetPosition(ctx, sellerPos)

	buyOrderID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{Trader: buyer, MarketId: marketID, Side: types.OrderSide_ORDER_SIDE_BUY_YES, OrderType: types.OrderType_ORDER_TYPE_LIMIT, Amount: "20", LimitPrice: "20000"})
	sellOrderID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{Trader: seller, MarketId: marketID, Side: types.OrderSide_ORDER_SIDE_SELL_YES, OrderType: types.OrderType_ORDER_TYPE_LIMIT, Amount: "20", LimitPrice: "10000"})

	_, err := k.ApplyTradeBatch(ctx, &types.MsgApplyTradeBatch{
		Sender:   resolver,
		MarketId: marketID,
		BatchId:  "batch-1",
		Trades: []types.TradeMatch{{
			TradeId:        "t-1",
			OrderAId:       buyOrderID,
			OrderBId:       sellOrderID,
			MatchAmount:    "10",
			ExecutionPrice: "10000",
		}},
	})
	require.NoError(t, err)

	buyOrder, _ := k.GetOrder(ctx, marketID, buyOrderID)
	sellOrder, _ := k.GetOrder(ctx, marketID, sellOrderID)
	require.Equal(t, types.OrderStatus_ORDER_STATUS_PARTIALLY_FILLED, buyOrder.Status)
	require.Equal(t, types.OrderStatus_ORDER_STATUS_PARTIALLY_FILLED, sellOrder.Status)
	require.Equal(t, "10", buyOrder.FilledAmount)
	require.Equal(t, "10", buyOrder.RemainingAmount)

	buyerPos, _ := k.GetPosition(ctx, marketID, sdk.MustAccAddressFromBech32(buyer))
	sellerPos, _ = k.GetPosition(ctx, marketID, sdk.MustAccAddressFromBech32(seller))
	require.Equal(t, "10", buyerPos.YesShares)
	require.Equal(t, "10", sellerPos.YesShares)

	require.Equal(t, sdkmath.NewInt(40_000), bank.AccountBalance(buyer, "upaxi"))
	require.Equal(t, sdkmath.NewInt(9_910), bank.AccountBalance(seller, "upaxi"))
	require.Equal(t, sdkmath.NewInt(100), bank.AccountBalance(resolver, "upaxi"))

	market, found := k.GetMarket(ctx, marketID)
	require.True(t, found)
	require.Equal(t, "10", market.TotalTradeVolume)
}

func TestCancelOrder(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(31)
	resolver := testAddress(32)
	trader := testAddress(33)
	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)

	orderID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     trader,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_NO,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "7",
		LimitPrice: "10000",
	})

	err := k.CancelOrder(ctx, &types.MsgCancelOrder{Trader: trader, MarketId: marketID, OrderId: orderID})
	require.NoError(t, err)

	_, found := k.GetOrder(ctx, marketID, orderID)
	require.False(t, found, "unfilled cancelled order should be deleted")
}

func TestCancelOrderPartialKeepsOrder(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(201)
	resolver := testAddress(202)
	buyer := testAddress(203)
	seller := testAddress(204)

	params := k.GetParams(ctx)
	params.MarketFeeBps = 100
	params.ResolverFeeSharePercent = 100
	k.SetParams(ctx, params)

	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)
	mustFund(bank, buyer, 50_000)

	sellerPos := k.getPositionOrDefault(ctx, marketID, sdk.MustAccAddressFromBech32(seller))
	k.mustSetPositionInts(sellerPos, sdkmath.NewInt(20), sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt())
	k.SetPosition(ctx, sellerPos)

	buyOrderID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     buyer,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "20",
		LimitPrice: "10000",
	})
	sellOrderID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     seller,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_SELL_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "20",
		LimitPrice: "10000",
	})

	_, err := k.ApplyTradeBatch(ctx, &types.MsgApplyTradeBatch{
		Sender:   resolver,
		MarketId: marketID,
		BatchId:  "partial-cancel-batch",
		Trades: []types.TradeMatch{{
			TradeId:        "partial-cancel-trade",
			OrderAId:       buyOrderID,
			OrderBId:       sellOrderID,
			MatchAmount:    "10",
			ExecutionPrice: "10000",
		}},
	})
	require.NoError(t, err)

	err = k.CancelOrder(ctx, &types.MsgCancelOrder{Trader: buyer, MarketId: marketID, OrderId: buyOrderID})
	require.NoError(t, err)

	order, found := k.GetOrder(ctx, marketID, buyOrderID)
	require.True(t, found)
	require.Equal(t, types.OrderStatus_ORDER_STATUS_CANCELLED, order.Status)
	require.Equal(t, "10", order.FilledAmount)
}

func TestMarketBestBidAskUpdatedOnPlaceAndCancel(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(211)
	resolver := testAddress(212)
	trader := testAddress(213)
	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)

	buyLowID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     trader,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "1",
		LimitPrice: "20000",
	})
	buyHighID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     trader,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_NO,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "1",
		LimitPrice: "30000",
	})
	askHighID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     trader,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_SELL_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "1",
		LimitPrice: "50000",
	})
	askLowID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     trader,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_SELL_NO,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "1",
		LimitPrice: "40000",
	})
	mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     trader,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
		OrderType:  types.OrderType_ORDER_TYPE_MARKET,
		Amount:     "1",
		WorstPrice: "60000",
	})

	market, found := k.GetMarket(ctx, marketID)
	require.True(t, found)
	require.Equal(t, "30000", market.BestBidPrice)
	require.Equal(t, "40000", market.BestAskPrice)

	require.NoError(t, k.CancelOrder(ctx, &types.MsgCancelOrder{Trader: trader, MarketId: marketID, OrderId: buyHighID}))
	market, _ = k.GetMarket(ctx, marketID)
	require.Equal(t, "20000", market.BestBidPrice)
	require.Equal(t, "40000", market.BestAskPrice)

	require.NoError(t, k.CancelOrder(ctx, &types.MsgCancelOrder{Trader: trader, MarketId: marketID, OrderId: askLowID}))
	market, _ = k.GetMarket(ctx, marketID)
	require.Equal(t, "20000", market.BestBidPrice)
	require.Equal(t, "50000", market.BestAskPrice)

	require.NoError(t, k.CancelOrder(ctx, &types.MsgCancelOrder{Trader: trader, MarketId: marketID, OrderId: buyLowID}))
	require.NoError(t, k.CancelOrder(ctx, &types.MsgCancelOrder{Trader: trader, MarketId: marketID, OrderId: askHighID}))
	market, _ = k.GetMarket(ctx, marketID)
	require.Empty(t, market.BestBidPrice)
	require.Empty(t, market.BestAskPrice)
}

func TestDuplicateTradeIDRejected(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(41)
	resolver := testAddress(42)
	buyer := testAddress(44)
	seller := testAddress(45)

	params := k.GetParams(ctx)
	k.SetParams(ctx, params)

	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)
	mustFund(bank, buyer, 20_000)

	sellerPos := k.getPositionOrDefault(ctx, marketID, sdk.MustAccAddressFromBech32(seller))
	k.mustSetPositionInts(sellerPos, sdkmath.NewInt(10), sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt())
	k.SetPosition(ctx, sellerPos)

	buyOrderID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{Trader: buyer, MarketId: marketID, Side: types.OrderSide_ORDER_SIDE_BUY_YES, OrderType: types.OrderType_ORDER_TYPE_LIMIT, Amount: "2", LimitPrice: "20000"})
	sellOrderID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{Trader: seller, MarketId: marketID, Side: types.OrderSide_ORDER_SIDE_SELL_YES, OrderType: types.OrderType_ORDER_TYPE_LIMIT, Amount: "2", LimitPrice: "10000"})

	_, err := k.ApplyTradeBatch(ctx, &types.MsgApplyTradeBatch{
		Sender:   resolver,
		MarketId: marketID,
		BatchId:  "batch-1",
		Trades: []types.TradeMatch{{
			TradeId:        "dup-trade",
			OrderAId:       buyOrderID,
			OrderBId:       sellOrderID,
			MatchAmount:    "1",
			ExecutionPrice: "10000",
		}},
	})
	require.NoError(t, err)

	_, err = k.ApplyTradeBatch(ctx, &types.MsgApplyTradeBatch{
		Sender:   resolver,
		MarketId: marketID,
		BatchId:  "batch-2",
		Trades: []types.TradeMatch{{
			TradeId:        "dup-trade",
			OrderAId:       buyOrderID,
			OrderBId:       sellOrderID,
			MatchAmount:    "1",
			ExecutionPrice: "10000",
		}},
	})
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrDuplicateTrade)
}

func TestApplyTradeBatchUpdatesLastTradePrice(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(214)
	resolver := testAddress(215)
	buyer := testAddress(216)
	seller := testAddress(217)

	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)
	mustFund(bank, buyer, 20_000)

	sellerPos := k.getPositionOrDefault(ctx, marketID, sdk.MustAccAddressFromBech32(seller))
	k.mustSetPositionInts(sellerPos, sdkmath.NewInt(1), sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt())
	k.SetPosition(ctx, sellerPos)

	buyOrderID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     buyer,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "1",
		LimitPrice: "30000",
	})
	sellOrderID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     seller,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_SELL_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "1",
		LimitPrice: "10000",
	})

	_, err := k.ApplyTradeBatch(ctx, &types.MsgApplyTradeBatch{
		Sender:   resolver,
		MarketId: marketID,
		BatchId:  "batch-last-price",
		Trades: []types.TradeMatch{{
			TradeId:        "t-last-price",
			OrderAId:       buyOrderID,
			OrderBId:       sellOrderID,
			MatchAmount:    "1",
			ExecutionPrice: "20000",
		}},
	})
	require.NoError(t, err)

	market, found := k.GetMarket(ctx, marketID)
	require.True(t, found)
	require.Equal(t, "20000", market.LastTradePrice)
	require.Equal(t, "1", market.TotalTradeVolume)
	require.Empty(t, market.BestBidPrice)
	require.Empty(t, market.BestAskPrice)
}

func TestSettlementInsufficientWalletBalance(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(51)
	resolver := testAddress(52)
	buyer := testAddress(54)
	seller := testAddress(55)

	params := k.GetParams(ctx)
	k.SetParams(ctx, params)

	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)
	mustFund(bank, buyer, 500)

	sellerPos := k.getPositionOrDefault(ctx, marketID, sdk.MustAccAddressFromBech32(seller))
	k.mustSetPositionInts(sellerPos, sdkmath.NewInt(10), sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt())
	k.SetPosition(ctx, sellerPos)

	buyOrderID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{Trader: buyer, MarketId: marketID, Side: types.OrderSide_ORDER_SIDE_BUY_YES, OrderType: types.OrderType_ORDER_TYPE_LIMIT, Amount: "1", LimitPrice: "10000"})
	sellOrderID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{Trader: seller, MarketId: marketID, Side: types.OrderSide_ORDER_SIDE_SELL_YES, OrderType: types.OrderType_ORDER_TYPE_LIMIT, Amount: "1", LimitPrice: "10000"})

	_, err := k.ApplyTradeBatch(ctx, &types.MsgApplyTradeBatch{
		Sender:   resolver,
		MarketId: marketID,
		BatchId:  "batch-insufficient-balance",
		Trades: []types.TradeMatch{{
			TradeId:        "t-wallet-1",
			OrderAId:       buyOrderID,
			OrderBId:       sellOrderID,
			MatchAmount:    "1",
			ExecutionPrice: "10000",
		}},
	})
	require.ErrorContains(t, err, "insufficient funds")
}

func TestSettlementInsufficientYesShares(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(61)
	resolver := testAddress(62)
	buyer := testAddress(64)
	seller := testAddress(65)

	params := k.GetParams(ctx)
	k.SetParams(ctx, params)

	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)
	mustFund(bank, buyer, 20_000)

	buyOrderID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{Trader: buyer, MarketId: marketID, Side: types.OrderSide_ORDER_SIDE_BUY_YES, OrderType: types.OrderType_ORDER_TYPE_LIMIT, Amount: "1", LimitPrice: "10000"})
	sellOrderID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{Trader: seller, MarketId: marketID, Side: types.OrderSide_ORDER_SIDE_SELL_YES, OrderType: types.OrderType_ORDER_TYPE_LIMIT, Amount: "1", LimitPrice: "10000"})

	_, err := k.ApplyTradeBatch(ctx, &types.MsgApplyTradeBatch{
		Sender:   resolver,
		MarketId: marketID,
		BatchId:  "batch-insufficient-yes",
		Trades: []types.TradeMatch{{
			TradeId:        "t-yes-1",
			OrderAId:       buyOrderID,
			OrderBId:       sellOrderID,
			MatchAmount:    "1",
			ExecutionPrice: "10000",
		}},
	})
	require.ErrorContains(t, err, "seller YES shares")
}

func TestSettlementInsufficientNoShares(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(71)
	resolver := testAddress(72)
	buyer := testAddress(74)
	seller := testAddress(75)

	params := k.GetParams(ctx)
	k.SetParams(ctx, params)

	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)
	mustFund(bank, buyer, 20_000)

	buyOrderID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{Trader: buyer, MarketId: marketID, Side: types.OrderSide_ORDER_SIDE_BUY_NO, OrderType: types.OrderType_ORDER_TYPE_LIMIT, Amount: "1", LimitPrice: "10000"})
	sellOrderID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{Trader: seller, MarketId: marketID, Side: types.OrderSide_ORDER_SIDE_SELL_NO, OrderType: types.OrderType_ORDER_TYPE_LIMIT, Amount: "1", LimitPrice: "10000"})

	_, err := k.ApplyTradeBatch(ctx, &types.MsgApplyTradeBatch{
		Sender:   resolver,
		MarketId: marketID,
		BatchId:  "batch-insufficient-no",
		Trades: []types.TradeMatch{{
			TradeId:        "t-no-1",
			OrderAId:       buyOrderID,
			OrderBId:       sellOrderID,
			MatchAmount:    "1",
			ExecutionPrice: "10000",
		}},
	})
	require.ErrorContains(t, err, "seller NO shares")
}

func TestTradeAfterMarketCloseRejected(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(81)
	resolver := testAddress(82)
	buyer := testAddress(84)
	seller := testAddress(85)

	params := k.GetParams(ctx)
	k.SetParams(ctx, params)

	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)
	mustFund(bank, buyer, 2_000)
	sellerPos := k.getPositionOrDefault(ctx, marketID, sdk.MustAccAddressFromBech32(seller))
	k.mustSetPositionInts(sellerPos, sdkmath.NewInt(1), sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt())
	k.SetPosition(ctx, sellerPos)

	buyOrderID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{Trader: buyer, MarketId: marketID, Side: types.OrderSide_ORDER_SIDE_BUY_YES, OrderType: types.OrderType_ORDER_TYPE_LIMIT, Amount: "1", LimitPrice: "10000"})
	sellOrderID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{Trader: seller, MarketId: marketID, Side: types.OrderSide_ORDER_SIDE_SELL_YES, OrderType: types.OrderType_ORDER_TYPE_LIMIT, Amount: "1", LimitPrice: "10000"})

	market, _ := k.GetMarket(ctx, marketID)
	market.CloseTime = ctx.BlockTime().Unix() - 1
	k.SetMarket(ctx, market)

	_, err := k.ApplyTradeBatch(ctx, &types.MsgApplyTradeBatch{
		Sender:   resolver,
		MarketId: marketID,
		BatchId:  "batch-closed",
		Trades: []types.TradeMatch{{
			TradeId:        "t-closed-1",
			OrderAId:       buyOrderID,
			OrderBId:       sellOrderID,
			MatchAmount:    "1",
			ExecutionPrice: "10000",
		}},
	})
	require.ErrorContains(t, err, "market must be OPEN")
}

func TestApplyTradeBatchOnlyResolverAuthorized(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(86)
	resolver := testAddress(87)
	attacker := testAddress(88)

	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)

	_, err := k.ApplyTradeBatch(ctx, &types.MsgApplyTradeBatch{
		Sender:   attacker,
		MarketId: marketID,
		BatchId:  "batch-unauthorized",
		Trades: []types.TradeMatch{{
			TradeId:        "t-unauthorized-1",
			OrderAId:       1,
			OrderBId:       2,
			MatchAmount:    "1",
			ExecutionPrice: "10000",
		}},
	})
	require.ErrorIs(t, err, types.ErrUnauthorized)
}

func TestResolveMarketAndResolverOnlyResolutionSourceUpdate(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(91)
	resolver := testAddress(92)
	attacker := testAddress(93)
	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)

	market, _ := k.GetMarket(ctx, marketID)
	market.ResolveTime = ctx.BlockTime().Unix()
	k.SetMarket(ctx, market)

	err := k.ResolveMarket(ctx, &types.MsgResolveMarket{
		Resolver:         attacker,
		MarketId:         marketID,
		WinningOutcome:   types.Outcome_OUTCOME_YES,
		ResolutionSource: "bad",
	})
	require.ErrorIs(t, err, types.ErrUnauthorized)

	err = k.ResolveMarket(ctx, &types.MsgResolveMarket{
		Resolver:         resolver,
		MarketId:         marketID,
		WinningOutcome:   types.Outcome_OUTCOME_NO,
		ResolutionSource: "official_source",
	})
	require.NoError(t, err)

	market, _ = k.GetMarket(ctx, marketID)
	require.Equal(t, types.MarketStatus_MARKET_STATUS_RESOLVED, market.Status)
	require.Equal(t, types.Outcome_OUTCOME_NO, market.WinningOutcome)
	require.Equal(t, "official_source", market.ResolutionSource)
}

func TestClaimPayoutSuccessAndDoubleClaim(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(101)
	resolver := testAddress(102)
	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)

	market, _ := k.GetMarket(ctx, marketID)
	market.Status = types.MarketStatus_MARKET_STATUS_RESOLVED
	market.WinningOutcome = types.Outcome_OUTCOME_YES
	market.TotalYesShares = "5"
	market.TotalNoShares = "0"
	k.SetMarket(ctx, market)

	pos := k.getPositionOrDefault(ctx, marketID, sdk.MustAccAddressFromBech32(creator))
	k.mustSetPositionInts(pos, sdkmath.NewInt(5), sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt())
	k.SetPosition(ctx, pos)

	before := bank.AccountBalance(creator, "upaxi")
	resp, err := k.ClaimPayout(ctx, &types.MsgClaimPayout{Creator: creator, MarketId: marketID})
	require.NoError(t, err)
	require.Equal(t, "5", resp.Payout)
	require.Equal(t, before.AddRaw(5), bank.AccountBalance(creator, "upaxi"))

	_, err = k.ClaimPayout(ctx, &types.MsgClaimPayout{Creator: creator, MarketId: marketID})
	require.ErrorIs(t, err, types.ErrAlreadyClaimed)
}

func TestClaimVoidRefundSuccessAndDoubleClaim(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(111)
	resolver := testAddress(112)
	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)

	market, _ := k.GetMarket(ctx, marketID)
	market.Status = types.MarketStatus_MARKET_STATUS_VOIDED
	market.TotalYesShares = "6"
	market.TotalNoShares = "2"
	k.SetMarket(ctx, market)

	pos := k.getPositionOrDefault(ctx, marketID, sdk.MustAccAddressFromBech32(creator))
	k.mustSetPositionInts(pos, sdkmath.NewInt(6), sdkmath.ZeroInt(), sdkmath.NewInt(2), sdkmath.ZeroInt())
	k.SetPosition(ctx, pos)

	before := bank.AccountBalance(creator, "upaxi")
	resp, err := k.ClaimVoidRefund(ctx, &types.MsgClaimVoidRefund{Creator: creator, MarketId: marketID})
	require.NoError(t, err)
	require.Equal(t, "4", resp.Refund)
	require.Equal(t, before.AddRaw(4), bank.AccountBalance(creator, "upaxi"))

	_, err = k.ClaimVoidRefund(ctx, &types.MsgClaimVoidRefund{Creator: creator, MarketId: marketID})
	require.ErrorIs(t, err, types.ErrAlreadyClaimed)
}

func TestUpdateParamsAuthority(t *testing.T) {
	authority := testAddress(240)
	nonAuthority := testAddress(241)
	k, ctx, _ := setupKeeperWithAuthority(t, authority)
	srv := NewMsgServerImpl(k)

	current := k.GetParams(ctx)
	paramsInput := types.ParamsInput{
		MaxBatchSize:            current.MaxBatchSize + 1,
		CreateMarketBond:        current.CreateMarketBond,
		CreateMarketBondDenom:   current.CreateMarketBondDenom,
		MarketFeeBps:            current.MarketFeeBps,
		ResolverFeeSharePercent: current.ResolverFeeSharePercent,
		MaxOrderLifetimeBh:      current.MaxOrderLifetimeBh,
		MaxOpenOrdersPerUser:    current.MaxOpenOrdersPerUser,
		MaxOpenOrdersPerMarket:  current.MaxOpenOrdersPerMarket,
		OrderPruneIntervalBh:    current.OrderPruneIntervalBh,
		OrderPruneRetainBh:      current.OrderPruneRetainBh,
		OrderPruneScanLimit:     current.OrderPruneScanLimit,
		OrderPruneDeleteLimit:   current.OrderPruneDeleteLimit,
	}

	_, err := srv.UpdateParams(sdk.WrapSDKContext(ctx), &types.MsgUpdateParams{
		Authority: nonAuthority,
		Params:    paramsInput,
	})
	require.ErrorIs(t, err, types.ErrUnauthorized)

	_, err = srv.UpdateParams(sdk.WrapSDKContext(ctx), &types.MsgUpdateParams{
		Authority: authority,
		Params:    paramsInput,
	})
	require.NoError(t, err)

	updated := k.GetParams(ctx)
	require.Equal(t, current.MaxBatchSize+1, updated.MaxBatchSize)
}
