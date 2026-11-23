# M1 Foundation Acceptance Ledger

M1 is an inert control-plane foundation, not an order engine, catalog, persisted strategy library or deployment approval. The platform aims to support all eligible Delta India crypto futures/options and RWA-linked derivatives; no hardcoded BTC assumptions.

## Mandatory acceptance gates

| Gate | Required evidence | Current state |
|---|---|---|
| Startup | Reject live mode, non-numeric-loopback binding (including DNS-based localhost aliases) and noncanonical/out-of-range TCP ports | Passed combined CI run 44905852868 |
| HTTP boundary | Liveness 200, readiness 503, trade_enabled explicitly false, representative trading routes return 404 | Passed combined CI run 44905852868 |
| API isolation | Browser security policy, no open CORS, server-generated unpredictable request IDs | Passed combined CI run 44905852868 |
| Go quality | Module verify/tidy, format, vet, race suite, build | Passed combined CI run 44905852868 |
| UI reproducibility | Exact direct dependencies, lockfile, npm ci, type-check, schema examples | Passed combined CI run 44905852868 |
| Runtime dependency audit | npm audit high severity for production dependencies | Passed combined CI run 44905852868 |
| Browser | Chrome screenshots at 360,768,1280,1920; horizontal overflow, labels, skip link, disabled live banner, navigation smoke | Passed CI 44905852868: 320px narrow reflow, 360/768/1280/1920, reduced-motion, AX accessible names, keyboard skip link, mobile drawer Enter/Shift+Tab/Escape and focus restoration; screenshots visually inspected; manual screen-reader and real 200% browser zoom still pending |
| Independent review | Approved by a person other than the author; no high/critical concerns | Not verified |
| Branch governance | Protect develop; required checks, independent review, no force push/delete or bypass | Owner-only setting, not verified (GitHub 403) |
| M1 handoff | Final SHA, CI run, review, acceptance evidence, merge SHA in CURRENT_STATUS | Pending |

## Security and scope

Live trading, public bind, real orders and credentials remain disallowed. CI uses read-only permissions and nonpersistent checkout credentials. GitHub Actions references are pinned to immutable commits verified against upstream v4/v5 GitHub tag refs on 2026-11-21. Review updates to these immutable pins as dependencies change. A Chrome screenshot smoke test is not equivalent to independent visual inspection, screen-reader evaluation or complete keyboard testing.

## Technical review observations (2026-11-21)

Author-side review found and fixed: literal README newline, off-screen mobile drawer remaining keyboard-focusable when closed, missing accessible menu focus handling and Escape/Tab trapping, untested narrow-viewport/accessible-name behavior, and incomplete security headers on an unlikely request-ID RNG failure. These changes passed both Go/UI jobs on 44905852868. **This review is not independent approval.**

## Repository owner action

In Settings > Rules > Rulesets or Branch protection, target develop and require one independent approving review plus the Go/UI checks; reject forced pushes/deletions and unauthorized bypass. Confirm those settings using owner access; the agent connection returns HTTP 403 for branch protection reads.

Combined code/CI checks passed; do not call M1 complete or move the dependency gate to M2 until independent review and develop governance are verified.

## Loopback DNS ambiguity hardening (2026-11-23)

M1 foundation startup must accept only numeric IPv4 `127.0.0.1` or IPv6 `::1` listeners. A hostname such as `localhost` depends on external resolver/hosts configuration and is rejected even if it usually resolves to loopback. Regression coverage includes IPv6 acceptance, DNS aliases and port bounds. This change requires fresh exact-head CI and independent review; it does not satisfy manual screen-reader or owner governance gates.

## Autonomous solo-maintainer M1 UI review — 2026-11-23

The user has **already delegated routine development decisions and eligible merges to `develop`**. Repeatedly asking the sole developer to personally approve each development milestone is not a requirement. This delegation does not supply an independent human code review or authorize bypassing a GitHub-enforced gate. Additional developers will later use the same feature-branch, PR, review, and CI workflow.

PR #46 now includes responsive mobile navigation changes: prevent keyboard focus escaping the open drawer; mark the backdrop as outside the accessibility and sequential tab order; close stale mobile drawer state when switching to desktop; and restore focus to the visible menu trigger when switching from desktop to mobile. Automated Chrome coverage includes both breakpoint directions and forced external focus, plus the existing 320px narrow-reflow and reduced-motion check.

CI investigation:
- https://github.com/rajkumar0rk/Delta-Crypto-Algo/actions/runs/45091601256 — Go passed; UI browser smoke failed solely on a new **headless CSS zoom assertion**. Chrome reported computed zoom `1` despite the inline style reading `200%`. This was **not valid evidence of native browser zoom**.
- The unreliable zoom assertion was removed; the deterministic 320 CSS-pixel reflow check remains, and **manual native browser zoom and screen-reader acceptance remain unverified**.
- Updated code head `d7bab0505fa23a789d4972d2dc0eccc881a51063` requires fresh CI; no green result is claimed here.

M1 remains pending independent code-review evidence, repository governance and manual assistive-technology acceptance. No live trading, production credentials, VPS deployment or `main` changes are authorized by this development delegation.
