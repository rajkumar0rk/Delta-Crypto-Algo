# Performance engineering

Performance is a measured contract, not an adjective.

## Backend
- Measure tick receive -> normalization -> candle/signal -> risk -> durable order intent -> exchange send. Record distinct p50/p95/p99 latencies and worst-case backlog.
- Bounded queues, backpressure policies, worker ownership, ring buffers for hot market data, batching where safe. Never silently drop events needed for trading state.
- Immutable event IDs, narrow allocations, profiling, race tests and context-based cancellation. Instrument-specific event ordering and gap recovery.
- Don't fan out every tick to every strategy or browser. Subscribe only to required instruments; aggregate UI data at slower rates.
- Separate large history queries from latency-sensitive execution workers.
- Reduce dependencies, optimize DB indices with query plans, outbox and retention/partitioning, cache only rebuildable data.
- Compare costs and p99 under the actual Oracle VPS CPU and 1GB/limited RAM; provide headroom targets before live trading.

## Frontend
- Route-based code splitting, virtualization for thousands of instruments, memoized row rendering, efficient state selectors and throttled visual updates.
- Responsive flex/grid with correct min-width=0, controlled overflow, clear alignment; accessible focus states.
- Preload only critical assets, avoid heavy chart instances unless visible, show stale states when disconnected.

## Benchmarks (proposed targets, *not measured*)
- Dashboard input interaction <100ms in normal conditions; initial usable paint <2s over a realistic network on a midrange device.
- Core in-process strategy validation p99 <1ms in benchmark; exclude network/DB.
- Zero missing order/fill events under supported workloads; bounded resource usage with tested overload behavior.
- Define measured workload and data integrity before setting a guarantee for external exchange latency.

## Evidence required
Benchmark hardware, dataset and number of symbols/strategies, repetitions, version/commit, median/p95/p99, heap/RSS/CPU, queue lag, rejects/retries and failure injection outcomes. Regression thresholds in CI.
