# Control UI foundation

Responsive React + TypeScript operations dashboard. **Prototype only.** No simulated balances, live market prices, exchange keys or order placement.

## Run
```bash
cd apps/control-ui
npm install
npm run build
npm run dev
```
Use http://127.0.0.1:5173. In subsequent milestones add typed API client, authentication, real market catalog, virtualized tables and real-time connection/status indicators.

## Visual QA contract
Verify 360px, 375px, 768px, 1024px, 1280px and 1920px widths; 200% zoom; keyboard-only and reduced-motion behavior; table scroll containment; no guessed prices. Design specification: `../../docs/UI_UX.md`.
