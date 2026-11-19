# Enterprise milestone plan

All work is sequential and test-gated. Each run: inspect current repository and latest PR/CI -> select next unblocked milestone or smallest issue -> implement on a feature branch -> test -> publish PR/evidence -> record remaining blockers. No credentials, trading or production changes without explicit human authorization.

| ID | Milestone | Exit criteria |
|---|---|---|
| M1 | Foundation / architecture | Repo, Go baseline, UI prototype, CI, security and market coverage architecture; PR reviewed |
| M2 | Universal product catalog | Delta India REST pagination, all crypto/RWA instruments, expiry/strike/tick/lot, stale and delisting behavior, contract tests |
| M3 | Market data gateway | Public WS channels, subscriptions, reconnect, gap detection, bounded queues, metrics and replay tests |
| M4 | Accounts + identity | Secure login/RBAC/MFA, account mapping, secure secrets, environment separation, private API tests |
| M5 | Strategy authoring | Fully versioned schema, multi-symbol/timeframe DSL, validation, dry-run and configuration editor |
| M6 | Deterministic signal engine | Efficient candle building, indicator/plugin interfaces, concurrency bounds, accurate timestamps |
| M7 | Risk engine | Account+portfolio+strategy limits, margin/expiry/volatility rules, stale data and kill switches, adversarial tests |
| M8 | Order management + exchange | Private WS, rate limits, durable OMS, idempotency, uncertain requests, fills, cancellation, reconciliation |
| M9 | P&L and positions | Multi-asset accounting, fees/GST/funding/option premiums, tax outputs, margin and liquidation awareness |
| M10 | Paper / backtest parity | Exchange simulator, cost/slippage models, regression fixtures, restart and event replay parity |
| M11 | Enterprise dashboard | Responsive trading cockpit, order/position/strategy panels, accessible UI design system, real-time state |
| M12 | Observability + admin | Tracing, latency histograms, alerts, audit, incident workflows and operators runbooks |
| M13 | Performance and security | Load/soak tests, p95/p99 latency budgets, vulnerability checks, chaos/failure drills, restore test |
| M14 | Oracle deployment + demo soak | Infrastructure as code, isolated demo credentials, backup, rollback, recovery drills, complete operating manual |
| M15 | Live-readiness review | Human approvals; no live order until safeguards independently tested and limits explicitly configured |

## Acceptance philosophy
- PR + green CI is necessary, not sufficient: require functional evidence for each feature.
- Prioritize safety/recovery over throughput optimization and decorative polish.
- No absolute zero-bug assurance. Maintain zero **known critical severity** defects as a release gate; retain operational monitoring.
- Hard gates apply even if milestone tasks run unattended.
- Every automated run must update GitHub issues/PRs, not just report vague progress.

## Initial performance budgets (targets, not measured facts)
- UI stable alignment at 360px, 768px, 1280px, 1920px with no overflow or clipped controls.
- Dashboard interactions target <100ms input responsiveness under typical load.
- Virtualized high-density tables and selective subscriptions, zero huge all-symbol render.
- In-process intent validation target p99 <1ms in a controlled benchmark (excludes network and disk).
- Track end-to-end tick-to-decision, risk-to-wire and exchange acknowledgment separately; report p50/p95/p99, queue lag, memory, GC and dropped messages.
- Gate any scaling claim on reproducible load tests with specified Oracle VPS hardware and market workload.
