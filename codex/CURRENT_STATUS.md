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
