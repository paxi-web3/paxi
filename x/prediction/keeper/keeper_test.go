package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

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

func (m *mockBankKeeper) setAccountBalance(addr, denom string, amount sdkmath.Int) {
	m.accountBalances[balanceKey(addr, denom)] = amount
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

func (m *mockBankKeeper) GetBalance(_ context.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	return sdk.NewCoin(denom, m.AccountBalance(addr.String(), denom))
}

type mockAccountKeeper struct{}

func (m mockAccountKeeper) GetModuleAddress(moduleName string) sdk.AccAddress {
	addr := make([]byte, 20)
	copy(addr, []byte(moduleName))
	return sdk.AccAddress(addr)
}

type mockPRC20QueryKeeper struct {
	allowances map[string]sdkmath.Int
	balances   map[string]sdkmath.Int
}

func newMockPRC20QueryKeeper() *mockPRC20QueryKeeper {
	return &mockPRC20QueryKeeper{
		allowances: make(map[string]sdkmath.Int),
		balances:   make(map[string]sdkmath.Int),
	}
}

func allowanceKey(contract, owner, spender string) string {
	return contract + "|" + owner + "|" + spender
}

func (m *mockPRC20QueryKeeper) setAllowance(contract, owner, spender string, allowance sdkmath.Int) {
	m.allowances[allowanceKey(contract, owner, spender)] = allowance
}

func (m *mockPRC20QueryKeeper) setBalance(contract, owner string, balance sdkmath.Int) {
	m.balances[contract+"|"+owner] = balance
}

func (m *mockPRC20QueryKeeper) QuerySmart(_ context.Context, contractAddress sdk.AccAddress, req []byte) ([]byte, error) {
	raw := map[string]json.RawMessage{}
	if err := json.Unmarshal(req, &raw); err != nil {
		return nil, err
	}
	if allowanceRaw, ok := raw["allowance"]; ok {
		var allowanceReq struct {
			Owner   string `json:"owner"`
			Spender string `json:"spender"`
		}
		if err := json.Unmarshal(allowanceRaw, &allowanceReq); err != nil {
			return nil, err
		}
		allowance, ok := m.allowances[allowanceKey(contractAddress.String(), allowanceReq.Owner, allowanceReq.Spender)]
		if !ok {
			allowance = sdkmath.ZeroInt()
		}
		return json.Marshal(prc20AllowanceResponse{Allowance: allowance.String()})
	}
	if balanceRaw, ok := raw["balance"]; ok {
		var balanceReq struct {
			Address string `json:"address"`
		}
		if err := json.Unmarshal(balanceRaw, &balanceReq); err != nil {
			return nil, err
		}
		balance, ok := m.balances[contractAddress.String()+"|"+balanceReq.Address]
		if !ok {
			balance = sdkmath.ZeroInt()
		}
		return json.Marshal(prc20BalanceResponse{Balance: balance.String()})
	}
	return nil, fmt.Errorf("unsupported query")
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

func TestCreateMarketTextCharLimits(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(220)
	resolver := testAddress(221)
	mustFund(bank, creator, 1_000_000)
	now := ctx.BlockTime().Unix()

	base := types.MsgCreateMarket{
		Creator:         creator,
		Resolver:        resolver,
		Title:           "T",
		Description:     "D",
		Rule:            "R",
		OutcomeType:     "BINARY",
		Outcomes:        []string{"YES", "NO"},
		CollateralType:  types.CollateralType_COLLATERAL_TYPE_NATIVE,
		CollateralDenom: "upaxi",
		OpenTime:        now - 10,
		CloseTime:       now + 3600,
		ResolveTime:     now + 7200,
	}

	tests := []struct {
		name    string
		apply   func(m *types.MsgCreateMarket)
		wantErr string
	}{
		{
			name: "title at max-1 chars is valid",
			apply: func(m *types.MsgCreateMarket) {
				m.Title = strings.Repeat("a", types.MaxMarketTitleChars-1)
			},
		},
		{
			name: "title at max chars is rejected",
			apply: func(m *types.MsgCreateMarket) {
				m.Title = strings.Repeat("a", types.MaxMarketTitleChars)
			},
			wantErr: "title must be < 512 characters",
		},
		{
			name: "description at max chars is rejected",
			apply: func(m *types.MsgCreateMarket) {
				m.Description = strings.Repeat("b", types.MaxMarketDescriptionChars)
			},
			wantErr: "description must be < 4096 characters",
		},
		{
			name: "rule at max chars is rejected",
			apply: func(m *types.MsgCreateMarket) {
				m.Rule = strings.Repeat("c", types.MaxMarketRuleChars)
			},
			wantErr: "rule must be < 4096 characters",
		},
		{
			name: "multibyte title counts characters not bytes",
			apply: func(m *types.MsgCreateMarket) {
				m.Title = strings.Repeat("你", types.MaxMarketTitleChars-1)
			},
		},
		{
			name: "multibyte title at max chars is rejected",
			apply: func(m *types.MsgCreateMarket) {
				m.Title = strings.Repeat("你", types.MaxMarketTitleChars)
			},
			wantErr: "title must be < 512 characters",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			msg := base
			tc.apply(&msg)
			_, err := k.CreateMarket(ctx, &msg)
			if tc.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.wantErr)
			}
		})
	}
}

func TestPlaceOrderWithoutEscrow(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(1)
	resolver := testAddress(2)
	trader := testAddress(3)
	mustFund(bank, trader, 200_000)

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

func TestSetMarketUpdatesCloseIndexOnCloseTimeChange(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(41)
	resolver := testAddress(42)
	mustFund(bank, creator, 1_000_000)
	now := ctx.BlockTime().Unix()

	marketID, err := k.CreateMarket(ctx, &types.MsgCreateMarket{
		Creator:         creator,
		Resolver:        resolver,
		Title:           "index refresh",
		Description:     "index refresh",
		Rule:            "rule",
		OutcomeType:     "BINARY",
		Outcomes:        []string{"YES", "NO"},
		CollateralType:  types.CollateralType_COLLATERAL_TYPE_NATIVE,
		CollateralDenom: "upaxi",
		OpenTime:        now - 10,
		CloseTime:       now + 100,
		ResolveTime:     now + 200,
	})
	require.NoError(t, err)

	market, found := k.GetMarket(ctx, marketID)
	require.True(t, found)
	oldClose := market.CloseTime
	oldKey := types.MarketCloseIndexKey(oldClose, marketID)

	store := k.storeService.OpenKVStore(ctx)
	oldIndex, err := store.Get(oldKey)
	require.NoError(t, err)
	require.NotNil(t, oldIndex)

	market.CloseTime = now + 300
	k.SetMarket(ctx, market)

	oldIndex, err = store.Get(oldKey)
	require.NoError(t, err)
	require.Nil(t, oldIndex)
	newIndex, err := store.Get(types.MarketCloseIndexKey(market.CloseTime, marketID))
	require.NoError(t, err)
	require.NotNil(t, newIndex)
}

func TestAutoCloseExpiredMarketsUsesCloseIndex(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(43)
	resolver := testAddress(44)
	mustFund(bank, creator, 2_000_000)
	now := ctx.BlockTime().Unix()

	create := func(title string, closeTime int64) uint64 {
		marketID, err := k.CreateMarket(ctx, &types.MsgCreateMarket{
			Creator:         creator,
			Resolver:        resolver,
			Title:           title,
			Description:     "auto close",
			Rule:            "rule",
			OutcomeType:     "BINARY",
			Outcomes:        []string{"YES", "NO"},
			CollateralType:  types.CollateralType_COLLATERAL_TYPE_NATIVE,
			CollateralDenom: "upaxi",
			OpenTime:        now - 10,
			CloseTime:       closeTime,
			ResolveTime:     closeTime + 100,
		})
		require.NoError(t, err)
		return marketID
	}

	expiredID := create("expired", now+1)
	openID := create("open", now+100)

	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(2 * time.Second))
	k.AutoCloseExpiredMarkets(ctx)

	expiredMarket, found := k.GetMarket(ctx, expiredID)
	require.True(t, found)
	require.Equal(t, types.MarketStatus_MARKET_STATUS_CLOSED, expiredMarket.Status)

	openMarket, found := k.GetMarket(ctx, openID)
	require.True(t, found)
	require.Equal(t, types.MarketStatus_MARKET_STATUS_OPEN, openMarket.Status)

	store := k.storeService.OpenKVStore(ctx)
	expiredIndex, err := store.Get(types.MarketCloseIndexKey(expiredMarket.CloseTime, expiredID))
	require.NoError(t, err)
	require.Nil(t, expiredIndex)

	openIndex, err := store.Get(types.MarketCloseIndexKey(openMarket.CloseTime, openID))
	require.NoError(t, err)
	require.NotNil(t, openIndex)
}

func TestAutoCloseExpiredMarketsBackfillsLegacyIndex(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(45)
	resolver := testAddress(46)
	mustFund(bank, creator, 1_000_000)
	now := ctx.BlockTime().Unix()

	marketID, err := k.CreateMarket(ctx, &types.MsgCreateMarket{
		Creator:         creator,
		Resolver:        resolver,
		Title:           "legacy backfill",
		Description:     "legacy backfill",
		Rule:            "rule",
		OutcomeType:     "BINARY",
		Outcomes:        []string{"YES", "NO"},
		CollateralType:  types.CollateralType_COLLATERAL_TYPE_NATIVE,
		CollateralDenom: "upaxi",
		OpenTime:        now - 10,
		CloseTime:       now + 1,
		ResolveTime:     now + 100,
	})
	require.NoError(t, err)

	market, found := k.GetMarket(ctx, marketID)
	require.True(t, found)

	store := k.storeService.OpenKVStore(ctx)
	k.mustDelete(store, types.MarketCloseIndexKey(market.CloseTime, marketID))
	k.mustDelete(store, types.MarketCloseIdxReady)
	k.mustDelete(store, types.MarketCloseIdxCursor)

	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(2 * time.Second))
	k.AutoCloseExpiredMarkets(ctx)

	updated, found := k.GetMarket(ctx, marketID)
	require.True(t, found)
	require.Equal(t, types.MarketStatus_MARKET_STATUS_CLOSED, updated.Status)

	indexBz, err := store.Get(types.MarketCloseIndexKey(updated.CloseTime, marketID))
	require.NoError(t, err)
	require.Nil(t, indexBz)
}

func TestPlaceOrderRequiresSufficientBalance(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(4)
	resolver := testAddress(5)
	trader := testAddress(6)

	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)
	mustFund(bank, trader, 50_000)

	_, err := k.PlaceOrder(ctx, &types.MsgPlaceOrder{
		Trader:     trader,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "10",
		LimitPrice: "10000",
		ExpireBh:   ctx.BlockHeight() + 100,
	})
	require.ErrorContains(t, err, "insufficient funds")
}

func TestPlaceOrderValidation(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(11)
	resolver := testAddress(12)
	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)
	trader := testAddress(13)
	mustFund(bank, trader, 1_000_000)

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
	prc20Query.setBalance(contract, buyer, sdkmath.NewInt(101_000))
	prc20Query.setAllowance(contract, buyer, moduleAddr.String(), sdkmath.NewInt(100_999))

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

	prc20Query.setAllowance(contract, buyer, moduleAddr.String(), sdkmath.NewInt(101_000))

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
	mustFund(bank, trader, 200_000)
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
	mustFund(bank, trader, 200_000)
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

func TestAutoPruneOrdersPrunesStaleUnusedOpenOrders(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	ctx = ctx.WithBlockHeight(300)

	creator := testAddress(104)
	resolver := testAddress(105)
	trader := testAddress(106)
	mustFund(bank, trader, 200_000)
	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)

	params := k.GetParams(ctx)
	params.OrderPruneIntervalBh = 1
	params.OrderPruneRetainBh = 2
	params.OrderPruneScanLimit = 100
	params.OrderPruneDeleteLimit = 100
	k.SetParams(ctx, params)

	// This order is never matched/cancelled and remains OPEN in storage, but
	// becomes effectively closed after expiry by height.
	staleOpenOrderID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     trader,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "1",
		LimitPrice: "10000",
		ExpireBh:   ctx.BlockHeight() + 1,
	})

	ctx = ctx.WithBlockHeight(303) // threshold = 301, order created at 300
	k.AutoPruneOrders(ctx)

	_, found := k.GetOrder(ctx, marketID, staleOpenOrderID)
	require.False(t, found)
	_, found = k.GetOrderByID(ctx, staleOpenOrderID)
	require.False(t, found)
}

func TestAutoPruneOrdersWithCursorAndDeleteLimit(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	ctx = ctx.WithBlockHeight(200)

	creator := testAddress(111)
	resolver := testAddress(112)
	trader := testAddress(113)
	mustFund(bank, trader, 200_000)
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

func TestApplyTradeBatchBuyYesBuyNoFeeChargedOnTop(t *testing.T) {
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
	mustFund(bank, yesBuyer, 600_000)
	mustFund(bank, noBuyer, 600_000)

	buyYesID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     yesBuyer,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "1",
		LimitPrice: "500000",
	})
	buyNoID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     noBuyer,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_NO,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "1",
		LimitPrice: "500000",
	})

	_, err := k.ApplyTradeBatch(ctx, &types.MsgApplyTradeBatch{
		Sender:   resolver,
		MarketId: marketID,
		BatchId:  "batch-buybuy-fee",
		Trades: []types.TradeMatch{{
			TradeId:           "t-buybuy-fee-1",
			OrderAId:          buyYesID,
			OrderBId:          buyNoID,
			MatchAmount:       "1",
			YesExecutionPrice: "500000",
			NoExecutionPrice:  "500000",
		}},
	})
	require.NoError(t, err)

	// For BUY_YES <-> BUY_NO, fee is charged on top of notional from each buyer.
	require.Equal(t, sdkmath.NewInt(95_000), bank.AccountBalance(yesBuyer, "upaxi"))
	require.Equal(t, sdkmath.NewInt(95_000), bank.AccountBalance(noBuyer, "upaxi"))
	require.Equal(t, sdkmath.NewInt(10_000), bank.AccountBalance(resolver, "upaxi"))
}

func TestApplyTradeBatchBuyYesBuyNoCostScalesWithMatchAmount(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(26)
	resolver := testAddress(27)
	yesBuyer := testAddress(28)
	noBuyer := testAddress(29)

	params := k.GetParams(ctx)
	params.MarketFeeBps = 0
	k.SetParams(ctx, params)

	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)
	mustFund(bank, yesBuyer, 6_000_000)
	mustFund(bank, noBuyer, 6_000_000)

	buyYesID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     yesBuyer,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "10",
		LimitPrice: "500000",
	})
	buyNoID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     noBuyer,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_NO,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "10",
		LimitPrice: "500000",
	})

	_, err := k.ApplyTradeBatch(ctx, &types.MsgApplyTradeBatch{
		Sender:   resolver,
		MarketId: marketID,
		BatchId:  "batch-buybuy-scale",
		Trades: []types.TradeMatch{{
			TradeId:           "t-buybuy-scale-1",
			OrderAId:          buyYesID,
			OrderBId:          buyNoID,
			MatchAmount:       "10",
			YesExecutionPrice: "500000",
			NoExecutionPrice:  "500000",
		}},
	})
	require.NoError(t, err)

	// Each side pays match_amount * side_execution_price = 10 * 500_000.
	require.Equal(t, sdkmath.NewInt(1_000_000), bank.AccountBalance(yesBuyer, "upaxi"))
	require.Equal(t, sdkmath.NewInt(1_000_000), bank.AccountBalance(noBuyer, "upaxi"))
}

func TestApplyTradeBatchBuyYesBuyNoMarketMarketNotAllowed(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(22)
	resolver := testAddress(23)
	yesBuyer := testAddress(24)
	noBuyer := testAddress(25)

	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)
	mustFund(bank, yesBuyer, 1_000_000)
	mustFund(bank, noBuyer, 1_000_000)

	buyYesID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     yesBuyer,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
		OrderType:  types.OrderType_ORDER_TYPE_MARKET,
		Amount:     "1",
		WorstPrice: "500000",
	})
	buyNoID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     noBuyer,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_NO,
		OrderType:  types.OrderType_ORDER_TYPE_MARKET,
		Amount:     "1",
		WorstPrice: "500000",
	})

	beforeYes := bank.AccountBalance(yesBuyer, "upaxi")
	beforeNo := bank.AccountBalance(noBuyer, "upaxi")

	resp, err := k.ApplyTradeBatch(ctx, &types.MsgApplyTradeBatch{
		Sender:   resolver,
		MarketId: marketID,
		BatchId:  "batch-buybuy-market-market",
		Trades: []types.TradeMatch{{
			TradeId:           "t-buybuy-market-market-1",
			OrderAId:          buyYesID,
			OrderBId:          buyNoID,
			MatchAmount:       "1",
			YesExecutionPrice: "500000",
			NoExecutionPrice:  "500000",
		}},
	})
	require.NoError(t, err)
	require.Equal(t, uint64(0), resp.SettledCount)
	require.False(t, k.HasAppliedTrade(ctx, marketID, "t-buybuy-market-market-1"))

	yesOrder, found := k.GetOrder(ctx, marketID, buyYesID)
	require.True(t, found)
	require.Equal(t, types.OrderStatus_ORDER_STATUS_OPEN, yesOrder.Status)
	require.Equal(t, "0", yesOrder.FilledAmount)
	require.Equal(t, "1", yesOrder.RemainingAmount)

	noOrder, found := k.GetOrder(ctx, marketID, buyNoID)
	require.True(t, found)
	require.Equal(t, types.OrderStatus_ORDER_STATUS_OPEN, noOrder.Status)
	require.Equal(t, "0", noOrder.FilledAmount)
	require.Equal(t, "1", noOrder.RemainingAmount)

	require.Equal(t, beforeYes, bank.AccountBalance(yesBuyer, "upaxi"))
	require.Equal(t, beforeNo, bank.AccountBalance(noBuyer, "upaxi"))
}

func TestApplyTradeBatchBuyYesBuyNoPriceMustSumToOne(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(225)
	resolver := testAddress(226)
	yesBuyer := testAddress(227)
	noBuyer := testAddress(228)

	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)
	mustFund(bank, yesBuyer, 1_000_000)
	mustFund(bank, noBuyer, 1_000_000)

	buyYesID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     yesBuyer,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "1",
		LimitPrice: "600000",
	})
	buyNoID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     noBuyer,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_NO,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "1",
		LimitPrice: "600000",
	})

	beforeYes := bank.AccountBalance(yesBuyer, "upaxi")
	beforeNo := bank.AccountBalance(noBuyer, "upaxi")

	resp, err := k.ApplyTradeBatch(ctx, &types.MsgApplyTradeBatch{
		Sender:   resolver,
		MarketId: marketID,
		BatchId:  "batch-buybuy-price-not-one",
		Trades: []types.TradeMatch{{
			TradeId:           "t-buybuy-price-not-one-1",
			OrderAId:          buyYesID,
			OrderBId:          buyNoID,
			MatchAmount:       "1",
			YesExecutionPrice: "600000",
			NoExecutionPrice:  "600000",
		}},
	})
	require.NoError(t, err)
	require.Equal(t, uint64(0), resp.SettledCount)
	require.False(t, k.HasAppliedTrade(ctx, marketID, "t-buybuy-price-not-one-1"))

	yesOrder, found := k.GetOrder(ctx, marketID, buyYesID)
	require.True(t, found)
	require.Equal(t, types.OrderStatus_ORDER_STATUS_OPEN, yesOrder.Status)

	noOrder, found := k.GetOrder(ctx, marketID, buyNoID)
	require.True(t, found)
	require.Equal(t, types.OrderStatus_ORDER_STATUS_OPEN, noOrder.Status)

	require.Equal(t, beforeYes, bank.AccountBalance(yesBuyer, "upaxi"))
	require.Equal(t, beforeNo, bank.AccountBalance(noBuyer, "upaxi"))
}

func TestApplyTradeBatchBuyYesBuyNoDualExecutionPrice(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(229)
	resolver := testAddress(230)
	yesBuyer := testAddress(231)
	noBuyer := testAddress(232)

	params := k.GetParams(ctx)
	params.MarketFeeBps = 0
	k.SetParams(ctx, params)

	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)
	mustFund(bank, yesBuyer, 1_000_000)
	mustFund(bank, noBuyer, 1_000_000)

	buyYesID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     yesBuyer,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "1",
		LimitPrice: "600000",
	})
	buyNoID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     noBuyer,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_NO,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "1",
		LimitPrice: "400000",
	})

	resp, err := k.ApplyTradeBatch(ctx, &types.MsgApplyTradeBatch{
		Sender:   resolver,
		MarketId: marketID,
		BatchId:  "batch-buybuy-dual-price",
		Trades: []types.TradeMatch{{
			TradeId:           "t-buybuy-dual-price-1",
			OrderAId:          buyYesID,
			OrderBId:          buyNoID,
			MatchAmount:       "1",
			YesExecutionPrice: "600000",
			NoExecutionPrice:  "400000",
		}},
	})
	require.NoError(t, err)
	require.Equal(t, uint64(1), resp.SettledCount)
	require.Equal(t, sdkmath.NewInt(400_000), bank.AccountBalance(yesBuyer, "upaxi"))
	require.Equal(t, sdkmath.NewInt(600_000), bank.AccountBalance(noBuyer, "upaxi"))

	yesOrder, found := k.GetOrder(ctx, marketID, buyYesID)
	require.True(t, found)
	require.Equal(t, "600000", yesOrder.SpentCollateral)
	require.Equal(t, "0", yesOrder.ReceivedCollateral)

	noOrder, found := k.GetOrder(ctx, marketID, buyNoID)
	require.True(t, found)
	require.Equal(t, "400000", noOrder.SpentCollateral)
	require.Equal(t, "0", noOrder.ReceivedCollateral)
}

func TestApplyTradeBatchTracksOrderSpentAndReceivedCollateral(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(120)
	resolver := testAddress(121)
	buyer := testAddress(122)
	seller := testAddress(123)

	params := k.GetParams(ctx)
	params.MarketFeeBps = 100
	params.ResolverFeeSharePercent = 100
	k.SetParams(ctx, params)

	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)
	mustFund(bank, buyer, 2_000_000)

	sellerPos := k.getPositionOrDefault(ctx, marketID, sdk.MustAccAddressFromBech32(seller))
	k.mustSetPositionInts(sellerPos, sdkmath.NewInt(2), sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt())
	k.SetPosition(ctx, sellerPos)

	buyOrderID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     buyer,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "2",
		LimitPrice: "500000",
	})
	sellOrderID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     seller,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_SELL_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "2",
		LimitPrice: "500000",
	})

	_, err := k.ApplyTradeBatch(ctx, &types.MsgApplyTradeBatch{
		Sender:   resolver,
		MarketId: marketID,
		BatchId:  "batch-order-collateral-tracking",
		Trades: []types.TradeMatch{{
			TradeId:           "t-order-collateral-tracking-1",
			OrderAId:          buyOrderID,
			OrderBId:          sellOrderID,
			MatchAmount:       "2",
			YesExecutionPrice: "500000",
			NoExecutionPrice:  "500000",
		}},
	})
	require.NoError(t, err)

	buyOrder, found := k.GetOrder(ctx, marketID, buyOrderID)
	require.True(t, found)
	require.Equal(t, "1000000", buyOrder.SpentCollateral)
	require.Equal(t, "0", buyOrder.ReceivedCollateral)

	sellOrder, found := k.GetOrder(ctx, marketID, sellOrderID)
	require.True(t, found)
	require.Equal(t, "0", sellOrder.SpentCollateral)
	require.Equal(t, "990000", sellOrder.ReceivedCollateral)
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
	mustFund(bank, trader, 2_000_000_000)

	moduleBefore := bank.ModuleBalance("upaxi")
	traderBefore := bank.AccountBalance(trader, "upaxi")
	shareUnit := sdkmath.NewInt(types.CollateralUnit)
	splitShares := sdkmath.NewInt(1000)
	mergeShares := sdkmath.NewInt(400)
	splitCollateral := splitShares.Mul(shareUnit)
	mergeCollateral := mergeShares.Mul(shareUnit)

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
	require.Equal(t, traderBefore.Sub(splitCollateral), bank.AccountBalance(trader, "upaxi"))
	require.Equal(t, moduleBefore.Add(splitCollateral), bank.ModuleBalance("upaxi"))

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
	require.Equal(t, traderBefore.Sub(splitCollateral).Add(mergeCollateral), bank.AccountBalance(trader, "upaxi"))
	require.Equal(t, moduleBefore.Add(splitCollateral).Sub(mergeCollateral), bank.ModuleBalance("upaxi"))
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
	mustFund(bank, trader, 200_000_000)

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
	mustFund(bank, buyer, 500_000)
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
			TradeId:           "t-1",
			OrderAId:          buyOrderID,
			OrderBId:          sellOrderID,
			MatchAmount:       "10",
			YesExecutionPrice: "10000",
			NoExecutionPrice:  "10000",
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

	require.Equal(t, sdkmath.NewInt(400_000), bank.AccountBalance(buyer, "upaxi"))
	require.Equal(t, sdkmath.NewInt(99_010), bank.AccountBalance(seller, "upaxi"))
	require.Equal(t, sdkmath.NewInt(1_000), bank.AccountBalance(resolver, "upaxi"))

	market, found := k.GetMarket(ctx, marketID)
	require.True(t, found)
	require.Equal(t, "100000", market.TotalTradeVolume)
}

func TestCancelOrder(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(31)
	resolver := testAddress(32)
	trader := testAddress(33)
	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)
	mustFund(bank, trader, 100_000)

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
	mustFund(bank, buyer, 300_000)

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
			TradeId:           "partial-cancel-trade",
			OrderAId:          buyOrderID,
			OrderBId:          sellOrderID,
			MatchAmount:       "10",
			YesExecutionPrice: "10000",
			NoExecutionPrice:  "10000",
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
	mustFund(bank, trader, 1_000_000)
	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)
	pos := k.getPositionOrDefault(ctx, marketID, sdk.MustAccAddressFromBech32(trader))
	k.mustSetPositionInts(pos, sdkmath.NewInt(10), sdkmath.ZeroInt(), sdkmath.NewInt(10), sdkmath.ZeroInt())
	k.SetPosition(ctx, pos)

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

func TestQueryOrdersByMarket(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(221)
	resolver := testAddress(222)
	trader := testAddress(223)
	other := testAddress(224)
	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)
	otherMarketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)
	mustFund(bank, trader, 1_000_000)
	mustFund(bank, other, 1_000_000)

	mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     trader,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "1",
		LimitPrice: "20000",
	})
	mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     other,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_NO,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "2",
		LimitPrice: "30000",
	})
	mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     trader,
		MarketId:   otherMarketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "3",
		LimitPrice: "40000",
	})

	qs := NewQueryServer(k)
	resp, err := qs.OrdersByMarket(ctx, &types.QueryOrdersByMarketRequest{MarketId: marketID})
	require.NoError(t, err)
	require.Len(t, resp.Orders, 2)
	for _, order := range resp.Orders {
		require.Equal(t, marketID, order.MarketId)
	}
}

func TestDuplicateTradeIDSkipped(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(41)
	resolver := testAddress(42)
	buyer := testAddress(44)
	seller := testAddress(45)

	params := k.GetParams(ctx)
	k.SetParams(ctx, params)

	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)
	mustFund(bank, buyer, 50_000)

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
			TradeId:           "dup-trade",
			OrderAId:          buyOrderID,
			OrderBId:          sellOrderID,
			MatchAmount:       "1",
			YesExecutionPrice: "10000",
			NoExecutionPrice:  "10000",
		}},
	})
	require.NoError(t, err)

	resp, err := k.ApplyTradeBatch(ctx, &types.MsgApplyTradeBatch{
		Sender:   resolver,
		MarketId: marketID,
		BatchId:  "batch-2",
		Trades: []types.TradeMatch{{
			TradeId:           "dup-trade",
			OrderAId:          buyOrderID,
			OrderBId:          sellOrderID,
			MatchAmount:       "1",
			YesExecutionPrice: "10000",
			NoExecutionPrice:  "10000",
		}},
	})
	require.NoError(t, err)
	require.Equal(t, uint64(0), resp.SettledCount)

	buyOrder, found := k.GetOrder(ctx, marketID, buyOrderID)
	require.True(t, found)
	require.Equal(t, "1", buyOrder.FilledAmount)
	require.Equal(t, "1", buyOrder.RemainingAmount)
}

func TestApplyTradeBatchUpdatesLastOutcomeTradePriceForYesPair(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(214)
	resolver := testAddress(215)
	buyer := testAddress(216)
	seller := testAddress(217)

	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)
	mustFund(bank, buyer, 40_000)

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
			TradeId:           "t-last-price",
			OrderAId:          buyOrderID,
			OrderBId:          sellOrderID,
			MatchAmount:       "1",
			YesExecutionPrice: "20000",
			NoExecutionPrice:  "20000",
		}},
	})
	require.NoError(t, err)

	market, found := k.GetMarket(ctx, marketID)
	require.True(t, found)
	require.Equal(t, "20000", market.LastYesTradePrice)
	require.Equal(t, "980000", market.LastNoTradePrice)
	require.Equal(t, "20000", market.TotalTradeVolume)
	require.Empty(t, market.BestBidPrice)
	require.Empty(t, market.BestAskPrice)
}

func TestApplyTradeBatchUpdatesLastOutcomeTradePricesBuyYesBuyNo(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(231)
	resolver := testAddress(232)
	yesBuyer := testAddress(233)
	noBuyer := testAddress(234)

	params := k.GetParams(ctx)
	params.MarketFeeBps = 0
	k.SetParams(ctx, params)

	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)
	mustFund(bank, yesBuyer, 600_000)
	mustFund(bank, noBuyer, 600_000)

	buyYesID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     yesBuyer,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "1",
		LimitPrice: "500000",
	})
	buyNoID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     noBuyer,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_NO,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "1",
		LimitPrice: "500000",
	})

	_, err := k.ApplyTradeBatch(ctx, &types.MsgApplyTradeBatch{
		Sender:   resolver,
		MarketId: marketID,
		BatchId:  "batch-last-outcome-prices",
		Trades: []types.TradeMatch{{
			TradeId:           "t-last-outcome-prices-1",
			OrderAId:          buyYesID,
			OrderBId:          buyNoID,
			MatchAmount:       "1",
			YesExecutionPrice: "500000",
			NoExecutionPrice:  "500000",
		}},
	})
	require.NoError(t, err)

	market, found := k.GetMarket(ctx, marketID)
	require.True(t, found)
	require.Equal(t, "500000", market.LastYesTradePrice)
	require.Equal(t, "500000", market.LastNoTradePrice)
}

func TestApplyTradeBatchCanonicalOutcomeLastPricesWeightedAcrossTrades(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(241)
	resolver := testAddress(242)
	yesBuyer := testAddress(243)
	noBuyer := testAddress(244)
	yesSeller := testAddress(245)
	noSeller := testAddress(246)

	params := k.GetParams(ctx)
	params.MarketFeeBps = 0
	k.SetParams(ctx, params)

	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)
	mustFund(bank, yesBuyer, 2_000_000)
	mustFund(bank, noBuyer, 500_000)

	yesSellerPos := k.getPositionOrDefault(ctx, marketID, sdk.MustAccAddressFromBech32(yesSeller))
	k.mustSetPositionInts(yesSellerPos, sdkmath.NewInt(2), sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt())
	k.SetPosition(ctx, yesSellerPos)

	noSellerPos := k.getPositionOrDefault(ctx, marketID, sdk.MustAccAddressFromBech32(noSeller))
	k.mustSetPositionInts(noSellerPos, sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.NewInt(1), sdkmath.ZeroInt())
	k.SetPosition(ctx, noSellerPos)

	buyYesID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     yesBuyer,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "2",
		LimitPrice: "600000",
	})
	sellYesID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     yesSeller,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_SELL_YES,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "2",
		LimitPrice: "10000",
	})
	buyNoID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     noBuyer,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_BUY_NO,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "1",
		LimitPrice: "200000",
	})
	sellNoID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
		Trader:     noSeller,
		MarketId:   marketID,
		Side:       types.OrderSide_ORDER_SIDE_SELL_NO,
		OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
		Amount:     "1",
		LimitPrice: "10000",
	})

	_, err := k.ApplyTradeBatch(ctx, &types.MsgApplyTradeBatch{
		Sender:   resolver,
		MarketId: marketID,
		BatchId:  "batch-canonical-last-prices",
		Trades: []types.TradeMatch{
			{
				TradeId:           "t-canonical-yes-1",
				OrderAId:          buyYesID,
				OrderBId:          sellYesID,
				MatchAmount:       "2",
				YesExecutionPrice: "600000",
				NoExecutionPrice:  "600000",
			},
			{
				TradeId:           "t-canonical-no-1",
				OrderAId:          buyNoID,
				OrderBId:          sellNoID,
				MatchAmount:       "1",
				YesExecutionPrice: "200000",
				NoExecutionPrice:  "200000",
			},
		},
	})
	require.NoError(t, err)

	// YES-view weighted price:
	// (600000*2 + (1000000-200000)*1) / 3 = 666666..., rounded to nearest 10000 = 670000.
	market, found := k.GetMarket(ctx, marketID)
	require.True(t, found)
	require.Equal(t, "670000", market.LastYesTradePrice)
	require.Equal(t, "330000", market.LastNoTradePrice)
	lastYes, ok := sdkmath.NewIntFromString(market.LastYesTradePrice)
	require.True(t, ok)
	lastNo, ok := sdkmath.NewIntFromString(market.LastNoTradePrice)
	require.True(t, ok)
	require.Equal(t, sdkmath.NewInt(types.CollateralUnit), lastYes.Add(lastNo))
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
	mustFund(bank, buyer, 20_000)

	sellerPos := k.getPositionOrDefault(ctx, marketID, sdk.MustAccAddressFromBech32(seller))
	k.mustSetPositionInts(sellerPos, sdkmath.NewInt(10), sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt())
	k.SetPosition(ctx, sellerPos)

	buyOrderID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{Trader: buyer, MarketId: marketID, Side: types.OrderSide_ORDER_SIDE_BUY_YES, OrderType: types.OrderType_ORDER_TYPE_LIMIT, Amount: "1", LimitPrice: "10000"})
	sellOrderID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{Trader: seller, MarketId: marketID, Side: types.OrderSide_ORDER_SIDE_SELL_YES, OrderType: types.OrderType_ORDER_TYPE_LIMIT, Amount: "1", LimitPrice: "10000"})
	bank.setAccountBalance(buyer, "upaxi", sdkmath.NewInt(500))

	resp, err := k.ApplyTradeBatch(ctx, &types.MsgApplyTradeBatch{
		Sender:   resolver,
		MarketId: marketID,
		BatchId:  "batch-insufficient-balance",
		Trades: []types.TradeMatch{{
			TradeId:           "t-wallet-1",
			OrderAId:          buyOrderID,
			OrderBId:          sellOrderID,
			MatchAmount:       "1",
			YesExecutionPrice: "10000",
			NoExecutionPrice:  "10000",
		}},
	})
	require.NoError(t, err)
	require.Equal(t, uint64(0), resp.SettledCount)
	_, found := k.GetOrder(ctx, marketID, buyOrderID)
	require.False(t, found, "insufficient buy order should be cancelled and deleted when unfilled")
	sellOrder, found := k.GetOrder(ctx, marketID, sellOrderID)
	require.True(t, found)
	require.Equal(t, types.OrderStatus_ORDER_STATUS_OPEN, sellOrder.Status)
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
	sellerPos := k.getPositionOrDefault(ctx, marketID, sdk.MustAccAddressFromBech32(seller))
	k.mustSetPositionInts(sellerPos, sdkmath.NewInt(1), sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt())
	k.SetPosition(ctx, sellerPos)

	buyOrderID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{Trader: buyer, MarketId: marketID, Side: types.OrderSide_ORDER_SIDE_BUY_YES, OrderType: types.OrderType_ORDER_TYPE_LIMIT, Amount: "1", LimitPrice: "10000"})
	sellOrderID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{Trader: seller, MarketId: marketID, Side: types.OrderSide_ORDER_SIDE_SELL_YES, OrderType: types.OrderType_ORDER_TYPE_LIMIT, Amount: "1", LimitPrice: "10000"})

	// Drain seller YES shares after order placement to simulate race at settlement.
	sellerPos, _ = k.GetPosition(ctx, marketID, sdk.MustAccAddressFromBech32(seller))
	k.mustSetPositionInts(sellerPos, sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt())
	k.SetPosition(ctx, sellerPos)

	resp, err := k.ApplyTradeBatch(ctx, &types.MsgApplyTradeBatch{
		Sender:   resolver,
		MarketId: marketID,
		BatchId:  "batch-insufficient-yes",
		Trades: []types.TradeMatch{{
			TradeId:           "t-yes-1",
			OrderAId:          buyOrderID,
			OrderBId:          sellOrderID,
			MatchAmount:       "1",
			YesExecutionPrice: "10000",
			NoExecutionPrice:  "10000",
		}},
	})
	require.NoError(t, err)
	require.Equal(t, uint64(0), resp.SettledCount)
	_, found := k.GetOrder(ctx, marketID, sellOrderID)
	require.False(t, found, "insufficient sell order should be cancelled and deleted when unfilled")
	buyOrder, found := k.GetOrder(ctx, marketID, buyOrderID)
	require.True(t, found)
	require.Equal(t, types.OrderStatus_ORDER_STATUS_OPEN, buyOrder.Status)
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
	sellerPos := k.getPositionOrDefault(ctx, marketID, sdk.MustAccAddressFromBech32(seller))
	k.mustSetPositionInts(sellerPos, sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.NewInt(1), sdkmath.ZeroInt())
	k.SetPosition(ctx, sellerPos)

	buyOrderID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{Trader: buyer, MarketId: marketID, Side: types.OrderSide_ORDER_SIDE_BUY_NO, OrderType: types.OrderType_ORDER_TYPE_LIMIT, Amount: "1", LimitPrice: "10000"})
	sellOrderID := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{Trader: seller, MarketId: marketID, Side: types.OrderSide_ORDER_SIDE_SELL_NO, OrderType: types.OrderType_ORDER_TYPE_LIMIT, Amount: "1", LimitPrice: "10000"})

	// Drain seller NO shares after order placement to simulate race at settlement.
	sellerPos, _ = k.GetPosition(ctx, marketID, sdk.MustAccAddressFromBech32(seller))
	k.mustSetPositionInts(sellerPos, sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt())
	k.SetPosition(ctx, sellerPos)

	resp, err := k.ApplyTradeBatch(ctx, &types.MsgApplyTradeBatch{
		Sender:   resolver,
		MarketId: marketID,
		BatchId:  "batch-insufficient-no",
		Trades: []types.TradeMatch{{
			TradeId:           "t-no-1",
			OrderAId:          buyOrderID,
			OrderBId:          sellOrderID,
			MatchAmount:       "1",
			YesExecutionPrice: "10000",
			NoExecutionPrice:  "10000",
		}},
	})
	require.NoError(t, err)
	require.Equal(t, uint64(0), resp.SettledCount)
	_, found := k.GetOrder(ctx, marketID, sellOrderID)
	require.False(t, found, "insufficient sell order should be cancelled and deleted when unfilled")
	buyOrder, found := k.GetOrder(ctx, marketID, buyOrderID)
	require.True(t, found)
	require.Equal(t, types.OrderStatus_ORDER_STATUS_OPEN, buyOrder.Status)
}

func TestApplyTradeBatchMixedValidityTable(t *testing.T) {
	type testCase struct {
		name                string
		buildInvalidTrade   func(marketID uint64, validA uint64, validB uint64) types.TradeMatch
		preApply            func(k Keeper, ctx sdk.Context, marketID uint64, validA uint64, validB uint64) sdk.Context
		expectedSettled     uint64
		expectedValidFilled string
	}

	tests := []testCase{
		{
			name: "invalid order id is skipped while valid trade settles",
			buildInvalidTrade: func(_ uint64, _ uint64, validB uint64) types.TradeMatch {
				return types.TradeMatch{
					TradeId:           "invalid-order-not-found",
					OrderAId:          99999999,
					OrderBId:          validB,
					MatchAmount:       "1",
					YesExecutionPrice: "500000",
					NoExecutionPrice:  "500000",
				}
			},
			expectedSettled:     1,
			expectedValidFilled: "1",
		},
		{
			name: "expired order is skipped while valid trade settles",
			buildInvalidTrade: func(_ uint64, validA uint64, validB uint64) types.TradeMatch {
				return types.TradeMatch{
					TradeId:           "invalid-expired-order",
					OrderAId:          validA,
					OrderBId:          validB,
					MatchAmount:       "1",
					YesExecutionPrice: "500000",
					NoExecutionPrice:  "500000",
				}
			},
			preApply: func(k Keeper, ctx sdk.Context, marketID uint64, validA uint64, validB uint64) sdk.Context {
				orderA, found := k.GetOrder(ctx, marketID, validA)
				require.True(t, found)
				orderB, found := k.GetOrder(ctx, marketID, validB)
				require.True(t, found)
				orderA.ExpireBh = 1
				orderB.ExpireBh = 1
				k.SetOrder(ctx, orderA)
				k.SetOrder(ctx, orderB)
				return ctx.WithBlockHeight(1)
			},
			expectedSettled:     1,
			expectedValidFilled: "1",
		},
		{
			name: "duplicate trade id is skipped while valid trade settles",
			buildInvalidTrade: func(marketID uint64, validA uint64, validB uint64) types.TradeMatch {
				return types.TradeMatch{
					TradeId:           "dup-trade-id",
					OrderAId:          validA,
					OrderBId:          validB,
					MatchAmount:       "1",
					YesExecutionPrice: "500000",
					NoExecutionPrice:  "500000",
				}
			},
			preApply: func(k Keeper, ctx sdk.Context, marketID uint64, _ uint64, _ uint64) sdk.Context {
				k.SetAppliedTrade(ctx, marketID, "dup-trade-id")
				return ctx
			},
			expectedSettled:     1,
			expectedValidFilled: "1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			k, ctx, bank := setupKeeper(t)
			creator := testAddress(241)
			resolver := testAddress(242)
			yesBuyer := testAddress(243)
			noBuyer := testAddress(244)
			yesBuyer2 := testAddress(245)
			noBuyer2 := testAddress(246)

			marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)
			mustFund(bank, yesBuyer, 700_000)
			mustFund(bank, noBuyer, 700_000)
			mustFund(bank, yesBuyer2, 700_000)
			mustFund(bank, noBuyer2, 700_000)

			validA := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
				Trader:     yesBuyer,
				MarketId:   marketID,
				Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
				OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
				Amount:     "1",
				LimitPrice: "500000",
				ExpireBh:   ctx.BlockHeight() + 100,
			})
			validB := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
				Trader:     noBuyer,
				MarketId:   marketID,
				Side:       types.OrderSide_ORDER_SIDE_BUY_NO,
				OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
				Amount:     "1",
				LimitPrice: "500000",
				ExpireBh:   ctx.BlockHeight() + 100,
			})

			invalidA := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
				Trader:     yesBuyer2,
				MarketId:   marketID,
				Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
				OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
				Amount:     "1",
				LimitPrice: "500000",
				ExpireBh:   ctx.BlockHeight() + 100,
			})
			invalidB := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
				Trader:     noBuyer2,
				MarketId:   marketID,
				Side:       types.OrderSide_ORDER_SIDE_BUY_NO,
				OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
				Amount:     "1",
				LimitPrice: "500000",
				ExpireBh:   ctx.BlockHeight() + 100,
			})

			invalidTrade := tc.buildInvalidTrade(marketID, invalidA, invalidB)
			if tc.preApply != nil {
				ctx = tc.preApply(k, ctx, marketID, invalidA, invalidB)
			}

			resp, err := k.ApplyTradeBatch(ctx, &types.MsgApplyTradeBatch{
				Sender:   resolver,
				MarketId: marketID,
				BatchId:  "mixed-validity",
				Trades: []types.TradeMatch{
					{
						TradeId:           "valid-trade",
						OrderAId:          validA,
						OrderBId:          validB,
						MatchAmount:       "1",
						YesExecutionPrice: "500000",
						NoExecutionPrice:  "500000",
					},
					invalidTrade,
				},
			})
			require.NoError(t, err)
			require.Equal(t, tc.expectedSettled, resp.SettledCount)

			validOrder, found := k.GetOrder(ctx, marketID, validA)
			require.True(t, found)
			require.Equal(t, tc.expectedValidFilled, validOrder.FilledAmount)
		})
	}
}

func TestApplyTradeBatchSettlementRaceAutoCancelTable(t *testing.T) {
	type testCase struct {
		name                  string
		setupOrdersAndMutate  func(t *testing.T, k Keeper, ctx sdk.Context, bank *mockBankKeeper, marketID uint64) (uint64, uint64)
		executionPrice        string
		expectedRemovedOrderA bool
		expectedRemovedOrderB bool
	}

	tests := []testCase{
		{
			name: "buy_yes_buy_no buyer balance drained at settlement",
			setupOrdersAndMutate: func(t *testing.T, k Keeper, ctx sdk.Context, bank *mockBankKeeper, marketID uint64) (uint64, uint64) {
				yesBuyer := testAddress(251)
				noBuyer := testAddress(252)
				mustFund(bank, yesBuyer, 600_000)
				mustFund(bank, noBuyer, 600_000)

				orderA := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
					Trader:     yesBuyer,
					MarketId:   marketID,
					Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
					OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
					Amount:     "1",
					LimitPrice: "500000",
					ExpireBh:   ctx.BlockHeight() + 100,
				})
				orderB := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
					Trader:     noBuyer,
					MarketId:   marketID,
					Side:       types.OrderSide_ORDER_SIDE_BUY_NO,
					OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
					Amount:     "1",
					LimitPrice: "500000",
					ExpireBh:   ctx.BlockHeight() + 100,
				})

				bank.setAccountBalance(yesBuyer, "upaxi", sdkmath.ZeroInt())
				return orderA, orderB
			},
			executionPrice:        "500000",
			expectedRemovedOrderA: true,
			expectedRemovedOrderB: false,
		},
		{
			name: "buy_yes_sell_yes buyer balance drained at settlement",
			setupOrdersAndMutate: func(t *testing.T, k Keeper, ctx sdk.Context, bank *mockBankKeeper, marketID uint64) (uint64, uint64) {
				buyer := testAddress(253)
				seller := testAddress(254)
				mustFund(bank, buyer, 20_000)
				sellerPos := k.getPositionOrDefault(ctx, marketID, sdk.MustAccAddressFromBech32(seller))
				k.mustSetPositionInts(sellerPos, sdkmath.NewInt(1), sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt())
				k.SetPosition(ctx, sellerPos)

				orderA := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
					Trader:     buyer,
					MarketId:   marketID,
					Side:       types.OrderSide_ORDER_SIDE_BUY_YES,
					OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
					Amount:     "1",
					LimitPrice: "10000",
					ExpireBh:   ctx.BlockHeight() + 100,
				})
				orderB := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
					Trader:     seller,
					MarketId:   marketID,
					Side:       types.OrderSide_ORDER_SIDE_SELL_YES,
					OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
					Amount:     "1",
					LimitPrice: "10000",
					ExpireBh:   ctx.BlockHeight() + 100,
				})

				bank.setAccountBalance(buyer, "upaxi", sdkmath.ZeroInt())
				return orderA, orderB
			},
			executionPrice:        "10000",
			expectedRemovedOrderA: true,
			expectedRemovedOrderB: false,
		},
		{
			name: "buy_no_sell_no seller shares drained at settlement",
			setupOrdersAndMutate: func(t *testing.T, k Keeper, ctx sdk.Context, bank *mockBankKeeper, marketID uint64) (uint64, uint64) {
				buyer := testAddress(55)
				seller := testAddress(56)
				mustFund(bank, buyer, 20_000)
				sellerPos := k.getPositionOrDefault(ctx, marketID, sdk.MustAccAddressFromBech32(seller))
				k.mustSetPositionInts(sellerPos, sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.NewInt(1), sdkmath.ZeroInt())
				k.SetPosition(ctx, sellerPos)

				orderA := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
					Trader:     buyer,
					MarketId:   marketID,
					Side:       types.OrderSide_ORDER_SIDE_BUY_NO,
					OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
					Amount:     "1",
					LimitPrice: "10000",
					ExpireBh:   ctx.BlockHeight() + 100,
				})
				orderB := mustPlaceOrder(t, k, ctx, &types.MsgPlaceOrder{
					Trader:     seller,
					MarketId:   marketID,
					Side:       types.OrderSide_ORDER_SIDE_SELL_NO,
					OrderType:  types.OrderType_ORDER_TYPE_LIMIT,
					Amount:     "1",
					LimitPrice: "10000",
					ExpireBh:   ctx.BlockHeight() + 100,
				})

				sellerPos, _ = k.GetPosition(ctx, marketID, sdk.MustAccAddressFromBech32(seller))
				k.mustSetPositionInts(sellerPos, sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt())
				k.SetPosition(ctx, sellerPos)
				return orderA, orderB
			},
			executionPrice:        "10000",
			expectedRemovedOrderA: false,
			expectedRemovedOrderB: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			k, ctx, bank := setupKeeper(t)
			creator := testAddress(57)
			resolver := testAddress(58)
			marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)

			orderA, orderB := tc.setupOrdersAndMutate(t, k, ctx, bank, marketID)

			resp, err := k.ApplyTradeBatch(ctx, &types.MsgApplyTradeBatch{
				Sender:   resolver,
				MarketId: marketID,
				BatchId:  "race-auto-cancel",
				Trades: []types.TradeMatch{{
					TradeId:           "race-t1",
					OrderAId:          orderA,
					OrderBId:          orderB,
					MatchAmount:       "1",
					YesExecutionPrice: tc.executionPrice,
					NoExecutionPrice:  tc.executionPrice,
				}},
			})
			require.NoError(t, err)
			require.Equal(t, uint64(0), resp.SettledCount)

			_, foundA := k.GetOrder(ctx, marketID, orderA)
			_, foundB := k.GetOrder(ctx, marketID, orderB)
			require.Equal(t, !tc.expectedRemovedOrderA, foundA)
			require.Equal(t, !tc.expectedRemovedOrderB, foundB)
		})
	}
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
	mustFund(bank, buyer, 20_000)
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
			TradeId:           "t-closed-1",
			OrderAId:          buyOrderID,
			OrderBId:          sellOrderID,
			MatchAmount:       "1",
			YesExecutionPrice: "10000",
			NoExecutionPrice:  "10000",
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
			TradeId:           "t-unauthorized-1",
			OrderAId:          1,
			OrderBId:          2,
			MatchAmount:       "1",
			YesExecutionPrice: "10000",
			NoExecutionPrice:  "10000",
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

func TestRequestResolveCreateAndUpdate(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(231)
	resolver := testAddress(232)
	other := testAddress(233)
	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)

	err := k.RequestResolve(ctx, &types.MsgRequestResolve{
		Creator:          other,
		MarketId:         marketID,
		RequestedOutcome: types.Outcome_OUTCOME_YES,
		RequestedSource:  "other-source",
	})
	require.ErrorIs(t, err, types.ErrUnauthorized)

	err = k.RequestResolve(ctx, &types.MsgRequestResolve{
		Creator:          creator,
		MarketId:         marketID,
		RequestedOutcome: types.Outcome_OUTCOME_YES,
		RequestedSource:  "source-v1",
	})
	require.NoError(t, err)

	req, found := k.GetResolveRequest(ctx, marketID)
	require.True(t, found)
	require.Equal(t, creator, req.Creator)
	require.Equal(t, types.Outcome_OUTCOME_YES, req.RequestedOutcome)
	require.Equal(t, "source-v1", req.RequestedSource)

	err = k.RequestResolve(ctx, &types.MsgRequestResolve{
		Creator:          creator,
		MarketId:         marketID,
		RequestedOutcome: types.Outcome_OUTCOME_NO,
		RequestedSource:  "source-v2",
	})
	require.NoError(t, err)

	req, found = k.GetResolveRequest(ctx, marketID)
	require.True(t, found)
	require.Equal(t, types.Outcome_OUTCOME_NO, req.RequestedOutcome)
	require.Equal(t, "source-v2", req.RequestedSource)

	market, _ := k.GetMarket(ctx, marketID)
	market.ResolveTime = 0
	k.SetMarket(ctx, market)

	err = k.ResolveMarket(ctx, &types.MsgResolveMarket{
		Resolver:         resolver,
		MarketId:         marketID,
		WinningOutcome:   types.Outcome_OUTCOME_NO,
		ResolutionSource: "final",
	})
	require.NoError(t, err)

	err = k.RequestResolve(ctx, &types.MsgRequestResolve{
		Creator:          creator,
		MarketId:         marketID,
		RequestedOutcome: types.Outcome_OUTCOME_YES,
		RequestedSource:  "source-after-resolve",
	})
	require.ErrorContains(t, err, "market must be OPEN or CLOSED")
}

func TestResolveMarketWithZeroResolveTime(t *testing.T) {
	k, ctx, bank := setupKeeper(t)
	creator := testAddress(94)
	resolver := testAddress(95)
	marketID := mustCreateMarket(t, k, ctx, bank, creator, resolver)

	market, _ := k.GetMarket(ctx, marketID)
	market.ResolveTime = 0
	k.SetMarket(ctx, market)

	err := k.ResolveMarket(ctx, &types.MsgResolveMarket{
		Resolver:         resolver,
		MarketId:         marketID,
		WinningOutcome:   types.Outcome_OUTCOME_YES,
		ResolutionSource: "instant_source",
	})
	require.NoError(t, err)

	market, _ = k.GetMarket(ctx, marketID)
	require.Equal(t, types.MarketStatus_MARKET_STATUS_RESOLVED, market.Status)
	require.Equal(t, types.Outcome_OUTCOME_YES, market.WinningOutcome)
	require.Equal(t, "instant_source", market.ResolutionSource)
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
	bank.addModuleBalance("upaxi", sdkmath.NewInt(5_000_000))

	before := bank.AccountBalance(creator, "upaxi")
	resp, err := k.ClaimPayout(ctx, &types.MsgClaimPayout{Creator: creator, MarketId: marketID})
	require.NoError(t, err)
	require.Equal(t, "5000000", resp.Payout)
	require.Equal(t, before.AddRaw(5_000_000), bank.AccountBalance(creator, "upaxi"))

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
	bank.addModuleBalance("upaxi", sdkmath.NewInt(4_000_000))

	before := bank.AccountBalance(creator, "upaxi")
	resp, err := k.ClaimVoidRefund(ctx, &types.MsgClaimVoidRefund{Creator: creator, MarketId: marketID})
	require.NoError(t, err)
	require.Equal(t, "4000000", resp.Refund)
	require.Equal(t, before.AddRaw(4_000_000), bank.AccountBalance(creator, "upaxi"))

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
