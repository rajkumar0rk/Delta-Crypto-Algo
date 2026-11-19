# Project context — stable requirements

Repository: `rajkumar0rk/Delta-Crypto-Algo`. New independent implementation from scratch; **do not copy or modify** the user's other Delta projects unless explicitly instructed. Backend: Go, frontend: React/TypeScript. Persistent system of record planned: PostgreSQL. Hosting target: Oracle Cloud VPS (small pilot first, not HA).

User wants an **enterprise-scale** application, fully configurable, scalable, high-performance, highly secure, user-friendly, precision-aligned UI for trading on **Delta Exchange India**. All crypto and supported RWA-linked derivative instruments are in scope; options and futures have distinct contract, risk, premium, margin and P&L models. Dynamic catalog, expiry handling and newly listed assets are mandatory.

The user is not expected to know or enumerate every small detail. Agents must **discover, design, document, prioritize and implement** the full professional lifecycle of configuring/testing/operating algo strategies. Do not wait to be explicitly told about edge cases, security, validation, operational recovery or UX polish.

## Target product journey
Browse strategies -> click **Configure Strategy** -> select blank/template -> step-by-step no-code editing -> validation and review -> save versioned draft -> reproducible backtest -> paper/demo deployment -> risk review/approvals -> monitored live deployment (only later, manually authorized) -> audit, analytics and rollback.

## Trust hierarchy
User approval/current explicit request > hard safety limits > `AGENTS.md`/codex guidelines > milestone backlog. A schema-valid draft is NOT safe/approved for trading; a backtest result is not a guarantee of profit.

## Current state pointer
Read [CURRENT_STATUS](CURRENT_STATUS.md) every run. The state changes often: follow repository actual branch/CI evidence over stale narratives.
