# Project status — verified integration

**Date:** 2026-11-19 (IST), M2 catalog domain merge. This file is a handoff snapshot; always check GitHub `develop`, new PRs and Actions before acting.

## Verified repository state
- **`develop`** is the unified development integration branch. After PR #20 merge, HEAD is `44df6ed85058fca8af5095b0eb24418e23d4cd94`.
- M1 foundation, initial M5 offline Strategy Builder, and the real Configure Strategy workspace/agent handbook are reachable from `develop`.
- **PR #18 merged**: https://github.com/rajkumar0rk/Delta-Crypto-Algo/pull/18
- GitHub Actions on `develop` HEAD **passed** (Go tests/race/static/build; UI typecheck/build; strategy JSON schema tests): https://github.com/rajkumar0rk/Delta-Crypto-Algo/actions/runs/44479713873
- Former PR #1 and #17 were closed as superseded; do not attempt to merge their already integrated commits again.
- User authorized autonomous **development-only merges to `develop`** after green CI and safety review. No automatic merging to `main`.

## M2 catalog domain integration (verified)
- PR #20 merged into `develop`: https://github.com/rajkumar0rk/Delta-Crypto-Algo/pull/20 (merge commit `44df6ed85058fca8af5095b0eb24418e23d4cd94`).
- Exact PR head `10e2ec4fcc51d5fcfa306e733ca24a5f4b052cc5` passed Go formatting, vet, race tests, build and UI/schema CI: https://github.com/rajkumar0rk/Delta-Crypto-Algo/actions/runs/44506163654
- `internal/domain/catalog` now validates typed futures/options and crypto/RWA-capable metadata with exact decimal strings, expiry, duplicate IDs/symbols, complete snapshots, and fail-closed freshness/product-state checks.
- This is only the domain model. **No actual Delta Exchange India REST catalog, pagination, account eligibility, exchange margin, or live order routing exists.** `Snapshot.Tradeable` is a metadata precondition, not permission to trade.
- Next M2 slice: public read-only Delta catalog fetch with bounded pagination, mapping, fixture contract tests, strict rejection of incomplete/stale data and unknown product capabilities.

## Current Configure Strategy experience
- Navigation **Strategies** opens the strategy workspace, with visible **Configure Strategy** button.
- A blank editor and four example templates (blank/MTF futures/options oscillator/price breakout) are available and open the same editor prefilled.
- Editor has eight steps: identity, universe, timeframes/session, long/short entry, long/short exit, execution/sizing, risk/protection, review.
- Uses schema v1, typed nested AND/OR rule groups and per-operand timeframes/offsets.
- JSON import/export, Ajv structural and basic semantic validation; futures and options examples.
- UI builds and checks pass. **No live strategies, backtests, account state, server persistence or exchange orders** exist.
- Responsive styles are written but browser visual comparison/a11y/E2E acceptance has **NOT** yet been performed.

## Critical missing work, in dependency order
1. **M2 product catalog**: accurate Delta India product metadata, expiry/strike/tick/lot/margin and dynamic crypto/RWA eligibility, including stale and delisted states.
2. **M5 authoring hardening**: normalized indicator registry, server semantic validation, no-lookahead rules, actual persisted versioned drafts, user roles, safer input sizing, full examples/validation fixtures, reliable edit/resume and browser a11y tests.
3. M3 public/private feed correctness, M4 auth/secrets, M6 deterministic runtime, M7 independent risk gate, M8 durable OMS, M9 P&L/fees, M10 backtest/paper parity.
4. M11 professional UI completion, M12 observability, M13 performance/security/regression, M14 Oracle pilot/demo soak, M15 separately approved live readiness.

## Next agent steps
- Read root `AGENTS.md`, `codex/AGENTS.md`, `codex/PRODUCT_REQUIREMENTS.md`, `codex/STRATEGY_BUILDER_SPEC.md`, `docs/FEATURE_MATRIX.md`, `docs/MILESTONES.md`.
- Inspect latest `develop` HEAD and CI to ensure this status is not stale.
- Implement the smallest unblocked part of M2 or M5 with actual tests, target a PR to `develop`, and merge only after gates pass.
- Maintain strict live-trading lock and no real credentials/deployments. Never claim benchmark or production readiness without measurements.

## M2 public REST pagination follow-up (2026-11-19)
- PR #22: https://github.com/rajkumar0rk/Delta-Crypto-Algo/pull/22 adds a read-only HTTPS-only `GET /v2/products` transport, bounded pagination, strict envelope/error handling, TLS fixture tests and design notes. CI must be checked on the exact head before merge; no live Delta fixture verified yet.
- Transport deliberately returns raw JSON only. Contract mapping, authoritative classification, freshness/atomic cache, capability rules, full venue fixture verification and UI integration are NOT done. No credentials, orders or VPS deployment.
- Next step after green CI: verified Delta India product response mapping with strict fail-closed unknown fields and typed catalog validation.

## M2 public catalog transport — verified merge (2026-11-19)
- PR #22 merged into `develop`: https://github.com/rajkumar0rk/Delta-Crypto-Algo/pull/22 (merge commit `880a8895538fde935b2e25677c4e3fef410e26a5`).
- Exact PR head `38c0517aad8d952a9fa1f5427806a13b0f6c8be6` passed Go formatting, tests, race/vet/build and UI/schema checks: https://github.com/rajkumar0rk/Delta-Crypto-Algo/actions/runs/44534347476
- New `internal/adapters/delta.PublicCatalog.Fetch` is read-only HTTPS transport with bounded pagination and failure-path tests. Raw product records are **not** validated for trading or exposed as eligible contracts.
- Next: verify authoritative Delta India product response fixtures, field-by-field mapping into `internal/domain/catalog`, venue capability matrix, atomic cache/freshness checks and contract tests. Review redirects/origin constraints before allowing configurable hosts.
- No credentials, live orders, deployment or changes to `main`.

## M2 redirect safety hardening (2026-11-19, pending PR validation)
- Feature branch `feature/m2-reject-catalog-redirects` branched from verified `develop` HEAD `2746cfa4153fc7e9647422162ed444d96fa0ddae`.
- Catalog transport now rejects HTTP 3xx without following redirects, including cross-origin redirects. It copies the configured HTTP client instead of changing caller policy. TLS fixture test asserts no redirected request reaches the target and partial results remain empty.
- **Not yet verified by CI at this handoff.** Do not merge without checking exact PR head Go/race/vet/build and UI/schema CI.
- M2 remains open: authoritative Delta India product fixture mapping, product capability/margin/fee accuracy, freshness/cache, UI integration, contract acceptance tests.
- No credentials, orders, deployment or `main` changes.

## M2 page-size hardening (2026-11-19; CI pending)
- Branch `feature/m2-enforce-catalog-page-byte-limit` from `develop` commit `f21a3751ea974353eeeefef10029d890b948499d`.
- Public catalog now measures the entire HTTP response body before JSON decoding and rejects bodies exceeding 4 MiB, including valid JSON followed by excess whitespace. Incomplete fetches return no partial products.
- Added regression fixture for valid JSON plus excessive whitespace. CI must verify Go formatting, tests/race/vet/build and UI/schema before merge; no local test execution claimed.
- M2 still needs authoritative venue field mapping, margin/fee/capability verification, atomic refresh/freshness and UI integration. No private API, credentials, live orders or deployment.

## M1 API response isolation hardening (2026-11-20; PR pending)
- Branch `feature/m1-api-content-security-policy` was created from verified `develop` HEAD `60fb41ed034ad86fe401623472d3d81b67d116bf`.
- Added API-only Content-Security-Policy (`default-src 'none'; frame-ancestors 'none'; base-uri 'none'; form-action 'none'`) and Cross-Origin-Resource-Policy (`same-origin`) on all foundation API responses, including missing paths and rejected methods. No browser UI is served by this handler.
- Added table-driven Go regression tests for liveness, fail-closed readiness, missing paths, and disallowed methods. **GitHub CI and independent review remain required before merging; do not mark M1 complete.**
- Existing M1 PRs #35–#39 and M2 PRs #29/#30 require safe integration/review. M1 still lacks branch protection, npm lockfile/reproducibility, browser a11y and security acceptance. No trading, credentials, VPS or `main` changes.

## M1 PR #40 CI formatting follow-up (2026-11-20)
- Exact PR head `54f95212bd3bf44471461b29dcd64026f9940dad` failed the Go formatting step; Go vet/race/build were not run, while the UI job passed: https://github.com/rajkumar0rk/Delta-Crypto-Algo/actions/runs/44728377585
- Corrected the table-field alignment in `internal/httpapi/server_security_test.go` with commit `9b0b4190c7b3b6556e73812f7b9d3b9562225372`; this is a formatting-only correction, not a new feature. CI on the corrected head and independent review must pass before merge.
- M1 remains incomplete. Do not merge this draft PR or treat the CSP as verified until exact-head CI succeeds; prioritize consolidation/review of existing M1 PRs #35–#40 before starting more overlapping changes.

## M1 API header CI remediation (2026-11-20, verification pending)
- PR #40 https://github.com/rajkumar0rk/Delta-Crypto-Algo/pull/40 was still failing Go formatting at exact head `afc66ac33b3fff912de3e2af502076a9f5870457`: https://github.com/rajkumar0rk/Delta-Crypto-Algo/actions/runs/44736740945. UI passed, Go vet/race/build were skipped. `internal/httpapi/server_security_test.go` was named in the formatting output.
- Corrective commit `37e56806e3d1618fd07be018be84e0bc68c05f51` normalizes the test's struct field declarations; local `gofmt -d` of the corrected test showed no changes. **GitHub CI for this exact head has not yet been observed**; PR stays draft/unmerged pending full green CI and independent review.
- `develop` remains `60fb41ed034ad86fe401623472d3d81b67d116bf` at this inspection. M1 still lacks independent reviews, branch protection, deterministic npm installs, browser a11y/visual tests and security acceptance. M2 catalog work must follow M1 exit gates; no M5 expansion.
- No live credentials, trading, deployments, main changes or changes to other schedules.

## M1 API security policy formatting correction — 2026-11-20
- PR #40 remains draft and unmerged: https://github.com/rajkumar0rk/Delta-Crypto-Algo/pull/40
- Previous exact-head CI `44745102756` **failed Go formatting** in `internal/httpapi/server_security_test.go`; UI passed, Go vet/race/build were skipped: https://github.com/rajkumar0rk/Delta-Crypto-Algo/actions/runs/44745102756
- Replaced the header map alignment with local `gofmt`-verified spacing in commit `17e8028c9cea4d30a1342a20ef863a37c2d26174`. CI for this commit was not yet available at correction time. Do not merge until exact latest-head CI is green and independent review is complete.
- `develop` remains `60fb41ed034ad86fe401623472d3d81b67d116bf` at the start of this run; M1 incomplete. Prioritize M1 acceptance and PR consolidation before M2, M3 or M5 expansion.
