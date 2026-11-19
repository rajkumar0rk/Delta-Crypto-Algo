# No-code Strategy Builder — M5 design specification

## Purpose and current status
This is a **design and offline draft editor**, not a strategy execution system. Never equate JSON/schema validation with exchange approval, sufficient margin, successful backtesting or authorization to trade. All orders remain disabled.

## Information architecture
| Workspace step | User controls | Safeguards |
|---|---|---|
| Identity | Name, description, tags, draft revision, clone | Immutable published revision later |
| Universe | Dynamic filter or exact Delta product IDs, product types, symbol/underlying filters, max products | Check product availability and account eligibility server-side |
| Time & sessions | Signal candle, fast protective timeframe, closed-candle confirmation, data age threshold, zone/weekdays/windows | Timezone-aware and explicit overnight windows |
| Entry | Long/short rules as nested ALL/ANY groups, typed left/right operands, crossovers, indicator timeframe & offset | Empty group disables its signal; no implicit 'true' |
| Exit | Long/short discretionary exits; separate mandatory protective stops | Reduce-only; reconcile after uncertain response |
| Execution & sizing | Order type, TIF, slippage cap, sizing by lots/notional/% risk, option DTE/strike/side | Risk sizing cannot run without a valid stop; product metadata authoritative |
| Risk | Stop loss, profit target, trailing activation+distance, notional/position/day/drawdown/trade-rate/streak/margin guards | Global account/portfolio risk rules always win |
| Review | Normalized JSON, validation errors, export/import, revision diff (planned) | Save offline draft only; no live activation |

## Editor workflow
Sidebar navigation > Strategies > Strategy Builder. Two-column desktop editor: step rail + controls + pinned JSON/validation pane. Narrow/mobile: single-column step navigation, collapsed preview and sticky draft actions. Field labels/units always visible; all actions keyboard-friendly; errors describe exact offending field.

1. Pick dynamic instrument universe, including all **eligible** crypto and RWA-linked derivatives, without hardcoding BTC/ETH.
2. Choose signal timeframe (e.g. 1h) separately from execution/protective monitoring timeframe (e.g. 5m).
3. Compose recursive condition groups: ALL, ANY, comparisons including crosses_above/below. Each operand can be price field, indicator output or literal value; reference timeframe is explicit.
4. Write distinct long/short entry/exit groups and conflict policy; no implicit mirror, default empty groups don't signal.
5. Define session time zone and days with clear entry window boundaries.
6. Select execution and sizing; configure product-specific option contract choice (DTE, strike and call/put selection) only for option scope.
7. Configure all protective limits; review the JSON and validation messages, export a draft. Running/backtesting is intentionally absent at design stage.

## Versioned schema rules
- `schemas/strategy.v1.schema.json` is Draft 2020-12. `examples/mtf-futures-draft.json` is a non-live sample.
- The strategy DSL is inert data, **no arbitrary JavaScript/Python expressions or user-supplied code**. Later custom indicators run through a reviewed, sandboxed plugin registry.
- Operators `crosses_above`/`crosses_below` compare previous and current aligned samples using finalized bars; under-specified alignment must reject rather than look ahead.
- Indicator `parameters` are extensible but actual names/outputs/parameter bounds are validated by a server-side indicator registry (not by the general JSON Schema).
- Numeric money/size/tick values use decimal strings; do not convert into floats in accounting/execution.
- `universe.selection=dynamic_filter` can select all eligible listed assets; `product_ids` targets exact listed contracts. Use catalog snapshots to stabilize replay.
- Fractional/derived timeframes must aggregate only complete source intervals. Strategy decisions use exchange-clock-aligned events, never chart-render refresh.
- All config revisions require review, effective time, content hash and robust migration; changes to active positions must have explicit behavior and never retroactively modify past fills.
- Risk objects are strategy-local limits only; deployment-wide account/exchange loss and exposure guards cannot be overridden by this document.
- Risk and exchange execution are not implemented by this design. Server must reject invalid combos, impossible option expiry ranges, unsupported venue features, oversized universes and unsatisfied risk budgets.

## Distinctions to preserve
- Options premium SL (% of premium paid) and futures price-move SL are **different**; show their basis/units in UX.
- Spot tick, futures quantity, option contract multiplier, mark price, fee and GST behavior come from product and account context.
- Same JSON must produce deterministic decisions under identical normalized data; simulator must model spread, fee, funding, latency and fill uncertainty.
- A strategy with 1h entry and 1m stop tracking must not wait for 1h close to enforce protective stop.
- Long and short exit rules are different by design, and an empty group means disabled, not match-all.

## Validation/testing matrix
- Schema validation, semantic validation, indicator capability checks, option selector constraints, session DST/overnight behavior, multi-timeframe alignment/no lookahead.
- UI: keyboard, SR labels, 360/768/1280/1920px, 200% zoom, no clipped risk controls, mobile action bar, import malformed/oversized JSON safely.
- Audit: roundtrip export/import, stable IDs, safe defaults, no leaking credentials, no actual orders.
- Runtime later: order uncertainty, risk halt, partial fills, slow feed, stale quotes and restart at candle boundary.
