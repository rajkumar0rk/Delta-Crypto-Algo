# M2 product-specific capability evidence (development contract)

This is a **metadata-only, fail-closed primitive**. It is not an exchange-verified mapping, account eligibility result, order authorization, or live risk control. There are no production capabilities populated in this repository.

## Evidence required before admitting a product

Each venue product ID (including each option strike and expiry) needs independently verified source, observation time, fee ratios, initial/maintenance margin ratios, and supported order types. A missing field is **unknown**, never zero. An explicitly verified zero is represented by the decimal string `0`. Stale evidence and unknown products are denied. Product IDs, not BTC/ETH or underlying-name defaults, are the lookup key.

Fee/margin ratios are bounded exact decimal strings in [0,1]; the current primitive intentionally cannot represent tier schedules, account-specific discounts, portfolio margin, or dynamic collateral haircuts. Those require additional versioned evidence and a separate account/risk layer. Do not substitute this static model for actual venue calculations.

## Acceptance before integration

- Confirm Delta India product endpoint response shape and pagination with authoritative fixtures.
- Verify all supported crypto and RWA derivative product mappings; reject unrecognized classifications.
- Confirm fees, GST applicability, tiering, margin modes, liquidation requirements, option writer obligations, and order-type capabilities against documented venue/account rules.
- Wire verified snapshots and capabilities into a bounded refresh/reconciliation process and UI read-only states.
- Maintain independent M7 risk and M8 OMS authorization; passing this primitive never places an order.

No real credentials, private API calls, orders, or deployment are involved.
