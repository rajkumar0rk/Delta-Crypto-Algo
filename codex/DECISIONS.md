# Architecture decision record (living)

| Date | Decision | Why | Alternative / limitation |
|---|---|---|---|
| 2026-11-19 | Build independently from scratch | Prevent cross-project coupling and accidental reuse | Other Delta repos are out of scope |
| 2026-11-19 | `develop` integrates all development branches | Single place for agents and coherent integration | `main` remains manual release branch |
| 2026-11-19 | Agents can merge tested, low-risk development PRs into `develop` | User explicitly authorized self-merging | No permission to live trade, deploy or merge `main` |
| 2026-11-19 | Go backend + React/TypeScript UI + planned PostgreSQL | Performance, typed contracts and maintainable boundaries | Not all infrastructure exists yet |
| 2026-11-19 | Dynamic full venue product catalog, not BTC-only | Futures/options/RWA coverage and new listings | Account eligibility and contract capabilities remain venue authoritative |
| 2026-11-19 | Configure Strategy = actual guided app workflow | User should not need to enumerate every field | Current UI only exports offline drafts |
| 2026-11-19 | Versioned declarative strategy DSL | Replay, audit, safety and UI generation | Avoid arbitrary execution; feature registry/server validation forthcoming |
| 2026-11-19 | Fail-closed live startup and independent risk gateway | Real-money safety | Production readiness requires separate explicit human go-live approval |
| 2026-11-19 | Oracle VPS pilot, not enterprise HA | User deployment preference, realistic resource limits | Scale-out/HA must be designed and tested later |

Changes to architecture or financial semantics require updating this record and relevant tests/documents.
