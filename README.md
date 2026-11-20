# Delta Crypto Algo

Enterprise-oriented algorithmic trading platform for Delta Exchange India crypto and RWA-linked derivatives.

> **M1 FOUNDATION — NOT LIVE-TRADING READY.** The backend has no order endpoints and explicitly rejects live mode. The dashboard never sends an order.

## Coverage
A dynamic exchange product catalog will cover **every instrument Delta Exchange India currently exposes**, not a hardcoded BTC/ETH list: crypto perpetual futures, dated futures, call/put options, and real-world-asset (RWA) derivatives such as tokenized gold, silver and stock underlyings. Contract types are capability-discovered per product. A listed symbol is not assumed tradable until product state, status, account entitlement, leverage/margin and risk rules are verified.

## Architecture
- Go modular-monolith execution and control services with hard domain boundaries; split services only on measured demand.
- React + TypeScript responsive operations dashboard.
- PostgreSQL immutable trade/audit ledger and outbox; streaming/cache only as accelerators.
- Exchange public/private WebSockets and REST, with dedicated reconciliation and rate-limit control.
- Oracle Cloud VPS pilot, with a migration design toward redundant production infrastructure.
- Never put trading credentials in GitHub; no unattended live trading activation.

## Status
- Foundation Go HTTP API: `GET /health/live` and `GET /health/ready`.
- Readiness intentionally returns `503` until backend dependencies and safety controls exist.
- Startup rejects live mode and non-loopback listening. `DELTA_LISTEN_ADDR` accepts only `127.0.0.1`, `::1` (bracketed), or `localhost` with an ASCII decimal TCP port in `1..65535`; service names, signs, zero, and out-of-range ports fail configuration validation before binding.
- Domain validation for order intents; no exchange execution implemented.
- Responsive control-panel prototype with clearly labeled unavailable/live-disabled states.
- Automatic tests + gated milestone process.

## Local checks
```bash
go test -race ./...
go vet ./...
go run ./cmd/control-api
# ./apps/control-ui: npm install && npm run build && npm run dev
```

Read [architecture](docs/ARCHITECTURE.md), [milestones](docs/MILESTONES.md), [security](docs/SECURITY.md), [UX standards](docs/UI_UX.md), [performance](docs/PERFORMANCE.md), and [Delta contracts](docs/DELTA_API.md).

## Development branch and agent handbook

**`develop` is the active integration branch.** Read [AGENTS.md](AGENTS.md), [codex/AGENTS.md](codex/AGENTS.md), [product requirements](codex/PRODUCT_REQUIREMENTS.md), [strategy builder specification](codex/STRATEGY_BUILDER_SPEC.md), and [feature status matrix](docs/FEATURE_MATRIX.md) before changing code.

The **Strategies → Configure Strategy** workspace lets users start from a blank/template draft, configure no-code multi-timeframe conditions, risk and option selection, then validate/export JSON. It is an **offline editor**, not a running trading strategy or backtest.

Autonomous agents may merge tested development PRs **into `develop`**. Merge/release into `main`, production deployment and real trades require separately explicit approval.
