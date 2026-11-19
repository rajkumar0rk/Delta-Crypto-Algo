# Delta Crypto Algo — agent instructions (authoritative)

Read this file **before modifying code**. Then read `codex/AGENTS.md`, `codex/PROJECT_CONTEXT.md`, `codex/CURRENT_STATUS.md`, `docs/MILESTONES.md`, and the relevant module specification.

## User intent
Build a **large, enterprise-quality, independent trading platform** for ALL products supported by **Delta Exchange India**, including available crypto and RWA-linked futures/options, not a BTC-only bot. The user should NOT need to invent missing requirements: agents must proactively reason through operational, financial, security, UX, performance, correctness and failure details. The interface must be accessible, responsive, consistently aligned, usable by non-programmers, and fast under measured loads.

## Git workflow — user authorized autonomous development merging
- **`develop` is the single integration branch.** All new work branches from `develop`, uses focused feature branches, and merges into `develop` after required reviews/green CI and tests. Agents may autonomously merge development-only PRs to `develop` under these gates.
- **`main` is a protected conceptual release branch.** Do not push or merge to `main` automatically. No rebase or force pushes to shared branches. Changing default branch or production settings requires explicit authorization.
- Review existing PRs and branch SHA before work to avoid duplicate or conflicting changes. No auto-merge with failed/unknown checks. Record completed work in `codex/CURRENT_STATUS.md` and issue/PR notes.
- Use GitHub checks, not confident prose, as proof. If CI or access is blocked, stop and report it.

## Hard non-negotiables
1. **NO REAL TRADING OR LIVE AUTHORIZATION.** Do not connect real trading credentials, submit real orders, enable live strategies, perform actual production deploys or withdraw funds without current explicit human approval for that specific action. This permission to merge code is NOT permission to trade.
2. Preserve fail-closed startup and risk boundaries. Risk controls must be independent and cannot be bypassed from the UI or strategy JSON.
3. No secrets, account data, production logs, customer PII or credentials in GitHub, UI bundles, test data or prompt files.
4. Do not claim bug-free, profitable, benchmarked, exchange-compatible or production-ready unless verified evidence exists.
5. Never hardcode BTC/ETH or assume tick/lot/margin/expiry/fees. Model venue-specific product metadata dynamically and accurately.
6. No arbitrary expression execution from JSON. Use typed, declarative, bounded DSL with a registry and deterministic tests.
7. Every change must include acceptance criteria, negative cases, edge cases, performance/accessibility/security impacts and docs. Prefer minimal coherent improvements to fake broad implementation.

## Code map
- Go: `cmd/`, `internal/`, versioned domain and future services.
- UI: `apps/control-ui/src/`; Strategy workspace `apps/control-ui/src/strategy/`.
- DSL schema and examples: `schemas/`, `examples/`.
- Product/architecture docs: `docs/`; agent handoff: `codex/`.

## Local gates
```bash
go test -race ./...
go vet ./...
cd apps/control-ui && npm install --no-audit --no-fund && npm run validate:schema && npm run build
```
Any incomplete or unavailable test gate must be explicitly disclosed. See `codex/ENGINEERING_STANDARDS.md`.

## Required delivery report
Milestone; exact commit / PR; what was implemented; tests and their actual outcomes; what was NOT done; known risks; next prioritized action. Do not silently merge a partially validated feature.
