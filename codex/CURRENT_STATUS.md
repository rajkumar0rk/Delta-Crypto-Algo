# Project status — M1 consolidated, acceptance pending

## 2026-11-21 — M1 author-side review and expanded accessibility verification

- PR #46 latest audited code SHA `6747cd825a1645f13beaa0ef7d499fa9546df081` passed both CI jobs: https://github.com/rajkumar0rk/Delta-Crypto-Algo/actions/runs/44905852868. Checks: Go format, module integrity, vet, race suite and build; UI pinned dependencies, npm ci, audit, schema, typecheck, browser viewport and keyboard tests.
- Browser tests now check 320px narrow reflow as a **zoom-equivalent approximation** (not a direct manual 200% browser zoom test), 360/768/1280/1920 widths, reduced-motion preference, accessibility-tree button names, skip-link keyboard focus, mobile drawer Enter open / Shift+Tab containment / Escape close / focus restoration. Six screenshots collected. This **does not constitute manual screen-reader acceptance**.
- Author-side code review found and fixed a broken literal README newline, hidden mobile navigation remaining focusable, missing mobile drawer keyboard/focus affordances, and absent hardening headers on an unlikely RNG-error response path. This **is not an independent human approval**.
- Branch governance: repository rulesets GET returned `[]`; classic protection GET returned HTTP 403 to this GitHub integration, so `develop` protection cannot be confirmed or configured in this connection.
- No independent PR reviewer is assigned or has approved. Therefore M1 **remains open**, PR #46 unmerged, `develop` unchanged; require independent reviewer, repository owner protection setup, and manual screen-reader/zoom acceptance. Do not bypass these gates or touch `main`/live.

**2026-11-21 (IST) verified handoff:** `develop` remains `60fb41ed034ad86fe401623472d3d81b67d116bf`; this document on feature branch is not yet integrated. `main` has not been touched.

## M1 integrated work proposed for develop
- Consolidated foundation PR **#46**: https://github.com/rajkumar0rk/Delta-Crypto-Algo/pull/46 (the **sole** M1 integration candidate; PRs #35–#42 and #45 overlap/superseded).
- Verified green full CI for code head `39c028d80249c71b52c25a6279ffe029becaca67`: https://github.com/rajkumar0rk/Delta-Crypto-Algo/actions/runs/44904413027 — Go module verify/tidy, formatting, vet, race tests, build; UI exact-version guard, lockfile npm ci, high-severity production dependency audit, JSON schema/TS build, Chrome viewport checks and screenshot artifact.
- Browser smoke successfully checked 360, 768, 1280, 1920 widths: safety messaging, skip link, visible button labels, no document horizontal overflow; also Strategies navigation and Configure Strategy CTA. The screenshots were downloaded and manually inspected for obvious alignment/clipping. This is **not** full screen-reader, 200% zoom or keyboard-only acceptance.
- Backend remains loopback-only, paper/demo-only, readiness 503, no trading routes or exchange credentials. API adds strict browser headers and random request IDs; no live trading enablement.
- Additional 2026-11-21 M1 CI hardening: checkout, setup-go, setup-node and artifact upload are pinned to verified upstream action commit SHAs. The preceding green full CI on updated browser runner head `c8c0551fb9aee7a53feec80b1b18227ae22bbb21` is https://github.com/rajkumar0rk/Delta-Crypto-Algo/actions/runs/44904755562. Immutable-action pin update passed exact-head Go+UI CI at `1f041d407dff179886a680a708001ffa69f2c246`: https://github.com/rajkumar0rk/Delta-Crypto-Algo/actions/runs/44904889053.
- M1 acceptance ledger: [docs/M1_ACCEPTANCE.md](../docs/M1_ACCEPTANCE.md).
- **Not complete:** PR #46 still needs an independent human security/code review and the owner to configure/verify `develop` branch-protection/rulesets with required reviews and checks. Github API read of protection returned HTTP 403; no proof it is configured. Combined code CI is green; final docs-only head CI still needs checking before merge.
- M2 next after M1 acceptance: reconcile #29/#30 and verify authoritative Delta India crypto/futures/options/RWA catalog mapping, margin/fees/expiry/capabilities; M5 remains deferred.
- Do not push/merge main, deploy VPS, add production keys, or send any real orders.

## Historical handoff (2026-11-19)

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

## 2026-11-21 — consolidated M1 acceptance batch (PR #46, not merged)
- Integration baseline: develop at 60fb41ed034ad86fe401623472d3d81b67d116bf. PR: https://github.com/rajkumar0rk/Delta-Crypto-Algo/pull/46
- Combined overlapping M1 proposals #35–#42/#45 into one feature branch: strict canonical loopback TCP port parsing, deny-trading-route regression tests, server-generated request IDs, CSP/CORP API headers, Go module integrity/tidy checks, Go toolchain from go.mod, no persisted checkout credentials, pinned UI dependency guard, committed npm lockfile and npm ci.
- At handoff exact-head CI results were not yet returned by GitHub; local tests could not run in this environment because the GitHub host was not reachable for cloning. Do not merge without Go formatting/vet/race/build and UI schema/build success plus independent review.
- M1 remains OPEN: independent review and branch governance, actual browser responsive/keyboard/a11y visual acceptance, further security verification and clean CI evidence. Branch protection endpoint returned 403 (integration permission); no rulesets were returned. Owner must configure/enforce governance outside this connector.
- No real orders, credentials, live enablement, production deploy or main branch change. M2 remains next after M1 acceptance.
