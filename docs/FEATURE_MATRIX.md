# Configure Strategy feature matrix (current vs planned)

**This is the source of truth for what exists.** Use `codex/PRODUCT_REQUIREMENTS.md` for the exhaustive target feature set.

| Feature | Current development state | Next implementation gate |
|---|---|---|
| Strategy workspace + Configure Strategy | Offline React UI and starter templates | Visual/accessibility tests, persisted drafts |
| Eight-step no-code editor | Offline UI implemented | Progressive basic/advanced layout, form-level error messages |
| Dynamic crypto/RWA instrument filtering | JSON fields + UI, no venue data | M2 Delta catalog + actual product eligibility |
| Nested AND/OR indicator conditions | JSON AST + offline editor | Registry/type semantics, NOT/sequence features, runtime parity |
| Per-operand timeframes and offsets | UI and schema fields | M6 deterministic aggregation + no-lookahead tests |
| Long/short entries and exits | Separate groups in UI/schema | M6 state machine + priority / simultaneous trigger handling |
| Sessions, weekdays, timezone | UI fields | Server timezone/DST, overnight and venue-session validation |
| Futures sizing and order preferences | JSON fields | M7 risk + M8 exchange capability checks |
| Options strike, DTE, CALL/PUT, BUY/SELL | Offline options selector | M2 metadata, M7 true margin, M8 venue integration |
| Stop/TP/trailing/risk controls | Editable strategy-level fields | Independent enforced M7 risk and M8 OMS order lifecycle |
| JSON import/export + local schema checks | Implemented | Authoritative backend semantics, revisioning, validation suite |
| Persistent strategy library, edit/history/publish | NOT implemented | M5 backend and auth |
| Backtest / paper run | NOT implemented | M10 deterministic replay and simulator |
| Live account connection / order execution | NOT implemented, blocked | Human-approved M15 after all prior gates |
| Performance at enterprise workload | NOT measured | M13 reproducible benchmarks |
| Full a11y and visual QA | Initial responsive CSS only | Browser tests, screen-reader and device regression |

## Definition of "Configure Strategy complete"
A user can configure and persist every supported strategy field, validate against actual venue capabilities and account context, run reproducible backtests, review risk and execution decisions, paper trade, and inspect auditable outcomes. A compiled frontend or valid schema alone is insufficient.
