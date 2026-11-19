# Configure Strategy — professional product and technical blueprint

## Why this file exists
The user asked for a real **Configure Strategy** experience inside the application, not a sample conversation button or only a schema. Nontechnical users should be guided through a safe workflow; power users can reveal advanced options. Agents are responsible for discovering and implementing all domain-specific details, not repeatedly asking the user which small fields to add.

## Primary route
**Strategies → Strategy Workspace → Configure Strategy → editor**.
User selects **Blank** or a clearly labeled **Illustrative template**; may import previously exported JSON. The editor supports identity, universe, time/session, entry rules, exit rules, execution/position sizing/option selector, risk, review/JSON export. Later backend milestones add safe persistent drafts, validation, versions, backtest, paper activation and controlled live workflow. Never show a fake “Deploy” control.

### Professional UX flow (target)
- Strategy Library: table/card view with status tags from server; search/sort/filter, owner/account, last edit, version, backtest summary, health and actions. **No placeholder trading P&L.**
- Configure: responsive stepper + main form + pinned preview/errors on large screens. Mobile: horizontal step selector, contextual drawers instead of 3 cramped columns, sticky Save Draft actions only when available.
- Progressive disclosure: Basic tab shows safe essentials; Advanced tab shows full order/margin/expiry/session/rule semantics. Advanced values are visible in review before publishing.
- Every value displays units, basis, constraints and help text. Examples: “Stop loss 15% of OPTION ENTRY PREMIUM” vs “200 BTCUSD price points” vs “risk 0.5% of account equity.”
- Rule builder: nested AND/OR groups, long and short split, explicit `previous` offsets and each indicator timeframe, searchable registry, output selection and parameter descriptions. Rule tree can disable a branch without deleting it (future).
- Validation: immediate cheap client validation, authoritative schema/semantic server validation, dry-run historical validation, venue capability and account checks. Show errors adjacent to fields, not only in JSON.
- Review: signal summary, risk summary, effective contracts, estimated max order/position, differences from previous revision, unvalidated data indicators, precise state/readiness.
- Save: versioned draft persisted after auth; optimistic concurrency/ETag, encrypted audit metadata, idempotent save; user controls how active positions react to configuration changes.
- Launch flow later: Backtest → inspect fills/edge cases/monthly results → paper/dry-run → permissions/margin review → operator-reviewed live approval. No default automatic promotion.

## DSL contract
Use `schemas/strategy.v1.schema.json` for CURRENT offline draft syntax; do not assume current schema is comprehensive or final. Migrations must preserve existing versions and past backtest reproducibility.

Typed AST nodes:
```text
SignalTree = Group(ALL|ANY, children) | Comparison(left_operand,operator,right_operand)
Operand    = Constant(decimal) | MarketSeries(field,product,timeframe,offset) | Indicator(id,params,output,product,timeframe,offset)
```
Future extensions: NOT, BETWEEN, arithmetic expressions, set membership, sequences, candle patterns, time filters, order book signals, custom reviewed plugins, multi-leg option spreads. Never allow arbitrary execution, cycles, unbounded nesting, unbounded series or unknown dependencies.

### Critical semantics not yet fully implemented
- Series alignment: given signal evaluation event T, every operand uses latest bar **whose close time ≤ T**, except explicitly authorized tick/intrabar protective feeds. With offset=0, it is the most recent FINALIZED candle, not the forming bar.
- Candle completion and 1m/5m/15m/1h aggregation; missing/incomplete source interval suspends affected signal.
- Indicator initialization warmup, NaN/unknown values, stale data, zero-volume VWAP, division by zero, series gaps; UNKNOWN never silently becomes TRUE.
- Crossovers require previous and current aligned samples. Crossover on partial forming candles must be explicit and separately backtested.
- Long/short conflict, multiple independent strategies, hedging/netting, concurrency, same-instrument signals and account margin reservations.
- Orders apply the exact published strategy revision to risk decisions, even if a user starts editing a later draft.
- Safe idempotent order IDs, cancel/replace and reconciliation after timeout; no duplicate opening order during restart.
- For options, CALL/PUT and BUY/SELL are separate axes with distinct margins and stop direction. Explicit product selection state must be recorded with the signal.
- Reference basis for percent stop, profit and trailing must be explicit, compatible with the product and determined from venue-confirmed cost/entry pricing. Spot, futures and options must never share an assumed “point value.”
- Both SL/TP touched within same bar: need specified conservative execution-order policy or tick replay. Different providers can have different stop trigger semantics.
- Time: user-selected timezone and sessions, including overnight start>end, DST, exchange maintenance, funding, options expiry cutoff and weekend rules.
- Config diff and revision migration: never change open position policy silently; active-position handoff requires explicit lifecycle choices (keep old rules, planned close, migrate with authorization).

## Strategy-level configuration sections with future schema expansion
1. **General:** name, ID, tags, version, author, account scope, notes, status, cloning, strategy template and allocation group.
2. **Market universe:** exchange/env, product type, underlying filters, whitelist/blacklist, currency, selection order, liquidity, spread, min volume, OI, min/max listing age, contract availability, max concurrent subscriptions.
3. **Timing:** entry signal timeframe, trend filter timeframe, independent protective monitor, candle-close vs next-open execution, session windows, time zone, day/week/month filters, holidays, first/last trading minutes, daily flatten policy.
4. **Signals:** long/short enabled, setup/confirmation, candle patterns (including doji thresholds), indicator expressions, breakout/pullback, strength/volume/volatility, multi-symbol correlation, stale-signal age, retest policy, one-entry-per-bar/trend.
5. **Options selector:** CALL/PUT action, buy/sell, expiry and DTE, strike ATM/ITM/OTM/delta/absolute, stable-strike selection policy, IV/Greeks filter, spread/quote depth, premium budget, lot rules, margin obligations, optional multi-leg.
6. **Futures sizing:** lot or notional, leverage/margin mode, stop-risk sizing, account allocation, max exposure, funding tolerance, tick/step rounding and contract multiplier.
7. **Order execution:** market/limit/stop variations as venue supported, maker/post-only, TIF, retries and uncertainty rules, slippage, price offset, delay, stale quote policy, cancel/reprice, reduce-only and partial-fill behavior.
8. **Exits:** fixed ticks/points/percent premium/ATR, signal-based, break-even, activation/trail, step-up profit lock, TP ladder, partial close, time/expiry/session/funding exit, gap policy, forced flatten and restart behavior.
9. **Risk:** max risk per trade/day, concurrent positions, strategy-account-portfolio caps, correlation, losses streak, concentration, liquidation buffers, margins, max orders/trades, circuit breakers, pause/resume, kill switch.
10. **Review:** compile/semantic/type checks, capability detection, risk budget preview, parameter/version hash, config diff, test status, approval state and dry-run execution trace.

## Interface state model and backend APIs (planned)
```text
GET    /v1/strategies            library, pagination, search, filter
POST   /v1/strategies            create versioned draft (idempotency key)
GET    /v1/strategies/{id}       draft/published revision metadata
PUT    /v1/strategies/{id}/draft conditional update using ETag
POST   /v1/strategies/validate   schema + semantic + capability validation
POST   /v1/strategies/{id}/runs  async backtest job
GET    /v1/strategies/{id}/runs/{runId}
POST   /v1/strategies/{id}/publish  approvals, no live enable
POST   /v1/strategies/{id}/paper   human-approved paper run
```
**No routes above are currently implemented.** They are contracts for later milestones, not promises that the UI will work against a server today.

## Testing and release matrix
- Schema: unknown/missing fields, invalid decimals, nested group limits, unsupported indicators, invalid strike/expiry/risk basis, hazardous option SELL sizing and unsupported order types.
- Determinism: 1H filter with 5M stop, next-candle gap, asynchronous market feeds, missing data, month boundary/DST, one candle competing entry/exit, missing expiry, exchange reconnect.
- UI: all step transitions, template prefill, import/export roundtrip, field-level errors, disabled future controls, keyboard navigation, touch, text overflow, data density, viewport 320/360/768/1024/1440/1920, 200% zoom, light/dark.
- Trading: no real credentials/orders in CI, fail-closed uncertain state, order deduplication, position accounting and negative cases via fake exchange.
- Performance: budget under target Oracle VPS constraints, benchmark p50/p95/p99 and GC/memory with realistic instrument and strategy counts.
