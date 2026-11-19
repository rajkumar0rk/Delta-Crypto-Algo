# Codex and autonomous agent operating contract

## Start here
1. Read repository root [AGENTS](../AGENTS.md), then [context](PROJECT_CONTEXT.md), [current status](CURRENT_STATUS.md), [product requirements](PRODUCT_REQUIREMENTS.md), [strategy design](STRATEGY_BUILDER_SPEC.md), [quality gates](ENGINEERING_STANDARDS.md), [milestone runbook](MILESTONE_RUNBOOK.md), and [decision log](DECISIONS.md).
2. Inspect `develop`, PRs, CI and issues in GitHub. `main` is a release boundary, NOT the development target.
3. Select one small implementation slice with evidence-based acceptance criteria. Avoid unrelated rewrites and avoid touching sibling Delta repositories.
4. Work on a feature branch from the CURRENT `develop` HEAD; test; review; merge to `develop` only after gates. Capture new details in docs and handoff files.

## Role split for complex work
- Architect: contracts, bounded contexts, failure semantics, backward compatibility.
- Frontend/UI: flows, labels/help, validation, responsive alignment, keyboard/screen readers, state handling.
- Trading domain: deterministic market semantics, risk/margin, order state machines, accounting.
- QA/SRE: high-pressure tests, race conditions, uptime/restart, benchmarks, backup/restore, adverse venue behavior.
All roles have a duty to raise missing requirements and security concerns without waiting for the user to identify them.

## Typical turn workflow
Read -> identify gap -> design tradeoffs -> implement -> run quality checks -> open PR -> review/merge to develop if green and low-risk -> update status and issues -> post concise facts. When blocked, record honest blockers. No empty "in progress" commits, fabricated tests or claimed performance.

## Special safety rule
Autonomous code integration is allowed; autonomous **live trading, key setup, VPS production deployment, financial transfers, main release and go-live approval are not allowed**. Always use demo/paper/stubs for testing unless the user explicitly authorizes a specific higher-risk action.
