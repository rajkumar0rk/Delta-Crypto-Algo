# M2 public product mapping and refresh

The public Product contract is documented at https://docs.delta.exchange/ (India API). The adapter recognizes the four documented futures/options types, maps nested asset symbols, tick size, contract value, settlement time and strike. Unknown or unsupported contract value models are rejected. Unknown/non-operational/reduce-only states cannot admit new risk. Per-product lot evidence is required; the documented Product schema does not establish a lot-step field. Synthetic tests are not captured exchange records.

PublicTerms preserves published fee/margin strings verbatim. It is not a fee calculator, a normalized margin ratio, account permission or execution authorization. Do not populate CapabilitySet from this alone. Lot evidence must be supplied from a trusted per-product source, tied to ID, symbol, contract value, observation time and lot increment. Refresh is read-only, returns terms and audit changes for the published observation, and does not establish durable storage. Invalid fetch/mapping leaves the last-good snapshot untouched; consumers must check age before using it. No private API or order endpoint is called.

Remaining acceptance: authoritative actual responses and product-specific size constraints; validate margin units, tiering, option writing and account eligibility, inverse/quanto models where supported, and durable refresh/reconciliation. No live trading allowed.

## Atomic cache and account metadata admission

CatalogCache publishes normalized products and unnormalized public fee/margin fields as one in-memory observation. Refresh failure preserves both. Read returns detached copies; callers still enforce freshness. Contract-change audit covers normalized Product fields; raw financial term changes are retained in the current view but are not yet separately persisted or classified. CheckMetadataAdmission requires explicit account/environment/product evidence tied to the exact catalog source and observation, plus independently verified capabilities no older than the current catalog. All evidence is synthetic in tests; no actual account or venue authorization is implied. Durable cache and audit persistence, authoritative fee/margin interpretation and real product eligibility remain acceptance work.

## Public contract observation workflow

The separate public-catalog-evidence workflow runs a standard-library Python probe against the official India public products endpoint, bounded to two live records per documented derivative type. It logs the observation timestamp, envelope pagination shape, and selected product fields, never credentials or private account state. It does not follow redirects or use results as trading permissions. External network availability is deliberately separated from deterministic foundation CI. The probe is sampling evidence, not a full listing, proof of lot-size or a fee/margin contract. A failed probe is a failed observation, not an accepted fixture.
