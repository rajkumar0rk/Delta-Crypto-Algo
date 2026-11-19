# Milestone execution playbook for scheduled agents

1. Inspect `develop` HEAD, open PRs/issues, CI and latest `codex/CURRENT_STATUS.md`. Resolve conflicting work before writing.
2. Choose highest-priority unblocked slice of the 15-milestone roadmap. User expects professional feature discovery: add missing requirements to the backlog before implementing code.
3. State acceptance tests and failure cases, implement on `feature/<milestone>-<topic>` from CURRENT `develop`, keep commits small and reversible.
4. Run unit/static/UI/semantic tests and review own diff for credential leaks, UI regressions, denial states, product units and state consistency.
5. Open PR `head=feature/*`, `base=develop` with complete evidence. Inspect actual GitHub CI status after the latest commit.
6. User **authorizes autonomous MERGE INTO develop** only when all applicable tests and review gates pass, with no known high/critical bug, no dangerous live behavior and no conflicts. Do not force merge. When gates cannot be verified, leave PR unmerged, record blockers and address in next scheduled run.
7. Never auto-merge to `main`; never publish/deploy production, connect live credentials, activate trading or place orders without a separate explicit user command.
8. Update `codex/CURRENT_STATUS.md` and relevant issue with merged commit, evidence, precise gaps, next milestone; avoid simultaneous agents clobbering the same status file.
9. Send concise user-facing report with links; avoid "done" if CI failed or only draft design exists.

## Schedule semantics
A scheduled ChatGPT automation invokes **one run of work at a time**, not a permanent daemon, and may require GitHub permissions available on that run. Prioritize useful PRs and honest status. Do not modify unrelated schedules.

## Workstream priority
- Finish M1/M5 integration and close superseded PRs cleanly.
- Complete M2 catalog + safe capability map before depending on market-specific assumptions.
- Complete M3 market data integrity, M4 auth/account and M5 persistent authoring before starting live execution.
- M6 indicators/signals, M7 independent risk, M8 OMS/exchange, M9 accounting, M10 backtest/paper, M11 UI depth, M12–M15 ops/perf/deployment/live signoff.

## Git conflicts
Do not force push shared integration branch. If branch moved, fetch latest, compare, create new PR from latest develop; if connector access can't perform a safe merge/rebase, leave work on PR, disclose blocker.
