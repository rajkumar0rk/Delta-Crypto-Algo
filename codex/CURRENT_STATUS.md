## 2026-11-24 — M2 atomic financial metadata cache and eligibility checks

- Previous mapping/refresh implementation head `7df86e2a02bba781ebc71b7ff6d2b00ff11b8b2d` passed exact Go/UI CI at https://github.com/rajkumar0rk/Delta-Crypto-Algo/actions/runs/45393761353.
- Added CatalogCache that publishes mapped products and raw fee/margin terms under a single lock, rejects invalid and nonmonotonic refreshes and returns detached read views. TLS tests exercise successive publications, audit changes, caller mutation and failed refresh preservation.
- Added independent account-scope/product/catalog-observation metadata admission, requiring current capabilities and explicit eligibility evidence. Fixed a capability freshness bug: individual evidence age is measured against time of use, not merely set publication. Synthetic tests cover wrong account/environment, stale/future evidence, delisting, missing capability and stale individual evidence.
- New head CI pending. No real account eligibility data, live venue fixture or verified product-specific fee/margin normalization has been obtained. Public endpoint retrieval in this environment was unavailable; only documented field names and synthetic fixtures are used. M1 manual accessibility remains unverified. No main, credentials, orders, production or deployment.

## 2026-11-24 — M2 documented product mapping and audited refresh

- Prior exact head `637d23d6ad69529d4d406d66a04a1536ae398ac5` passed both Go/UI CI: https://github.com/rajkumar0rk/Delta-Crypto-Algo/actions/runs/45298028013. New head CI is pending.
- Implemented documented India Product mapping for vanilla non-quanto futures and call/put options, without symbol-specific rules. Public margin and commission values are retained verbatim, not normalized or interpreted as account fees. Required per-product lot evidence is external to the documented Product schema; missing/stale/mismatched evidence denies publication.
- Added atomic Store.PublishWithChanges and read-only Fetch → MapSnapshot → PublishWithChanges orchestration with synthetic TLS/contract regression fixtures. Source changes, stale observations and invalid mappings fail closed. New local full Go/format checks not performed in this connector call; exact-head CI required. Synthetic tests are not live venue captures.
- Official reference: https://docs.delta.exchange/ (Products and Product schema). Outstanding: real venue captures incl. RWA, verified lot/margin units, inverse/quanto handling, actual fee tiers and account eligibility, durable refresh/cache integration, M1 manual accessibility. Independent human review/protection waived. No main, production, credentials or orders.

## 2026-11-24 — Exact-head CI formatting recovery

- PR #54 head `1ea0e6e5b280f2eb1ea839ee3ec7c5da6598cd38` failed Go formatting in `internal/adapters/delta/catalog_test.go` at [run #45270309502](https://github.com/rajkumar0rk/Delta-Crypto-Algo/actions/runs/45270309502); UI passed. The Go test/build stages did not run.
- Reproduced the missing blank line between the preceding test and the terminal-cursor test's documentation comment with actual local Go 1.23.2 `gofmt -d` on the exact affected test-function excerpt. Applied the resulting one-line whitespace change to the fetched repository file. The complete repository formatting and test suite still require exact-head CI; do not claim local full-suite validation.
- Independent human review and server-side develop protection were explicitly waived by the maintainer for routine development integration on 2026-11-23. They are not completed checks or outstanding approval requests. M1 native 200% zoom and manual screen-reader acceptance remain unverified; M2 official venue product mapping and fee/margin/eligibility/cache acceptance remain incomplete. No merge, live credentials, orders, main or production deployment.

## 2026-11-24 — M2 explicit pagination terminal cursor hardening (PR #54)

- Verified prior exact-head `9b71955f17b9e778274b05a112f241a517d2cf0d` passed GitHub Actions Go/UI at https://github.com/rajkumar0rk/Delta-Crypto-Algo/actions/runs/45260837908. This does **not** validate subsequent commits.
- Fixed public catalog transport to reject an empty JSON-string pagination cursor rather than treating it as terminal. Only explicit JSON `null` proves completion; an empty string must fail closed without publishing partial records.
- Added deterministic HTTPS fixture tests for empty, whitespace, absent and empty-metadata terminal markers. No live venue call or authoritative contract mapping claimed.
- Implementation commits: `31af4c4bb0f6fb7c6b305ab5e8d279f5eadd82e4`, `9090dcb404908175ac3c080a3801da5da22590a0`. **Exact-new-head CI not yet verified** at this handoff; PR remains draft and unmerged. `develop` unchanged. M1 manual governance/accessibility and M2 venue product fee/margin/expiry/account eligibility remain blocked.

## 2026-11-24 — M2 exact-head CI failure triage and formatting repair (draft PR #54)

- Inspected exact-head GitHub Actions for prior commit `a29d53b78915f2529c95a987910d4cd8e7fc99b8`: [run #45240094938](https://github.com/rajkumar0rk/Delta-Crypto-Algo/actions/runs/45240094938) **failed** Go formatting (step 5); Go vet/race/build skipped. UI job succeeded. Previous green CI was for an older commit and does not validate this head.
- Repaired Go formatting in `internal/domain/catalog/store_rejection_test.go` and `internal/adapters/delta/catalog_test.go` on the existing M2 consolidation branch. Final formatting-repair code commit: [ee816582](https://github.com/rajkumar0rk/Delta-Crypto-Algo/commit/ee816582d34509573168f65e2db091a32385ee11). No new functional capability or venue-authoritative fixture is claimed.
- **Current head tests/CI unverified** at handoff; inspect the new workflow, fix remaining failures, then extend M2 contract mapping and fixture coverage. Do not merge PR #54 while M2 acceptance is incomplete. `develop` unchanged. M1 manual assistive-technology acceptance and governance evidence still unresolved. No live credentials, orders, production deployment or `main` changes.

## 2026-11-23 — M2 refresh rejection regression batch (unmerged)
- Added `internal/domain/catalog/store_rejection_test.go` at commit `ea959bf913f70a612e29b00469e1b22b5049f635` on the existing consolidated M2 branch / PR #54. Nine table-driven malformed refresh scenarios assert that the previously published catalog remains intact, stays usable only while fresh, and becomes non-tradeable when stale.
- Tests and exact-head CI for this commit are **not yet verified**; do not treat the new tests as passed or merge M2. No exchange-authoritative contract fixture, account-specific eligibility, fee/margin rates, or full M2 acceptance established.
- M1 PR #46 is already merged into `develop`, but formal manual screen-reader/200% zoom and independent governance acceptance evidence remains missing. Earlier historical entries below describe previous states and must not be read as current.

## 2026-11-23 M1 merge and M2 consolidated preparation
- **M1 PR #46 merged into `develop`** at `e3a24aca6b6b870bb0bcccf06ed7066786e5207a`: https://github.com/rajkumar0rk/Delta-Crypto-Algo/pull/46 . PR source head `7f2d86c` passed Go and UI CI https://github.com/rajkumar0rk/Delta-Crypto-Algo/actions/runs/45111951903 . Post-merge CI has not yet been verified in this status entry.
- Maintainer expressly waived **independent review** and **server-side branch protection** for development integration on 2026-11-23. These controls are NOT configured/verified. Manual screen-reader and native browser 200% zoom acceptance remain **unverified**; M1 is integrated with partial acceptance.
- **M2 consolidation branch** from this exact develop commit incorporates PR #29 (atomic validated store), #30 (plain exact decimals), #47 (audited snapshot transitions), #49 (fee/margin/order capability evidence) and PR #53 (synthetic TLS cross-layer fail-closure and pagination atomicity tests).
- **First combined CI failed Go test fixtures:** [run #45215879440](https://github.com/rajkumar0rk/Delta-Crypto-Algo/actions/runs/45215879440); Go formatting and vet passed, UI passed. New strict content-representation guard exposed synthetic successful JSON handlers with Go's implicit text/plain sniffing; they failed before reaching cursor/timestamp assertions. JSON test fixtures now explicitly declare application/json; fresh exact-head CI remains required. No authoritative Delta India product contract fixtures, full REST mapping, account eligibility, real margin/fee rates, or durable refresh/cache are claimed. Transport safeguards from PRs #48/#50/#51 have now been reconciled into this draft: cursor validation, timestamp/cancellation atomicity, nested JSON key rejection, and HTTP content representation validation. Original PRs remain open until combined-head CI and integration review. No M2 merge before further acceptance.
- Metadata-only tests, no credentials, orders, live-trading enablement, production deployment, `main` changes or other automation modifications.

## 2026-11-24 — M2 Go formatting CI follow-up

- PR #54 exact head `a6231e75fa0027c4c3cca1e10045ac108b31d13a` failed Go formatting (`internal/domain/catalog/store_rejection_test.go`) in [run #45249421975](https://github.com/rajkumar0rk/Delta-Crypto-Algo/actions/runs/45249421975); UI passed, Go vet/race/build skipped. This is the **second** formatting failure, not a passing validation.
- Expanded the two duplicate-product fixture closures into gofmt-compatible multiline bodies in commit `deb0d78c8f1ca1b47f1f8e669682beefb2bb090e` on existing M2 consolidation PR #54. Exact new-head CI is pending; no claim of passing tests.
- M1 independent review, `develop` protection and manual screen-reader/native zoom acceptance remain unverified. M2 still lacks authoritative venue contract/fee/margin/eligibility fixtures. No merge, `main`, live orders, secrets or deployment.

# Project status — M1 consolidated, acceptance pending

## 2026-11-24 — M2 catalog transport preflight hardening (draft PR #54)

- On existing consolidation branch `feature/m2-catalog-domain-consolidated-20261123`, reject opaque HTTPS origins and reject already-canceled contexts before each pagination request. Added deterministic counting-transport regression tests proving no network I/O for canceled requests or malformed origins.
- Code commit: https://github.com/rajkumar0rk/Delta-Crypto-Algo/commit/7cf6b93fdc1b2499bcb74e67cf4447276048acc2 ; tests commit: https://github.com/rajkumar0rk/Delta-Crypto-Algo/commit/54d187c2bf9bac4a71bba69b3d6aa2cb2d0a5f16 .
- **CI for these new commits not yet verified at this handoff**; older green run #45216268002 covers an earlier head only. No local test execution claimed.
- M1 code merged via #46 but formal independent review, protection and manual assistive-tech/zoom acceptance remain unverified. M2 still lacks authoritative venue contract/fee/margin/expiry mapping and eligibility proof. Do not merge #54 or enable trading on this evidence.


## 2026-11-23 — M1 accessibility, review governance and current CI

- Existing PR #46: https://github.com/rajkumar0rk/Delta-Crypto-Algo/pull/46 targeting `develop`, not merged.
- Implemented drawer `role=dialog`/`aria-modal`, inert background/skip link, deferred focus restoration and Chrome regression tests for keyboard containment, modal state and viewport transitions.
- Added `.github/CODEOWNERS`, PR template, `docs/CONTRIBUTING.md`; clarified delegated development-only merge authority versus separate live/main permissions in `docs/SECURITY.md`.
- Verified automated Go+UI CI **passed** at code/governance head `7b4320fe91c73cdc28a3d8cde48af9ae8c4e7c5f`: https://github.com/rajkumar0rk/Delta-Crypto-Algo/actions/runs/45092813051. Revalidate exact final HEAD after these documentation commits.
- External acceptance gates unresolved: zero independent reviews on PR #46, `develop` observed unprotected and rulesets absent, protection administration unavailable to GitHub integration, human screen-reader/native 200% zoom evidence absent.
- `develop` checked at `60fb41ed034ad86fe401623472d3d81b67d116bf`. No merge attempted; no main, real credentials, orders or deployment.
- **Decision:** PR ready for independent acceptance, not M1-complete. Resume M2 independently testable preparation if external acceptance gates remain unavailable; do not call M1 complete or waive required review under solo-maintainer delegation.


## 2026-11-23 — M1 responsive navigation and automated acceptance follow-up

- **User governance clarification:** the sole developer has delegated routine development decisions and eligible `develop` merges; do not repeatedly request personal milestone approval. The existing independent-review/branch-protection/assistive-technology gates are not silently marked passed. Preserve team-ready governance for future contributors.
- **Code:** existing M1 PR #46 adds mobile drawer keyboard containment (including programmatically displaced focus), excludes the scrim from tab/AX navigation, closes stale drawer state on the desktop breakpoint, and restores focus to the visible menu trigger on the mobile breakpoint.
- **Tests:** Chrome browser smoke now exercises both responsive transitions and keyboard focus recovery; it retains 320px narrow-reflow, viewport snapshots and AX button-name checks. Do not label this real 200% browser-zoom or screen-reader acceptance.
- **Observed CI:** https://github.com/rajkumar0rk/Delta-Crypto-Algo/actions/runs/45091601256 had successful Go checks but UI failed because headless Chrome did not report an applied CSS zoom value. The nonrepresentative assertion was removed; updated code head `d7bab0505fa23a789d4972d2dc0eccc881a51063` has **unverified exact-head CI** at the time of this handoff. Check newer documentation-only commits and the resulting final-head CI before merge.
- **Integration:** `develop` remained `60fb41ed034ad86fe401623472d3d81b67d116bf` before these PR changes, PR #46 remains open and unmerged. No trading, credentials, `main`, or production deployment.


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

## 2026-11-23 — M1 numeric loopback hardening (PR #46, pending exact-head CI)
- Foundation listener now accepts only literal `127.0.0.1` or `::1`, rejecting `localhost`/case aliases because resolver/hosts mappings are not guaranteed loopback. Demo IPv6 acceptance and alias rejection covered in deterministic configuration tests.
- Implementation commits `bdd9d584dc77cc486aa1631cfb9de480394b9db9`, `5c3547dbf9013fef3eaae6c63fe84e3029a61c4f`, acceptance docs `420d8853281758cab04c847384b49156cf017964`.
- GitHub Actions lookup returned **no runs** for documentation head at inspection time; no local Go tests were executed. Do not treat previous green runs as validation of these changes.
- M1 remains blocked by independent review, owner-only `develop` governance, manual screen-reader and 200% zoom acceptance. No merge, `main` changes, credentials, live orders or deployment.

## 2026-11-23 — M1 HTTP response security regression (PR #46; not merged)
- Added deterministic table-driven tests on the M1 consolidation branch covering CSP, CORP, nosniff, frame denial, referrer policy, cache disabling, generated request IDs, and no CORS grant across liveness, unready, denied-method, missing-order and preflight paths. Commit `6bb80c4d72344afcc1b6ef74c5d474da62966607`.
- Prior PR head `c8a6ab4ed221ffb128285a908c74ab44fc885b4d` has successful Actions run https://github.com/rajkumar0rk/Delta-Crypto-Algo/actions/runs/45067324239. New test commit had no workflow run returned at inspection; do not call it validated or merge until exact-head CI passes.
- M1 acceptance still requires independent review, owner-verified develop branch governance, and manual screen-reader/200% zoom acceptance. No live access, orders, deployment, or main changes.

## M1 CI formatting repair (2026-11-23)
- PR #46 exact prior head `39b26cd6d248ba54d491216f3df1ef9cd5bfbd54` failed Go formatting (only `internal/httpapi/server_test.go`); UI succeeded: https://github.com/rajkumar0rk/Delta-Crypto-Algo/actions/runs/45072767981
- Commit `3c7f6c1b35d3f3389c1c0f60c04f8e8f3ceed864` aligns the security-header test table with gofmt. **Fresh exact-head CI remains required**; the failed run does not validate Go tests/race/vet/build.
- M1 remains unaccepted until independent review, owner-confirmed `develop` governance and manual screen-reader/200% zoom acceptance. Do not merge or begin M2 integration on green CI alone.

## M1 health contract regression batch (2026-11-23)
- PR #46: added table-driven coverage for JSON health state/trade_enabled=false, liveness/readiness statuses, HEAD routing, unsafe method denials, unknown order paths, no-store and denied cross-origin responses.
- Code head `5f454852446b081c3b3d20974696fa12412bf6b3` was pushed; no exact-head workflow returned immediately after push. Go tests were NOT run locally. This is pending CI and independent review, not accepted M1.
- Owner-only branch protection verification and manual screen-reader/200%-zoom acceptance remain outstanding. No `develop`, `main`, credentials, orders or production changes.


## M1 security regression batch — 2026-11-23
- PR #46 remains the M1 integration candidate; new commits add a deny-by-default browser Permissions-Policy and Go regression tests for server-generated unique request IDs across successful/denied/unknown routes, client-supplied request ID spoofing, denied unsafe methods and absence of permissive CORS preflight headers.
- **Exact new-head CI not yet verified.** Prior green run #45092945033 validates an older head only and MUST NOT be used to merge this update. Local Go/UI test execution was not performed through the GitHub connector.
- M1 remains blocked by independent qualified review, develop governance/protection verification, manual screen-reader/200% zoom acceptance and exact-head CI. Do not merge until gates are satisfied.
- M2 PRs #29/#30/#47–#51 require consolidation after M1; authoritative Delta India contract/margin/fee/expiry fixtures remain missing. No live credentials, orders, production or main changes.

## M1 exact-head CI diagnosis — 2026-11-23
- PR #46 head `1e29961dc3eafcd7779499b33296074d0ba128be` failed Go formatting while UI passed: https://github.com/rajkumar0rk/Delta-Crypto-Algo/actions/runs/45106026995. The Go job explicitly reported `internal/httpapi/server_security_test.go` from `gofmt -l`; no Go tests were reached.
- Corrected the `Permissions-Policy` map literal's one-column alignment in commit `95762feb4550c50670ca81fef3eb830a3a463364`; this is a targeted formatting fix, not evidence of successful CI. The latest head must be rechecked after this status commit.
- Independent review count remains zero; branch governance and native screen-reader/200%-zoom manual acceptance remain external gates. M1 not accepted and no merge authorized until all gates pass. M2 integration remains deferred.
