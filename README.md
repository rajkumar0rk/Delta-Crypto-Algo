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
