# M1 Foundation Acceptance Ledger

M1 is an inert control-plane foundation, not an order engine, catalog, persisted strategy library or deployment approval. The platform aims to support all eligible Delta India crypto futures/options and RWA-linked derivatives; no hardcoded BTC assumptions.

## Mandatory acceptance gates

| Gate | Required evidence | Current state |
|---|---|---|
| Startup | Reject live mode, non-loopback binding and noncanonical/out-of-range TCP ports | Passed combined CI run 44904889053 |
| HTTP boundary | Liveness 200, readiness 503, trade_enabled explicitly false, representative trading routes return 404 | Passed combined CI run 44904889053 |
| API isolation | Browser security policy, no open CORS, server-generated unpredictable request IDs | Passed combined CI run 44904889053 |
| Go quality | Module verify/tidy, format, vet, race suite, build | Passed combined CI run 44904889053 |
| UI reproducibility | Exact direct dependencies, lockfile, npm ci, type-check, schema examples | Passed combined CI run 44904889053 |
| Runtime dependency audit | npm audit high severity for production dependencies | Passed combined CI run 44904889053 |
| Browser | Chrome screenshots at 360,768,1280,1920; horizontal overflow, labels, skip link, disabled live banner, navigation smoke | Browser CI passed 44904889053; screenshots visually inspected; full accessibility review pending |
| Independent review | Approved by a person other than the author; no high/critical concerns | Not verified |
| Branch governance | Protect develop; required checks, independent review, no force push/delete or bypass | Owner-only setting, not verified (GitHub 403) |
| M1 handoff | Final SHA, CI run, review, acceptance evidence, merge SHA in CURRENT_STATUS | Pending |

## Security and scope

Live trading, public bind, real orders and credentials remain disallowed. CI uses read-only permissions and nonpersistent checkout credentials. GitHub Actions references are pinned to immutable commits verified against upstream v4/v5 GitHub tag refs on 2026-11-21. Review updates to these immutable pins as dependencies change. A Chrome screenshot smoke test is not equivalent to independent visual inspection, screen-reader evaluation or complete keyboard testing.

## Repository owner action

In Settings > Rules > Rulesets or Branch protection, target develop and require one independent approving review plus the Go/UI checks; reject forced pushes/deletions and unauthorized bypass. Confirm those settings using owner access; the agent connection returns HTTP 403 for branch protection reads.

Combined code/CI checks passed; do not call M1 complete or move the dependency gate to M2 until independent review and develop governance are verified.
