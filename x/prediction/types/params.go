package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "prediction"
	StoreKey   = ModuleName
	RouterKey  = ModuleName
)

var (
	KeyParams = []byte("prediction_params")
)

const (
	DefaultCreateMarketBond      = "20000000"
	DefaultCreateMarketBondDenom = "uusdt"
	DefaultMarketFeeBps          = uint64(50)
	DefaultResolverFeeSharePct   = uint64(80)
	DefaultMaxOrderLifetimeBh    = uint64(216_000)
	DefaultMaxOpenOrdersPerUser  = uint64(1_000)
	DefaultMaxOpenOrdersPerMkt   = uint64(100)
	DefaultOrderPruneIntervalBh  = uint64(1_200)
	DefaultOrderPruneRetainBh    = uint64(216_000)
	DefaultOrderPruneScanLimit   = uint64(20_000)
	DefaultOrderPruneDeleteLimit = uint64(2_000)
)

// Params defines prediction module parameters.
type Params struct {
	MaxBatchSize            uint64 `json:"max_batch_size" yaml:"max_batch_size"`
	CreateMarketBond        string `json:"create_market_bond" yaml:"create_market_bond"`
	CreateMarketBondDenom   string `json:"create_market_bond_denom" yaml:"create_market_bond_denom"`
	MarketFeeBps            uint64 `json:"market_fee_bps" yaml:"market_fee_bps"`
	ResolverFeeSharePercent uint64 `json:"resolver_fee_share_percent" yaml:"resolver_fee_share_percent"`
	MaxOrderLifetimeBh      uint64 `json:"max_order_lifetime_bh" yaml:"max_order_lifetime_bh"`
	MaxOpenOrdersPerUser    uint64 `json:"max_open_orders_per_user" yaml:"max_open_orders_per_user"`
	MaxOpenOrdersPerMarket  uint64 `json:"max_open_orders_per_market" yaml:"max_open_orders_per_market"`
	OrderPruneIntervalBh    uint64 `json:"order_prune_interval_bh" yaml:"order_prune_interval_bh"`
	OrderPruneRetainBh      uint64 `json:"order_prune_retain_bh" yaml:"order_prune_retain_bh"`
	OrderPruneScanLimit     uint64 `json:"order_prune_scan_limit" yaml:"order_prune_scan_limit"`
	OrderPruneDeleteLimit   uint64 `json:"order_prune_delete_limit" yaml:"order_prune_delete_limit"`
}

type GenesisState struct {
	Params       Params      `json:"params" yaml:"params"`
	Markets      []*Market   `json:"markets" yaml:"markets"`
	Orders       []*Order    `json:"orders" yaml:"orders"`
	Positions    []*Position `json:"positions" yaml:"positions"`
	NextMarketID uint64      `json:"next_market_id" yaml:"next_market_id"`
	NextOrderID  uint64      `json:"next_order_id" yaml:"next_order_id"`
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:       DefaultParams(),
		Markets:      []*Market{},
		Orders:       []*Order{},
		Positions:    []*Position{},
		NextMarketID: 1,
		NextOrderID:  1,
	}
}

func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	if gs.NextMarketID == 0 {
		return fmt.Errorf("next_market_id must be >= 1")
	}
	if gs.NextOrderID == 0 {
		return fmt.Errorf("next_order_id must be >= 1")
	}

	for _, m := range gs.Markets {
		if m == nil {
			return fmt.Errorf("market cannot be nil")
		}
		if m.Id == 0 {
			return fmt.Errorf("market id must be > 0")
		}
	}

	for _, p := range gs.Positions {
		if p == nil {
			return fmt.Errorf("position cannot be nil")
		}
		if p.MarketId == 0 {
			return fmt.Errorf("position market_id must be > 0")
		}
		if _, err := sdk.AccAddressFromBech32(p.Address); err != nil {
			return fmt.Errorf("invalid position address %q: %w", p.Address, err)
		}
	}

	for _, o := range gs.Orders {
		if o == nil {
			return fmt.Errorf("order cannot be nil")
		}
		if o.Id == 0 {
			return fmt.Errorf("order id must be > 0")
		}
		if o.MarketId == 0 {
			return fmt.Errorf("order market_id must be > 0")
		}
		if _, err := sdk.AccAddressFromBech32(o.Trader); err != nil {
			return fmt.Errorf("invalid order trader %q: %w", o.Trader, err)
		}
	}

	return nil
}

func DefaultParams() Params {
	return Params{
		MaxBatchSize:            500,
		CreateMarketBond:        DefaultCreateMarketBond,
		CreateMarketBondDenom:   DefaultCreateMarketBondDenom,
		MarketFeeBps:            DefaultMarketFeeBps,
		ResolverFeeSharePercent: DefaultResolverFeeSharePct,
		MaxOrderLifetimeBh:      DefaultMaxOrderLifetimeBh,
		MaxOpenOrdersPerUser:    DefaultMaxOpenOrdersPerUser,
		MaxOpenOrdersPerMarket:  DefaultMaxOpenOrdersPerMkt,
		OrderPruneIntervalBh:    DefaultOrderPruneIntervalBh,
		OrderPruneRetainBh:      DefaultOrderPruneRetainBh,
		OrderPruneScanLimit:     DefaultOrderPruneScanLimit,
		OrderPruneDeleteLimit:   DefaultOrderPruneDeleteLimit,
	}
}

func (p Params) Validate() error {
	if p.MaxBatchSize == 0 {
		return fmt.Errorf("max_batch_size must be > 0")
	}

	if p.CreateMarketBondDenom == "" {
		return fmt.Errorf("create_market_bond_denom cannot be empty")
	}
	bondAmount, ok := sdkmath.NewIntFromString(p.CreateMarketBond)
	if !ok || !bondAmount.IsPositive() {
		return fmt.Errorf("create_market_bond must be a positive integer")
	}

	if p.MarketFeeBps > BPSDenominator {
		return fmt.Errorf("market_fee_bps cannot exceed %d", BPSDenominator)
	}
	if p.ResolverFeeSharePercent > 100 {
		return fmt.Errorf("resolver_fee_share_percent cannot exceed 100")
	}
	if p.MaxOrderLifetimeBh == 0 {
		return fmt.Errorf("max_order_lifetime_bh must be > 0")
	}
	if p.MaxOpenOrdersPerUser == 0 {
		return fmt.Errorf("max_open_orders_per_user must be > 0")
	}
	if p.MaxOpenOrdersPerMarket == 0 {
		return fmt.Errorf("max_open_orders_per_market must be > 0")
	}
	if p.MaxOpenOrdersPerMarket > p.MaxOpenOrdersPerUser {
		return fmt.Errorf("max_open_orders_per_market cannot exceed max_open_orders_per_user")
	}
	if p.OrderPruneIntervalBh == 0 {
		return fmt.Errorf("order_prune_interval_bh must be > 0")
	}
	if p.OrderPruneRetainBh == 0 {
		return fmt.Errorf("order_prune_retain_bh must be > 0")
	}
	if p.OrderPruneScanLimit == 0 {
		return fmt.Errorf("order_prune_scan_limit must be > 0")
	}
	if p.OrderPruneDeleteLimit == 0 {
		return fmt.Errorf("order_prune_delete_limit must be > 0")
	}
	if p.OrderPruneDeleteLimit > p.OrderPruneScanLimit {
		return fmt.Errorf("order_prune_delete_limit cannot exceed order_prune_scan_limit")
	}

	return nil
}
