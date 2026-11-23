# Engineering contribution and review workflow

This repository currently has one primary maintainer, but is designed for additional engineers. The user has already authorized routine autonomous development and eligible, **reviewed** merges to `develop`; additional personal approval of each development step is not needed. This does **not** waive independent review requirements or authorize live trading, production deploys or `main` releases.

## Branches and code ownership

- `develop` is the sole feature integration branch; start each feature/fix branch from its freshly verified HEAD.
- `main` is reserved for explicitly authorized releases. Never force-push or use automated merges to `main`.
- CODEOWNERS routes review requests; it is **not proof** that an owner (especially the PR author) independently approved the change.
- When more developers join, assign bounded context owners and at least one qualified reviewer who is **not** the author for risk, OMS, security and money-sensitive changes. Document reviewer credentials and review scope in the PR.
- Repository administrators should protect `develop` with required Go/UI checks, a non-author review requirement, blocked force pushes/deletions and restricted bypass. CI and CODEOWNERS files alone **cannot enforce** GitHub server-side rules.

## Minimum pre-merge evidence

1. Review instructions in root `AGENTS.md` and `codex/`, inspect open PRs and the current `develop` commit to avoid duplicate implementation.
2. State observable acceptance criteria, security impacts, negative cases and rollback plan. Include typed contract units and source/provenance for venue-specific fields.
3. Run `go test -race -count=1 ./...`, `go vet ./...` and Go build; for UI run `npm ci --no-audit --no-fund`, `npm run validate:schema`, `npm run build` and browser smoke where relevant.
4. Record actual test outcomes and the exact-head CI run URL; never infer green CI from an earlier commit.
5. Obtain any mandated independent approval and manual acceptance; document unresolved deficiencies honestly. A self-review can help but must not be represented as independent verification.
6. Merge only a reviewed, green, low-risk PR to `develop` through GitHub's normal merge API, preserving repository protection and checks.

## Safety model

No deployment, exchange keys, real orders, live mode, withdrawals or `main` releases without separate explicit approval for that action. Strategy JSON is not allowed to bypass hard risk gates. Failed or unknown exchange acknowledgments must not trigger blind opening-order retries. Unsupported catalog metadata remains unknown, not a guessed default.

## Accessibility evidence

Automated Chrome accessibility-tree, keyboard and narrow-reflow checks are useful, but not the same as a human screen-reader session or native 200% browser zoom test. Record actual browser, viewport, assistive technology, tested flow, problems and reproduction steps; leave unperformed checks pending rather than inferring success.
