# Architecture v0.1 — M1 foundation

## Core decisions
- Modular Go service with explicit data-plane, execution-plane and control-plane boundaries. Preserve ability to split services without prematurely operating microservices.
- React + TypeScript responsive, keyboard-accessible control UI.
- PostgreSQL canonical persistence, append-only audit, transactional outbox, and rebuildable projections. No trade state exclusively in RAM or a cache.
- Account and strategy isolation; no cross-account positions or exposure leaks.
- Versioned configuration, event schemas, exchange metadata and strategies.
- All market coverage is dynamic, **never hardcoded to BTC**. Supports every **Delta-listed** perpetual/dated future, call/put option and RWA-linked derivative. Eligibility varies per product and account.
- Fixed-point price/quantity units, exchange tick/lot/currency metadata and accurate contract-specific P&L, fees, taxes, funding, liquidation buffers.
- A failed acknowledgment is UNKNOWN, not rejected. Reconcile before retry.
- Performance budgets and load tests precede tuning; optimize end-to-end execution rather than UI animations.

## Modules
```
Public exchange WebSocket/REST -> Market data gateway -> normalized events
                                                    -> strategy runtimes
                                                     -> pre-trade risk
                                      durable intent -> OMS / reconciliation
Private exchange WebSocket/REST <-> authenticated exchange adapter
                          OMS -> positions / accounting / audit / alerts
PostgreSQL <-> all stateful modules via transaction boundary / outbox
React dashboard <-> authenticated control API <-> configuration / monitoring
```

### Major bounded contexts
Identity, roles and approvals; accounts/subaccounts; exchange catalog; market data; strategy builder/runtime; backtests and paper simulation; risk policies; order management (OMS); execution gateway; positions, cash and margin; accounting/fees/tax/funding; reporting; observability; alerts; audit and compliance; deployment.

### Catalog / RWA
Catalog sync calls `GET /v2/products` and pages through `meta.after`, storing complete exact specs. Classify by authoritative exchange contract type and available metadata, with user-defined tags for crypto / RWA (asset metadata may not supply a reliable category). Account permissions, legal eligibility, market state and product-specific risks must be checked before trade enablement. Expired/delisted/restricted products may be kept for history, never assumed orderable.

### Order lifecycle
Intent -> deterministic rule evaluation -> per-account risk + kill switches -> reserve capacity -> durable ID -> exchange request -> response or uncertain state -> fill/event reconciliation -> accounting update -> audit. Dedupe by exchange event/order identifiers; never infer position merely from strategy intent. Cancel, cancel-replace, partial fill, late event and restart recovery receive separate testing.

### Reliability
Stale data, inconsistent exchange private feeds, unknown order outcome, database loss or time drift stop new risk. Emergency shutdown has separate semantics from cancel-open and flatten-positions; no guarantee a venue order will execute during an outage. Backups and restore tests mandatory.

### Deployment
An Oracle VPS is an inexpensive pilot, **not high availability**. Separate persistent storage, automatic encrypted backups, monitoring, rollback and strict ingress. Later: redundant execution/DB services and failover tested under realistic faults.

## M1 actual scope
Only documentation, CI, pure domain validation, a loopback-only health-only Go API and an inert dashboard prototype. **No market-data gateway, account connection, risk engine, OMS, WebSocket consumer, trading engine, or live deployment exists yet.**
