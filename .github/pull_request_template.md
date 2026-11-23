## Scope and dependency milestone

- Milestone / issue:
- Base branch: `develop` (never automatic `main` integration)
- Previous PRs reviewed to prevent overlapping work:

## Behavior and failure modes

- What changes, including defaults, constraints, units and authorization boundaries:
- Failure injection / invalid-input / edge cases:
- Risk/OMS/account isolation and fixed-point semantics, if relevant:
- UI behavior at keyboard, screen-reader and responsive widths, if relevant:

## Evidence (paste actual commands, output and exact commit CI URL)

- [ ] Go format, vet, race tests and build verified where relevant
- [ ] UI lockfile install, schema validation, TypeScript build and browser checks verified where relevant
- [ ] Security scan/dependency/audit reviewed; no known unaddressed high/critical issue
- [ ] Updated docs, tests, feature matrix and status ledger
- [ ] Exact-head GitHub Actions checks green (link):
- [ ] Independent review/required ownership completed (link):
- [ ] Confirm required `develop` ruleset/checks (do not infer from mergeability):

## Non-negotiable release boundaries

- [ ] No real credentials, trading orders, live-mode activation or deployment introduced
- [ ] No automatic `main` merge/release
- [ ] Unknown order/data state fails closed; risk and execution boundaries remain independent
- [ ] Product-specific contract/margin/fee/expiry assumptions are exchange-evidence-backed, or explicitly unknown and unavailable

## Known blockers / follow-up

List outstanding checks and manual acceptance; do not claim milestone completion from green CI alone.
