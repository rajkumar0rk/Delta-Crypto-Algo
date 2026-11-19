# M2 Catalog normalization — initial domain slice

This module defines the **internal product metadata contract** for all venue-listed underlying types, including crypto and RWA-linked derivatives. It deliberately does **not** hardcode BTC/ETH symbols, contract fees, leverage or margin rules.

## Current code
`internal/domain/catalog` validates immutable snapshots and per-product fields, including exact positive decimal strings for tick/lot/multiplier, call/put strike and expiry, dated vs perpetual futures, unique IDs/symbols, complete pagination, staleness and active status. The `Tradeable` check is a **necessary but insufficient** condition for any future order; an independent account capability, margin and risk gateway is mandatory.

## Not implemented
No Delta REST client, `meta.after` pagination, persistence, HTTP catalog API, per-account restrictions, contract classification mapping, fee/margin API, RWA asset taxonomy or live order routing. Do not use `Tradeable` as authorization. Exchange product status and field names require adapter fixtures from actual read-only responses and API documentation before ingestion. No real network requests or credentials are needed for this slice.

## Next slice
Implement a bounded, cancellable read-only REST adapter that handles every page and partial failures, snapshots raw responses, maps authoritative fields and fails closed on unknown product types. Add httptest fixtures for cursor loops, duplicate pages, rate limiting, 5xx, stale and delisted products. Keep external exchange data untrusted and never infer missing tick/lot/multiplier.
