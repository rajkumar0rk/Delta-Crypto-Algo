# Engineering quality gates

## Definition of done for each milestone
- Scope and expected behavior documented, UI/ API cross-referenced, no ambiguity in units, timezone, data origin or edge cases.
- Code compiles; automated unit tests and boundary tests pass; lint/static checks pass; integration/contract tests when relevant.
- Security and error-path review; no new high/critical vulnerabilities; credentials never exposed.
- UI accessible and visually aligned at phone/tablet/laptop/desktop sizes; keyboard and responsive table checks.
- Reproducible performance measurements where the change touches hot execution paths; regressions understood.
- Observability (metrics/log context/error alerts) and recovery documentation updated where relevant.
- GitHub PR targets `develop`, green CI, meaningful description/checklist, no unresolved safety blockers. Merge only after verifying all claimed checks.
- `codex/CURRENT_STATUS.md` updated *after* actual result, not prior to tests. Handoff summarizes open issues honestly.

## Tests by layer
- Pure domain: Go unit + race, fixed-point exactness and property-based invariants.
- APIs: input schemas, invalid auth, idempotency, 429/rate-limit, TLS, CSRF, role-specific access, concurrency.
- Strategy DSL: schema, semantic type system, indicator registry, no lookahead, deterministic crossovers, missing-data UNKNOWN behavior.
- Exchange: public/private feed disconnects, order timeout UNKNOWN handling, partial fills, cancel/replace, reconciliation after restart.
- UI: a11y labels, line wrapping, overflow, browser smoke tests, interaction responsiveness, state/error messaging.
- Security: secret detection, dependency pinning/lockfiles, release signing, least-privilege workflows.

## Required commands (initial state)
```bash
go test -race -count=1 ./...
go vet ./...
cd apps/control-ui
npm install --no-audit --no-fund
npm run validate:schema
npm run build
```
Record exact GitHub Actions run link; do not state tests passed before completion. Where lockfile is missing, add a reproducible lockfile and switch to `npm ci` in a separate tracked improvement. Add actual E2E/screenshot testing before calling the UI visually verified.

## Trading non-negotiables
Protective risk gates must be independent of strategy JSON. Production boundaries require explicit human authorization. Under uncertain order or reconciliation status, stop NEW risk and investigate. Never silently rerun possibly accepted orders or hide partial fills.
