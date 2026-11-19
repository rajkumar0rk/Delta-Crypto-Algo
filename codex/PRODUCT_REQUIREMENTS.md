# Enterprise product requirements and capability inventory

**Design intent**, not a claim of already-implemented capabilities. Every capability needs configuration, permissions, UX, backend invariants, tests, telemetry and incident operations. Split deliverables into small milestones; do not prematurely expose controls that cannot function.

## 1. Strategy workspace and management
- Strategy library, search/filter/tag/favorites, templates, account/environment ownership, history, folders, cloning, import/export, archived drafts and version diffs.
- Roles (viewer/researcher/operator/admin), author, approver, immutable published revisions, publish/rollback/activate/pause/suspend/retire, effective-at timestamps, acknowledgement of open positions after config changes.
- User-friendly "Configure Strategy" call to action, predictable back/breadcrumb, draft state saved with conflict resolution, unsaved-change warning, inline validation and contextual help.
- First-run quick-start vs advanced mode; do not overwhelm with every control at once. All modes compile to one versioned DSL.
- Distinguish local draft / server saved / validated / backtested / paper-ready / approved / live. Status is a truthful backend state, never color or wishful thinking.

## 2. Instruments / account
- Delta India catalog pagination, refresh cadence, delist/expired instruments, contract spec snapshots, permissions, restrictions and market status.
- Asset type crypto/RWA/other, linear/inverse (where supported), perpetual/dated future, option right, strike, expiry, tick and lot, contract multiplier and quote/settlement currency.
- Multi-symbol allowlist/blacklist, dynamic screeners, liquidity/volume/spread/OI, funding/expiry bounds, top-N caps, per-account routing.
- Exchange account/environment selection, isolated credential vault, subaccounts if supported; account permissions and balance/margin reservations.
- Product capability matrix: allowed order types, reduce-only, post-only, price bands and fees cannot be guessed.

## 3. Market data and time
- WebSocket subscriptions with lifecycle ownership, heartbeats, sequence numbers, deduplication, snapshot+delta merge, out-of-order events, gaps, resync, stale quotes, reconnect, backpressure and bounds.
- Candles based on exact exchange or UTC aligned intervals (open/high/low/close/volume), partial candle vs confirmed candle, session gaps, timezone and DST, missing trades, wick-touch semantics, indicator warm-up.
- Independent timeframes for instrument selection, trend filter, entry trigger, order placement, stop monitor and exit. Never infer engine timeframe from chart UI.
- High-resolution timing, exchange clock drift, holiday/market sessions, funding/expiry and exchange maintenance windows, admission after session ends.
- Historical vs live data provenance, contract roll and expiry continuity, corporate-action complications for RWA where applicable.

## 4. No-code signal DSL
- Condition tree with typed ALL/ANY/NOT and safe bounded grouping; numeric operators, crosses, between, changes, thresholds, candle patterns, indicator output, volume/order-book features, multi-candle sequences.
- Reusable variables, arithmetic expressions via **restricted AST** (never eval), derived series with units, type checking, strategy parameters and parameter sweeps with guardrails.
- Indicator registry: SMA/EMA/VWAP/RSI/ATR/ADX/MACD/Stoch RSI/Supertrend/Donchian/Bollinger/pivots/volume ratio and custom versioned reviewed plugins; outputs/params/timeframes explicit.
- Long and short entry paths, setup state, confirmation/next-candle entry, re-entry/cooldown, one-per-trend/day/per-bar, priority/conflicts, stale-signal rejection, price crossing and gap rules.
- Multi-timeframe semantics: complete candle aggregation, no lookahead, confirmed vs intrabar signals, previous candles, equality/touch rounding, exact event-time ordering, replay parity.

## 5. Execution intent and order management
- Order routing/account/symbol, market/limit/stop/stop-limit/post-only/reduce-only **as supported by venue**, TIF, maker/taker assumptions, limit-price offsets, price collars, slippage caps.
- Sizing units: contracts/lots, quote notional, percent budget/equity, risk-per-stop, margin cap, minimum order/step, precision conversion, rounding direction and rejection reasons.
- Scheduled entry times, latency tolerance, debounce, signal expiration, max chase/reprice attempts, spread/depth checks and partial fills.
- OMS: client idempotency keys, order state transitions, unknown ack, retries only after reconciliation, cancel/replace races, partial fill, duplicate fill, disconnect/restart recovery, broker-as-source-of-truth matching.
- Full order event audit with state, timings and effective config revision; explainable "why order not placed" diagnostics.

## 6. Options-specific rules
- CALL/PUT; BUY/SELL **explicit separate mappings**, different margin and adverse price direction.
- Select underlying, exact or relative strike, ATM/ITM/OTM/delta, nearest eligible expiry, exact expiry, min/max DTE, strike-first/expiry-second policy if configured, replacement behavior when a contract is missing.
- Option chain/Greeks/IV/OI/spread/liquidity filtering with reliable quote timestamps; exercise/settlement, expiry cutoff, premium, intrinsic/time value, margin obligations for writing options.
- Stop and profit references: entry premium, current option quote, underlying movement, delta-adjusted value, absolute quote, % premium, account equity. SELL option margin is **not** capped by received premium.
- Multi-leg spreads/hedges/combos only as a separately reviewed milestone with atomicity limitations and portfolio-margin awareness.

## 7. Futures-specific rules
- Linear/inverse contract interpretation (if venue supports), point value, contract multiplier, mark/index/last source, cross/isolated margin, leverage, maintenance buffer, funding and liquidation exposure.
- Exact futures SL/TP in ticks, points, percent entry, ATR or reference candle, price bands, exchange triggers and rounding.

## 8. Exit/risk behaviors
- Independent stop, take-profit, trailing activation+distance, breakeven stop, profit-lock steps, multi-target scaling, partial reductions, time-based exit, volatility/funding/expiry exit, signal flip and portfolio-level liquidation policies.
- Price touch vs confirmed close vs event tick semantics, intrabar vs bar-based stop checks, stop priority when SL+TP touched in same OHLC bar, conservative backtest fill when chronology unknowable.
- Per position/strategy/account/portfolio max loss and exposure, max orders/trades, max daily drawdown, leverage, open order caps, notional/lot/budget, correlation, concentration, consecutive losses and cooldown.
- Hard risk gateway independently enforces global limits, reserve/free margin and venue status; no strategy JSON can weaken hard limits. Fail-closed on uncertainty.
- Emergency halt, cancel open, flatten portfolio are **different actions** requiring confirmations, idempotency, clear status and execution reporting.

## 9. Backtest / research / paper
- Same versioned signal/risk domain logic across backtest and real-time where possible; model venue-specific fill/slippage/fees/GST/funding/margin and contract delistings.
- Reproducible dataset manifest, candle integrity, warmup, missing data, lookahead/multi-timeframe, latency, same-bar SL/TP collision policy.
- Robustness: in/out sample, walk-forward, bootstrap/Monte Carlo, parameter overfit detection, capacity testing, realistic capital.
- Metrics: monthly PF, gross/net P&L, equity curve, drawdown, total return, win rate, exposure, avg win/loss, fees, slippage, latency, tail risk. Never label hypothetical profit as realized.
- Demo/paper trade lifecycle, provider differences, deployment dry-run, parity diagnostics, simulation failure injection.

## 10. UI/UX
- Strategy library and guided builder; different tabs on desktop vs mobile, progress indicator, field units, validation next to control, searchable indicators/conditions, condition preview, enable/disable switches.
- Type-driven controls so an option premium field never masquerades as futures points; contextual tooltips/examples, safe defaults, descriptive error messages, read-only JSON preview and advanced JSON import/export.
- 8px grid/4px micro spacing; 44px touch targets; accessible labels, keyboard focus, role semantics, right-aligned numeric columns, contrast, empty/loading/error/stale/partial states.
- Charts, reporting, order/position panels, alerts, optimization reports, reconciliation warnings; do not fabricate balances, prices, latency, P&L or live connection.
- Virtualized high-cardinality tables, filtered subscriptions, throttled UI updates, no render-blocking massive JSON, test 320/360/375/768/1024/1440/1920px, zoom and screen readers.

## 11. Security and operations
- Credential encryption/vault, TLS/VPN, MFA/RBAC, trusted sessions, nonrepudiable audit, segregated environments, CSRF, rate limits, replay/auth checks.
- Structured logs with redaction, PII/secrets scanning, metrics p50/p95/p99, traces, rate-limit budgets, alerting, strategy-health checks, graceful shutdown.
- HA roadmap vs tiny Oracle VPS pilot, backup/restore, DR, disk pressure, clock sync, deployment rollback and release approvals.
- Compatibility and change management for Delta API and exchange rule/fee/margin updates.
- Runbooks for outage, partial state, data loss, margin calls, cutover, emergency stop, liquidations and unknown order status.

## Planning rule
Never call the application "fully configurable", "enterprise complete", or "very fast" without an explicit feature matrix, capability flags, measured load test, and honest gaps. Complete advanced modules incrementally with testable contracts and minimal coupling.
