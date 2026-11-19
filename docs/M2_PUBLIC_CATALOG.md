# M2 public catalog pagination — read-only transport

`internal/adapters/delta.PublicCatalog.Fetch` requests `GET /v2/products` over HTTPS and follows `meta.after`. This is a **transport-only** slice, not a production-ready exchange integration.

Properties:
- Explicit client timeout and request context, HTTPS-only origin, no credentials, no private endpoint, GET only.
- Max 100 pages, 100,000 records and 4 MiB JSON per page; caller can reduce caps, never raise them.
- Rejects HTTP errors, failed API envelope, malformed JSON/products, empty pages, oversized bodies, invalid/repeated cursors, missing pagination termination and context cancellation.
- Returns no partial snapshot on failure. Raw records are deliberately **not** admitted to the validated domain catalog or labeled tradeable.
- Unit tests use a local TLS fixture; they do **not** prove live Delta response compatibility.

Next: verify Delta India response fixtures, implement exact field-by-field mapper and strict product-kind/capability normalization, validate all required contract fields and status, verify pagination semantics, then integrate an atomic complete snapshot refresh with staleness and observability. Only after that can the UI display authoritative catalog values. No live order or credentials allowed.
