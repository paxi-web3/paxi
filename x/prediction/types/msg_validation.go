package types

import (
	"fmt"
	"strings"
	"unicode/utf8"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const BPSDenominator = uint64(10_000)

const (
	MaxMarketTitleChars       = 512
	MaxMarketDescriptionChars = 4096
	MaxMarketRuleChars        = 4096
)

func (m *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return fmt.Errorf("invalid authority address: %w", err)
	}

	params := Params{
		MaxBatchSize:            m.Params.MaxBatchSize,
		CreateMarketBond:        m.Params.CreateMarketBond,
		CreateMarketBondDenom:   m.Params.CreateMarketBondDenom,
		MarketFeeBps:            m.Params.MarketFeeBps,
		ResolverFeeSharePercent: m.Params.ResolverFeeSharePercent,
		MaxOrderLifetimeBh:      m.Params.MaxOrderLifetimeBh,
		MaxOpenOrdersPerUser:    m.Params.MaxOpenOrdersPerUser,
		MaxOpenOrdersPerMarket:  m.Params.MaxOpenOrdersPerMarket,
		OrderPruneIntervalBh:    m.Params.OrderPruneIntervalBh,
		OrderPruneRetainBh:      m.Params.OrderPruneRetainBh,
		OrderPruneScanLimit:     m.Params.OrderPruneScanLimit,
		OrderPruneDeleteLimit:   m.Params.OrderPruneDeleteLimit,
	}
	if err := params.Validate(); err != nil {
		return err
	}
	return nil
}

func (m *MsgCreateMarket) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}
	if _, err := sdk.AccAddressFromBech32(m.Resolver); err != nil {
		return fmt.Errorf("invalid resolver address: %w", err)
	}
	if strings.TrimSpace(m.Title) == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if utf8.RuneCountInString(m.Title) >= MaxMarketTitleChars {
		return fmt.Errorf("title must be < %d characters", MaxMarketTitleChars)
	}
	if utf8.RuneCountInString(m.Description) >= MaxMarketDescriptionChars {
		return fmt.Errorf("description must be < %d characters", MaxMarketDescriptionChars)
	}
	if utf8.RuneCountInString(m.Rule) >= MaxMarketRuleChars {
		return fmt.Errorf("rule must be < %d characters", MaxMarketRuleChars)
	}
	if m.OpenTime <= 0 || m.CloseTime <= 0 {
		return fmt.Errorf("open_time and close_time must be > 0")
	}
	if m.OpenTime > m.CloseTime {
		return fmt.Errorf("open_time cannot be greater than close_time")
	}
	if m.ResolveTime < 0 {
		return fmt.Errorf("resolve_time must be >= 0")
	}
	// resolve_time = 0 means no earliest resolve-time restriction.
	if m.ResolveTime > 0 && m.CloseTime > m.ResolveTime {
		return fmt.Errorf("close_time cannot be greater than resolve_time")
	}

	if err := validateCollateralTypeFields(m.CollateralType, m.CollateralDenom, m.CollateralContractAddr); err != nil {
		return err
	}

	if strings.TrimSpace(m.OutcomeType) == "" {
		return fmt.Errorf("outcome_type cannot be empty")
	}
	if len(m.Outcomes) != 2 {
		return fmt.Errorf("binary market must contain exactly two outcomes")
	}

	return nil
}

func (m *MsgPlaceOrder) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Trader); err != nil {
		return fmt.Errorf("invalid trader address: %w", err)
	}
	if m.MarketId == 0 {
		return fmt.Errorf("market_id must be > 0")
	}
	if m.Side == OrderSide_ORDER_SIDE_UNSPECIFIED {
		return fmt.Errorf("side cannot be unspecified")
	}
	if m.OrderType == OrderType_ORDER_TYPE_UNSPECIFIED {
		return fmt.Errorf("order_type cannot be unspecified")
	}

	amount, ok := sdkmath.NewIntFromString(m.Amount)
	if !ok || !amount.IsPositive() {
		return fmt.Errorf("amount must be a positive integer")
	}

	switch m.OrderType {
	case OrderType_ORDER_TYPE_LIMIT:
		if _, err := ParsePriceTicks(m.LimitPrice, "limit_price"); err != nil {
			return err
		}
		if strings.TrimSpace(m.WorstPrice) != "" {
			return fmt.Errorf("worst_price must be empty for limit order")
		}
	case OrderType_ORDER_TYPE_MARKET:
		if _, err := ParsePriceTicks(m.WorstPrice, "worst_price"); err != nil {
			return err
		}
		if strings.TrimSpace(m.LimitPrice) != "" {
			return fmt.Errorf("limit_price must be empty for market order")
		}
	default:
		return fmt.Errorf("invalid order_type")
	}

	if m.ExpireBh <= 0 {
		return fmt.Errorf("expire_bh must be > 0")
	}

	return nil
}

func (m *MsgCancelOrder) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Trader); err != nil {
		return fmt.Errorf("invalid trader address: %w", err)
	}
	if m.MarketId == 0 {
		return fmt.Errorf("market_id must be > 0")
	}
	if m.OrderId == 0 {
		return fmt.Errorf("order_id must be > 0")
	}
	return nil
}

func (m *MsgSplitPosition) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Trader); err != nil {
		return fmt.Errorf("invalid trader address: %w", err)
	}
	if m.MarketId == 0 {
		return fmt.Errorf("market_id must be > 0")
	}
	amount, ok := sdkmath.NewIntFromString(m.Amount)
	if !ok || !amount.IsPositive() {
		return fmt.Errorf("amount must be a positive integer")
	}
	return nil
}

func (m *MsgMergePosition) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Trader); err != nil {
		return fmt.Errorf("invalid trader address: %w", err)
	}
	if m.MarketId == 0 {
		return fmt.Errorf("market_id must be > 0")
	}
	amount, ok := sdkmath.NewIntFromString(m.Amount)
	if !ok || !amount.IsPositive() {
		return fmt.Errorf("amount must be a positive integer")
	}
	return nil
}

func (m *MsgApplyTradeBatch) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return fmt.Errorf("invalid sender address: %w", err)
	}
	if m.MarketId == 0 {
		return fmt.Errorf("market_id must be > 0")
	}
	if strings.TrimSpace(m.BatchId) == "" {
		return fmt.Errorf("batch_id cannot be empty")
	}
	if len(m.Trades) == 0 {
		return fmt.Errorf("trades cannot be empty")
	}
	for i := range m.Trades {
		if err := m.Trades[i].ValidateBasic(); err != nil {
			return fmt.Errorf("invalid trade at index %d: %w", i, err)
		}
	}
	return nil
}

func (m TradeMatch) ValidateBasic() error {
	if strings.TrimSpace(m.TradeId) == "" {
		return fmt.Errorf("trade_id cannot be empty")
	}
	if m.OrderAId == 0 || m.OrderBId == 0 {
		return fmt.Errorf("order ids must be > 0")
	}
	if m.OrderAId == m.OrderBId {
		return fmt.Errorf("order_a_id and order_b_id cannot be the same")
	}
	matchAmount, ok := sdkmath.NewIntFromString(m.MatchAmount)
	if !ok || !matchAmount.IsPositive() {
		return fmt.Errorf("match_amount must be a positive integer")
	}
	if _, err := ParsePriceTicks(m.ExecutionPrice, "execution_price"); err != nil {
		return err
	}
	return nil
}

func (m *MsgResolveMarket) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Resolver); err != nil {
		return fmt.Errorf("invalid resolver address: %w", err)
	}
	if m.MarketId == 0 {
		return fmt.Errorf("market_id must be > 0")
	}
	if m.WinningOutcome != Outcome_OUTCOME_YES && m.WinningOutcome != Outcome_OUTCOME_NO {
		return fmt.Errorf("winning_outcome must be YES or NO")
	}
	return nil
}

func (m *MsgRequestResolve) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}
	if m.MarketId == 0 {
		return fmt.Errorf("market_id must be > 0")
	}
	if m.RequestedOutcome != Outcome_OUTCOME_YES && m.RequestedOutcome != Outcome_OUTCOME_NO {
		return fmt.Errorf("requested_outcome must be YES or NO")
	}
	if strings.TrimSpace(m.RequestedSource) == "" {
		return fmt.Errorf("requested_source cannot be empty")
	}
	return nil
}

func (m *MsgVoidMarket) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Resolver); err != nil {
		return fmt.Errorf("invalid resolver address: %w", err)
	}
	if m.MarketId == 0 {
		return fmt.Errorf("market_id must be > 0")
	}
	return nil
}

func (m *MsgClaimPayout) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}
	if m.MarketId == 0 {
		return fmt.Errorf("market_id must be > 0")
	}
	return nil
}

func (m *MsgClaimVoidRefund) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}
	if m.MarketId == 0 {
		return fmt.Errorf("market_id must be > 0")
	}
	return nil
}

func validateCollateralTypeFields(collateralType CollateralType, denom, contractAddr string) error {
	switch collateralType {
	case CollateralType_COLLATERAL_TYPE_NATIVE:
		if strings.TrimSpace(denom) == "" {
			return fmt.Errorf("collateral_denom is required for native collateral")
		}
		if strings.TrimSpace(contractAddr) != "" {
			return fmt.Errorf("collateral_contract_addr must be empty for native collateral")
		}
	case CollateralType_COLLATERAL_TYPE_PRC20:
		if strings.TrimSpace(contractAddr) == "" {
			return fmt.Errorf("collateral_contract_addr is required for prc20 collateral")
		}
		if _, err := sdk.AccAddressFromBech32(contractAddr); err != nil {
			return fmt.Errorf("invalid collateral_contract_addr: %w", err)
		}
		if strings.TrimSpace(denom) != "" {
			return fmt.Errorf("collateral_denom must be empty for prc20 collateral")
		}
	default:
		return fmt.Errorf("invalid collateral_type")
	}

	return nil
}
