# Delta Exchange India API integration contract

Authoritative docs: https://docs.delta.exchange/

- Base REST URL: `https://api.india.delta.exchange/v2`.
- `GET /products`: includes perpetual/dated futures and calls/puts; accepts `states`, `contract_types`, `after` and `page_size`.
- Product pagination response includes `meta.after`; keep fetching until null; detect duplicate cursor and partial sync.
- Carry exchange-assigned product ID, symbol, state, trading_status, contract_type, tick_size, contract_value, margin, fee and funding specifications, expiry and strike where relevant.
- Public WebSocket newer endpoint: `wss://public-socket.india.delta.exchange`. Legacy channel transitions must not be assumed backwards compatible.
- Private channels and auth need independent tests. Use Delta product metadata rather than assuming BTC price point value applies to ETH, RWA or other products.
- `client_order_id` max 32 chars; exact idempotency/replace semantics must be verified.
- Products can become `upcoming`, `expired`, disrupted, post-only, cancel-only, or reduce-only; default to NO NEW RISK.
- REST order status, fees, account positions and fills are authoritative only after reconciliation.

**M1 has no exchange-connected components.** These are researched design inputs for the implementation/verification milestones.
