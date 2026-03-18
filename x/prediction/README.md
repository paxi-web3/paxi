# x/prediction MVP Design (Non-Escrow)

This module implements a binary YES/NO prediction market with on-chain market/order/position accounting.

## Scope
- On-chain: market state, order state, positions, settlement application, resolution, claims.
- Off-chain: orderbook, matching engine, relayer logic.

## Core model
- Users place orders from wallet via `MsgPlaceOrder` without pre-deposit.
- Placing/canceling orders does not transfer funds.
- `MsgSplitPosition` and `MsgMergePosition` are available and do not charge settlement fee.
- Funds and share ownership are enforced only at settlement (`MsgApplyTradeBatch`).
- YES/NO are internal state entries (`Position.yes_shares`/`no_shares`), not bank tokens.

## Settlement model
- Market resolver submits `MsgApplyTradeBatch` with matched trades (`trade_id`, `order_a_id`, `order_b_id`).
- `trade_id` already applied in the same market is skipped (`TradeSkipped: duplicate_trade_id`).
- Supported matches:
  - `BUY_YES <-> BUY_NO` (mint pair, both buyers pay collateral)
  - `BUY_YES <-> SELL_YES`
  - `BUY_NO <-> SELL_NO`
- `BUY_YES <-> BUY_NO` does not allow `MARKET <-> MARKET`; at least one side must be `LIMIT`.
- Every trade must provide both `yes_execution_price` and `no_execution_price`.
- For `BUY_YES <-> BUY_NO`, `yes_execution_price + no_execution_price` must equal `1_000_000`.
- Settlement uses side-specific notional:
  - YES side notional = `match_amount * yes_execution_price`
  - NO side notional = `match_amount * no_execution_price`
- Fee per matched leg is computed from notional and split by `resolver_fee_share_percent` (remainder to market creator).
- For `BUY_YES <-> BUY_NO`, each buyer pays `trade_notional + fee` (fee charged on top, collateral remains fully backed in module).
- For `BUY_<OUTCOME> <-> SELL_<OUTCOME>`, buyer pays `trade_notional`, seller receives `trade_notional - fee`.
- `market.total_trade_volume` is accumulated in collateral notional (micro units, e.g. `uusdt`), not share count.
- `last_yes_trade_price` / `last_no_trade_price` are updated from batch-level amount-weighted canonical price and always satisfy `last_yes_trade_price + last_no_trade_price = 1_000_000`.

## Resolution and claims
- Resolver can `ResolveMarket` (YES/NO) or `VoidMarket`.
- Resolver can set/update `resolution_source` when resolving (if non-empty).
- RESOLVED payout: `winning_shares * CollateralUnit` (1 winning share = 1.0 collateral unit, e.g. `1_000_000` micro).
- VOIDED refund: `(yes_shares + no_shares) * CollateralUnit / 2` (integer truncation).
- Claims are one-time via `claimed_payout` / `claimed_void_refund`.

## Status flow
`OPEN -> CLOSED -> RESOLVED` or `OPEN/CLOSED -> VOIDED`

Markets auto-close at `close_time` during `BeginBlock`.

## Gas benchmark
- Run:
```bash
GOTOOLCHAIN=go1.25.5 go test ./x/prediction/keeper -run '^$' -bench 'BenchmarkApplyTradeBatch500(Native|PRC20)$' -benchtime=1x -count=1
```
- Reported metrics include:
  - `gas/batch`, `gas/trade`
  - `paxi/batch@0.05`, `paxi/trade@0.05` (assuming `min-gas-prices=0.05upaxi`)
- `BenchmarkApplyTradeBatch500PRC20` uses a mock PRC20 execute keeper with fixed gas per execute to give a controllable estimate.
