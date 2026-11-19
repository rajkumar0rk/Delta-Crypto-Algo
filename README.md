# Delta Crypto Algo

Enterprise-oriented, safety-first algorithmic trading platform for **Delta Exchange India** futures and options.

> **STATUS: FOUNDATION ONLY — LIVE TRADING DISABLED.** No production trading or funded deployment is authorized by this repository.

This project is being developed independently from scratch, using a milestone-based, test-gated workflow.

- Architecture: modular Go services / ports-and-adapters, event-driven integration, React + TypeScript control panel, PostgreSQL system of record.
- Delivery: one milestone at a time, via feature branches and pull requests; no automatic production deployment.
- Security: no API keys, secrets, account data, or real orders in the repository.
- Runtime: configurable paper/demo/live environments; the `live` mode requires explicit human activation and passed release gates.

See `docs/` on the milestone feature branch for detailed design and acceptance criteria.
