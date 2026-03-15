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
- Duplicate `trade_id` is rejected per market.
- Supported matches:
  - `BUY_YES <-> BUY_NO` (mint pair, both buyers pay collateral)
  - `BUY_YES <-> SELL_YES`
  - `BUY_NO <-> SELL_NO`
- Fees are charged during settlement and split by `resolver_fee_share_percent`.

## Resolution and claims
- Resolver can `ResolveMarket` (YES/NO) or `VoidMarket`.
- Resolver can update `resolution_source` when resolving.
- RESOLVED payout: winning shares redeem 1:1 collateral.
- VOIDED refund: `(yes + no) / 2` collateral (integer truncation).
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
