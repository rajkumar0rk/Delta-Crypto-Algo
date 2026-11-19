# UI / UX design contract

## Layout
- 8px spacing system with 4px micro-alignment; consistent control heights (40px), min 44px touch targets on mobile.
- Desktop 256px navigation, single persistent topbar, scrollable content, card grid, full-width data tables with sticky headers.
- Tablet compact navigation; mobile stacked flow with an explicit menu, never hidden mandatory risk status or inaccessible emergency actions.
- Support 360/375/768/1024/1280/1440/1920 widths, 200% zoom, dark/light themes, keyboard tab, visible focus, screen readers and reduced motion.
- Define tokens for typography, spacing, elevation, borders, neutral surfaces, semantic alerts, gains/losses and status. Never convey critical state by color alone.
- Numeric columns right-aligned using tabular figures. Currency/units, source timestamps, stale-data ages and environment are always visible.
- High-density trading tables support sorting, filtering, virtualization, empty/loading/error/stale states, copy/export and user column preferences.
- No fake account balance, fabricated strategy return, or implied live connectivity in prototype. Every control must work or clearly say unavailable.

## Main application areas
Overview; Markets (all supported instruments, filters for crypto/RWA/perp/future/call/put); Strategies; Orders; Positions; Risk; Backtests; Reports; Audit; Settings; Operations. Watchlists, alerts, account selectors, multi-tenant RBAC, environment banners, order detail and incident controls.

## Design QA
Visual screenshot regression and accessibility checks at multiple widths and zooms. Performance budgets for interaction latency, websocket fan-out, table row count, chart recycling, memory, reduced motion and network degradation.
