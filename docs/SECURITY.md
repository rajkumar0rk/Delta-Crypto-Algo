# Safety and security release requirements

## Threat model
Credential theft; malicious config or strategies; unauthorized sessions; price feed spoofing/staleness; WS sequence gaps; partial fills; uncertain order responses; account cross-talk; mistaken product units; margin shocks; exchange/rate-limit outages; public exposure; insider abuse; database loss; compromised CI/dependency.

## Controls
1. No secrets in repo, logs, backups without encryption, CI variables printed to logs, frontend bundles, examples or test fixtures. Prefer managed vault; otherwise strict-permission OS secret injection. Rotate and IP-restrict keys.
2. Trade-only API key, withdrawals off where possible; never add credentials until explicit human authorization.
3. Private admin UI behind VPN/Tailscale and TLS; allowlist origins, auth, MFA, RBAC, session expiry, CSRF on cookie APIs and strong security headers.
4. Authenticated, auditable activation of each strategy and live trading separately; full version history and two-person approval when feasible.
5. No live orders without active margin, risk, market-status, identity, account and audit checks.
6. Limits per account/product/strategy/portfolio: exposure, lots, outstanding requests, leverage, max loss, drawdown, price bands/spread/depth, stale quote and rate limits. Enforce pre-trade and in-flight.
7. Kill switch distinction: stop-new-risk, cancel-open, close-positions; each needs confirmation and retry/reconciliation semantics.
8. Idempotency, uncertain-order handling and restart recovery invariants; never retry a timed-out new order blindly.
9. Audit append-only events, trustworthy clock, alerts, metrics and signed releases. Regular verified backups and restore drills.
10. Vendor risk: monitor Delta India API/terms, RWA eligibility and jurisdictional compliance changes. RWA derivatives can have distinct legal and availability constraints.
11. Feature branches + PR approvals, protected main, least-privilege CI, code scanning and dependency reviews.
12. The repository owner already delegates routine **development-only** merges to `develop` when independent review, exact-head CI, required rules and low-risk checks pass. This standing delegation is not authority to bypass a review gate, merge into `main`, deploy to production, connect real keys, enable live mode or submit real orders; those higher-risk actions require separate current explicit approval.

## Mandatory release gates
[ ] Fault injection / recovery pass
[ ] Risk and reconciliation tested independently
[ ] Demo/paper soak passed with real operational limits
[ ] Measured throughput, latency and memory tested on target VPS
[ ] Credential, secret, authorization and compliance review
[ ] Rollback and encrypted backup restore tested
[ ] Human confirms risk caps, account, instrument scope and go-live
