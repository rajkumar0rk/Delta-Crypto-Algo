# Project status — update after each verified milestone

**Last documented:** 2026-11-19, creation of unified `develop` integration branch. Read GitHub for newer commits, PRs, and CI; this file may lag.

## Verified baseline at creation
- `develop` originally points at commit `73d65c5203015a3a49adff2d7c16851287ad3f5b`, containing M1 foundation plus the M5 offline Strategy Builder prototype.
- M1 old PR #1 (against main) and M5 old PR #17 (against M1) contain history now reachable from `develop`; reconcile/close as superseded once newer integration checks pass.
- Latest baseline GitHub Actions CI on M5 commit: successful. New work needs fresh CI; do not reuse baseline result as evidence for later changes.
- Existing UI can configure a versioned JSON draft; no backend strategy storage, market account data, backtesting or live order execution exists.
- New Strategy Workspace is being integrated: it provides an explicit Configure Strategy entry, templates, and editor navigation.

## Incomplete by design
- Market catalog and exchange integration, indicator registry/runtime, strategy persistence, external account auth, risk engine, OMS, P&L, backtests/paper, realistic load testing, full E2E/a11y and real deployment.
- Current Go backend binds localhost with only health endpoints and rejects live mode.
- Oracle VPS is target only; no deployment performed in this work.

## Next check
1. Ensure strategy workspace + agent handbook PR passes Go, TypeScript and schema CI.
2. Merge development-only PR to `develop` once verified and update this file and GitHub issue.
3. Add deterministic UI integration tests/lockfile and complete M2 universal catalog.
