import StrategyWorkspace from "./strategy/StrategyWorkspace";
import { useEffect, useRef, useState } from "react";
import {
  Activity, AlertTriangle, BarChart3, BookOpenCheck, Boxes, ChevronDown,
  ChevronRight, ClipboardList, Database, FileText, Gauge, LayoutDashboard,
  LockKeyhole, Menu, Search, Settings2, ShieldCheck, SlidersHorizontal,
  Terminal, Wallet, X
} from "lucide-react";
import type { LucideIcon } from "lucide-react";

type Page = "Overview" | "Markets" | "Strategies" | "Orders" | "Positions" | "Risk Center" | "Backtests" | "Reports" | "Audit Log" | "Settings" | "Operations";
type Group = { title: string; items: { name: Page; icon: LucideIcon }[] };
const groups: Group[] = [
  { title: "WORKSPACE", items: [
    { name: "Overview", icon: LayoutDashboard }, { name: "Markets", icon: BarChart3 },
    { name: "Strategies", icon: Boxes }, { name: "Backtests", icon: Gauge }
  ]},
  { title: "TRADING", items: [
    { name: "Orders", icon: ClipboardList }, { name: "Positions", icon: Wallet },
    { name: "Risk Center", icon: ShieldCheck }
  ]},
  { title: "OPERATIONS", items: [
    { name: "Reports", icon: FileText }, { name: "Audit Log", icon: BookOpenCheck },
    { name: "Operations", icon: Activity }, { name: "Settings", icon: Settings2 }
  ]}
];
const marketTypes = ["All instruments", "Crypto", "RWA-linked", "Perpetual", "Futures", "Calls", "Puts"] as const;
const routeDescriptions: Record<Page, string> = {
  Overview: "Your operations dashboard, risk visibility, and system health in one place.",
  Markets: "Discover every eligible contract dynamically, including crypto and RWA-linked derivatives.",
  Strategies: "Build, test, version, and supervise your strategies.",
  Orders: "Inspect execution activity and reconciliation status.",
  Positions: "Monitor open exposure, margin, and marked positions.",
  "Risk Center": "Protect accounts with strict limits, approvals, and emergency controls.",
  Backtests: "Run deterministic research and validate paper/live parity.",
  Reports: "Explore reconciled P&L, fees, funding and performance.",
  "Audit Log": "Review traceable operator and trading system events.",
  Operations: "Observe data feeds, latency, errors, and system recovery.",
  Settings: "Manage environment settings and authorized connections."
};

function PanelEmpty({ title, description, icon: Icon }: { title: string; description: string; icon: LucideIcon }) {
  return (
    <div className="empty-state">
      <div className="empty-icon"><Icon size={24} strokeWidth={1.7} aria-hidden="true" /></div>
      <h3>{title}</h3>
      <p>{description}</p>
      <span className="label-chip">Planned milestone</span>
    </div>
  );
}
function Placeholder({ label, note }: { label: string; note: string }) {
  return <div className="metric-card"><span>{label}</span><strong className="metric-pending">—</strong><small>{note}</small></div>;
}
export default function App() {
  const [page, setPage] = useState<Page>("Overview");
  const [menuOpen, setMenuOpen] = useState(false);
  const [marketType, setMarketType] = useState<(typeof marketTypes)[number]>("All instruments");
  const [marketSearch, setMarketSearch] = useState("");
  const menuButtonRef = useRef<HTMLButtonElement>(null);
  const closeMenuButtonRef = useRef<HTMLButtonElement>(null);
  const mobileNavRef = useRef<HTMLElement>(null);
  useEffect(() => {
    if (!menuOpen) return;
    // Focus after the mobile drawer becomes visible and its slide transition ends.
    // Moving focus while the element is still translated off-screen is unreliable.
    const focusHandle = window.setTimeout(() => closeMenuButtonRef.current?.focus(), 220);
    function handleNavigationKeys(event: KeyboardEvent) {
      if (event.key === "Escape") {
        event.preventDefault();
        setMenuOpen(false);
        menuButtonRef.current?.focus();
      }
      if (event.key !== "Tab") return;
      const buttons = Array.from(mobileNavRef.current?.querySelectorAll<HTMLButtonElement>("button:not(:disabled)") ?? [])
        .filter(button => getComputedStyle(button).display !== "none");
      const first = buttons[0];
      const last = buttons[buttons.length - 1];
      if (!first || !last) return;
      // A focusable element outside the drawer must never become a keyboard
      // escape route while the mobile navigation is acting as an overlay.
      if (!mobileNavRef.current?.contains(document.activeElement)) {
        event.preventDefault();
        (event.shiftKey ? last : first).focus();
      } else if (event.shiftKey && document.activeElement === first) {
        event.preventDefault();
        last.focus();
      } else if (!event.shiftKey && document.activeElement === last) {
        event.preventDefault();
        first.focus();
      }
    }
    document.addEventListener("keydown", handleNavigationKeys);
    return () => {
      window.clearTimeout(focusHandle);
      document.removeEventListener("keydown", handleNavigationKeys);
    };
  }, [menuOpen]);
  useEffect(() => {
    // A viewport resize or zoom-induced breakpoint change must not leave
    // a stale mobile backdrop or drawer state over the desktop navigation.
    const desktop = window.matchMedia("(min-width: 981px)");
    function handleBreakpointChange() {
      if (desktop.matches) {
        setMenuOpen(false);
        if (document.activeElement === closeMenuButtonRef.current) {
          mobileNavRef.current?.querySelector<HTMLButtonElement>(".nav-link.selected")?.focus();
        }
      } else if (mobileNavRef.current?.contains(document.activeElement)) {
        // Moving from the desktop sidebar to the hidden mobile drawer must
        // not strand keyboard focus inside invisible navigation controls.
        setMenuOpen(false);
        menuButtonRef.current?.focus();
      }
    }
    desktop.addEventListener("change", handleBreakpointChange);
    return () => desktop.removeEventListener("change", handleBreakpointChange);
  }, []);
  function closeMenu() {
    setMenuOpen(false);
    menuButtonRef.current?.focus();
  }
  useEffect(() => {
    document.title = `${page} · Delta Crypto Algo`;
  }, [page]);
  function navigate(next: Page) {
    setPage(next);
    setMenuOpen(false);
    window.scrollTo({ top: 0, behavior: "instant" });
    window.requestAnimationFrame(() => document.getElementById("main")?.focus());
  }
  return (
    <div className="app">
      <a className="skip-link" href="#main">Skip to main content</a>
      {menuOpen && <button aria-label="Close navigation" aria-hidden="true" tabIndex={-1} className="mobile-scrim" onClick={closeMenu} />}
      <aside className={`sidebar ${menuOpen ? "sidebar-open" : ""}`} ref={mobileNavRef} id="primary-navigation" aria-label="Primary navigation">
        <div className="brand">
          <div className="brand-mark" aria-hidden="true">D</div>
          <div className="brand-word"><strong>delta<span>crypto</span></strong><small>ALGO PLATFORM</small></div>
          <button ref={closeMenuButtonRef} className="close-nav" onClick={closeMenu} aria-label="Close menu"><X size={19} /></button>
        </div>
        <div className="sidebar-scroll">
          <div className="workspace-tile">
            <div className="workspace-orbit"><Terminal size={17} /></div>
            <div><strong>Workspace</strong><small>Foundation environment</small></div>
            <ChevronDown size={15} className="muted" aria-hidden="true" />
          </div>
          <nav>
            {groups.map(group => (
              <div className="nav-group" key={group.title}>
                <div className="group-label">{group.title}</div>
                {group.items.map(({name, icon: Icon}) => (
                  <button key={name} onClick={() => navigate(name)} aria-current={page === name ? "page" : undefined} className={`nav-link ${page === name ? "selected" : ""}`}>
                    <Icon size={18} strokeWidth={1.8} aria-hidden="true" /><span>{name}</span>{page === name && <span className="active-pill" />}
                  </button>
                ))}
              </div>
            ))}
          </nav>
        </div>
        <div className="sidebar-bottom">
          <span className="safety-tag"><LockKeyhole size={14} /> LIVE TRADING LOCKED</span>
          <small>Foundation v0.1 · No accounts connected</small>
        </div>
      </aside>
      <div className="workspace">
        <header className="topbar">
          <div className="top-left">
            <button ref={menuButtonRef} className="hamburger" aria-controls="primary-navigation" aria-expanded={menuOpen} aria-label="Open navigation" onClick={() => setMenuOpen(true)}><Menu size={21} /></button>
            <span className="top-location">Control Center</span><ChevronRight size={15} className="muted" />
            <strong>{page}</strong>
          </div>
          <div className="top-actions">
            <div className="session-status"><span className="status-dot amber" /> Offline prototype</div>
            <div className="env-tag"><LockKeyhole size={14} /> PAPER ONLY</div>
            <div className="profile-chip" title="User authentication not connected">U</div>
          </div>
        </header>
        <main id="main" className="content" tabIndex={-1}>
          <div className="page-heading">
            <div><span className="eyebrow">DELTA CRYPTO ALGO / CONTROL CENTER</span>
              <h1>{page === "Overview" ? "Trading overview" : page}</h1>
              <p>{routeDescriptions[page]}</p>
            </div>
            <span className="foundation-badge"><AlertTriangle size={15} /> Foundation build</span>
          </div>
          <div className="warning-strip" role="status">
            <div className="warning-icon"><LockKeyhole size={18} /></div>
            <div><strong>Live trading is disabled.</strong><span>This is a UI foundation prototype. No account, market stream, trading engine or order gateway is connected.</span></div>
          </div>
          {page === "Overview" ? (
            <>
              <section aria-labelledby="health-heading" className="overview-section">
                <div className="section-heading"><h2 id="health-heading">Account snapshot</h2><span>Awaiting authenticated data</span></div>
                <div className="metric-grid">
                  <Placeholder label="Net equity" note="No account connected" />
                  <Placeholder label="Unrealized P&L" note="No positions loaded" />
                  <Placeholder label="Margin utilization" note="No balances loaded" />
                  <Placeholder label="Active strategies" note="Runtime not deployed" />
                </div>
              </section>
              <div className="dashboard-grid">
                <section className="surface" aria-label="Portfolio overview">
                  <div className="surface-heading"><div><h2>Portfolio exposure</h2><p>Across supported crypto and RWA derivatives</p></div><span className="soft-chip">Not connected</span></div>
                  <PanelEmpty title="Portfolio data unavailable" description="Connect a verified account in a later milestone to see exposure, P&L, and asset-level risk here." icon={Wallet} />
                </section>
                <section className="surface" aria-label="System readiness">
                  <div className="surface-heading"><div><h2>System readiness</h2><p>Operational components and safety gates</p></div></div>
                  <div className="system-list">
                    {[
                      ["Control UI", "Prototype", "amber"],
                      ["Trading engine", "Not implemented", "neutral"],
                      ["Exchange streams", "Disconnected", "neutral"],
                      ["Order gateway", "Disabled", "neutral"],
                      ["Risk gateway", "Not implemented", "neutral"]
                    ].map(([label, status, tone]) => <div className="system-row" key={label}><span><span className={`status-dot ${tone}`} />{label}</span><small>{status}</small></div>)}
                  </div>
                  <div className="system-footer"><ShieldCheck size={16} /> No outbound orders can be placed in this build</div>
                </section>
              </div>
              <section className="surface quick-panel">
                <div className="surface-heading"><div><h2>Platform workspaces</h2><p>Navigate major enterprise capabilities</p></div></div>
                <div className="workspace-cards">
                  {([{name:"Markets", icon:Database, desc:"All eligible instruments"},{name:"Strategies",icon:SlidersHorizontal,desc:"Configurations & versions"},{name:"Risk Center",icon:ShieldCheck,desc:"Limits & safety gates"},{name:"Operations",icon:Activity,desc:"Health & diagnostics"}] as {name:Page;icon:LucideIcon;desc:string}[]).map(({name,icon:Icon,desc})=>(
                    <button className="workspace-card" onClick={() => navigate(name)} key={name}><Icon size={19}/><strong>{name}</strong><small>{desc}</small><ChevronRight size={16} className="card-arrow"/></button>
                  ))}
                </div>
              </section>
            </>
          ) : page === "Markets" ? (
            <section className="surface markets-surface">
              <div className="surface-heading"><div><h2>Instrument explorer</h2><p>Dynamic exchange catalog — integration scheduled for M2</p></div><span className="soft-chip">0 loaded</span></div>
              <div className="market-toolbar">
                <label className="search-box"><Search size={18} aria-hidden="true" /><span className="sr-only">Search instruments</span><input placeholder="Search symbol or underlying…" value={marketSearch} onChange={e => setMarketSearch(e.target.value)} /></label>
                <div className="type-filters" role="group" aria-label="Instrument type">
                  {marketTypes.map(type => <button key={type} aria-pressed={marketType===type} className={marketType===type?"active-filter":""} onClick={()=>setMarketType(type)}>{type}</button>)}
                </div>
              </div>
              <div className="table-scroller"><table><thead><tr><th>Instrument</th><th>Market</th><th>Contract</th><th>Expiry</th><th>Trading status</th><th className="align-right">Price</th></tr></thead><tbody><tr><td className="no-rows" colSpan={6}>No instruments loaded. Market catalog service is not connected yet.</td></tr></tbody></table></div>
              <div className="table-footnote"><Database size={14} /> Search and filters will apply to the complete Delta instrument catalog once connected.</div>
            </section>
          ) : page === "Strategies" ? (
            <StrategyWorkspace />
          ) : (
            <section className="surface route-placeholder">
              <div className="surface-heading"><div><h2>{page} workspace</h2><p>Full implementation scheduled in the project roadmap</p></div></div>
              <PanelEmpty icon={page==="Risk Center"?ShieldCheck:page==="Operations"?Activity:page==="Settings"?Settings2:Boxes} title={`${page} is not connected`} description="This module is intentionally inactive during the foundation milestone. Its backend, access controls and tests must be completed before activation." />
            </section>
          )}
          <footer className="footer"><span>© Delta Crypto Algo · Independent trading software</span><span>Foundation v0.1 <span className="footer-separator">·</span> <strong>Real-money orders unavailable</strong></span></footer>
        </main>
      </div>
    </div>
  );
}
