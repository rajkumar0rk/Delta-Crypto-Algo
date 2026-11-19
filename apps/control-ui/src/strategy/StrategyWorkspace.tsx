import { useMemo, useState } from "react";
import { Activity, ArrowRight, BarChart3, Boxes, ChevronRight, FileJson, GitBranch, LockKeyhole, Plus, Search, ShieldCheck, } from "lucide-react";
import StrategyBuilder from "./StrategyBuilder";
import { createDraft, makeComparison } from "./model";
import type { StrategyDefinition } from "./model";
import "./StrategyWorkspace.css";

type Preset = "blank" | "mtf" | "options" | "breakout";
type TemplateInfo = {
  id: Preset;
  icon: typeof Boxes;
  title: string;
  subtitle: string;
  symbol: string;
  timeframes: string;
  category: string;
};
const templates: TemplateInfo[] = [
  { id: "blank", icon: Plus, title: "Create from scratch", subtitle: "Start with a clean, configurable strategy draft", symbol: "All eligible assets", timeframes: "Your choice", category: "Blank" },
  { id: "mtf", icon: GitBranch, title: "Multi-timeframe trend", subtitle: "Slow trend filter + fast confirmation + separate protective monitoring", symbol: "Crypto futures", timeframes: "1H / 15M / 5M", category: "Futures" },
  { id: "options", icon: Activity, title: "Options oscillator", subtitle: "Call/put signals with explicit contract actions and premium-based stop", symbol: "Listed crypto options", timeframes: "1H / 1M", category: "Options" },
  { id: "breakout", icon: BarChart3, title: "Price breakout", subtitle: "Current confirmed close versus an earlier candle high", symbol: "Crypto / RWA futures", timeframes: "15M / 1M", category: "Futures" }
];
const onboarding = [
  ["01", "Choose instruments", "Dynamic listings, futures or options, DTE and strike rules"],
  ["02", "Define your signals", "Indicators, nested AND/OR rules, long and short setups"],
  ["03", "Set execution and protection", "Sizing, limit/market preferences, stops and hard risk limits"],
  ["04", "Validate and review", "Inspect errors and export a JSON draft"],
  ["05", "Backtest and paper trade", "Planned — unavailable until the simulation engine is verified"],
  ["06", "Controlled deployment", "Future milestone — requires risk review and manual approval"]
];
function preset(id: Preset): StrategyDefinition {
  const draft = createDraft();
  if (id === "blank") return draft;
  if (id === "mtf") {
    draft.id = "trend_pullback_draft";
    draft.metadata.name = "Multi-timeframe trend pullback";
    draft.metadata.description = "Illustrative configuration, not evaluated or profitable by default.";
    draft.clock.signal_timeframe = "15m";
    draft.clock.execution_timeframe = "5m";
    draft.signals.long_entry.children = [
      { ...makeComparison("1h"), id: "trend_filter", left: { kind: "indicator", indicator: "ema", parameters: { period: 50 }, output: "value", timeframe: "1h", offset: 0 },
        operator: "gt", right: { kind: "indicator", indicator: "ema", parameters: { period: 200 }, output: "value", timeframe: "1h", offset: 0 } },
      { ...makeComparison("15m"), id: "ema_touch", left: { kind: "price", field: "low", timeframe: "15m", offset: 0 },
        operator: "lte", right: { kind: "indicator", indicator: "ema", parameters: { period: 20 }, output: "value", timeframe: "15m", offset: 0 } }
    ];
    draft.risk.stop_loss.mode = "atr_multiple";
    draft.risk.stop_loss.value = "2";
    return draft;
  }
  if (id === "options") {
    draft.id = "options_oscillator_draft";
    draft.metadata.name = "Options oscillator";
    draft.metadata.description = "Research-only CALL/PUT contract selection with a premium-reference stop.";
    draft.universe.product_types = ["call_option", "put_option"];
    draft.clock.signal_timeframe = "1h";
    draft.clock.execution_timeframe = "1m";
    draft.signals.long_entry.children = [
      { type: "compare", id: "options_oversold", left: { kind: "indicator", indicator: "stoch_rsi", parameters: { period: 14 }, output: "k", timeframe: "1h", offset: 0 },
        operator: "lt", right: { kind: "constant", value: "20" } }
    ];
    draft.signals.short_entry.children = [
      { type: "compare", id: "options_overbought", left: { kind: "indicator", indicator: "stoch_rsi", parameters: { period: 14 }, output: "k", timeframe: "1h", offset: 0 },
        operator: "gt", right: { kind: "constant", value: "80" } }
    ];
    draft.risk.stop_loss = { mode: "percent", value: "15", basis: "option_entry_premium" };
    return draft;
  }
  draft.id = "confirmed_breakout_draft";
  draft.metadata.name = "Confirmed candle breakout";
  draft.metadata.description = "Price action breakout research template with fully confirmed candle inputs.";
  draft.clock.signal_timeframe = "15m";
  draft.signals.long_entry.children = [
    { type: "compare", id: "confirmed_breakout", left: { kind: "price", field: "close", timeframe: "15m", offset: 0 },
      operator: "gt", right: { kind: "price", field: "high", timeframe: "15m", offset: 1 } }
  ];
  return draft;
}
export default function StrategyWorkspace() {
  const [activeDraft, setActiveDraft] = useState<StrategyDefinition | null>(null);
  const [query, setQuery] = useState("");
  const filtered = useMemo(() => templates.filter(item =>
    [item.title, item.subtitle, item.symbol, item.category].some(s => s.toLowerCase().includes(query.trim().toLowerCase()))
  ), [query]);
  function openEditor(id: Preset) { setActiveDraft(preset(id)); }
  function goBack() {
    if (window.confirm("Leave Configure Strategy? This draft is not saved. Export it first if you want to keep it.")) setActiveDraft(null);
  }
  if (activeDraft) return <div className="sw-full-editor">
    <div className="sw-backline"><button onClick={goBack} aria-label="Return to strategy workspace">Strategies <ChevronRight size={14}/> <strong>Configure strategy</strong></button></div>
    <StrategyBuilder initialDraft={activeDraft} onBack={goBack}/>
  </div>;
  return <div className="sw-root">
    <div className="sw-heading">
      <div><div className="sw-eyebrow">STRATEGY WORKSPACE · OFFLINE DESIGN</div>
        <h2>Build and configure strategies</h2>
        <p>Professional rules without coding. Choose instruments, entries, exits, options, order sizing, and protection—all in one guided editor.</p></div>
      <button className="sw-primary" onClick={() => openEditor("blank")}><Plus size={17}/> Configure Strategy</button>
    </div>
    <div className="sw-banner" role="status"><LockKeyhole size={18}/><div><strong>Trading is disabled in this development build.</strong>
      <span>This workspace creates local, downloadable JSON drafts only. No strategy is deployed, and no account data is connected.</span></div></div>
    <div className="sw-summary">
      <div><span>Saved strategies</span><strong>0</strong><small>Backend persistence not implemented</small></div>
      <div><span>Live strategies</span><strong>0</strong><small>Live mode locked</small></div>
      <div><span>Included templates</span><strong>{templates.length}</strong><small>Examples, not tested trading edges</small></div>
      <div><span>Draft format</span><strong>v1</strong><small>Versioned JSON schema</small></div>
    </div>
    <div className="sw-section-heading"><div><h3>Strategy templates</h3><p>Use a starting point or begin from scratch. Every field remains editable.</p></div>
      <label className="sw-search"><Search size={16}/><span className="sw-sr-only">Search templates</span>
        <input value={query} onChange={e => setQuery(e.target.value)} placeholder="Search templates..." /></label></div>
    <div className="sw-templates">{filtered.map(item => { const Icon=item.icon; return <button key={item.id} className="sw-template" onClick={() => openEditor(item.id)}>
      <div className="sw-template-header"><span className="sw-template-icon"><Icon size={21}/></span><span className="sw-template-cat">{item.category}</span></div>
      <strong>{item.title}</strong><p>{item.subtitle}</p>
      <div className="sw-template-meta"><span>{item.symbol}</span><span>{item.timeframes}</span></div>
      <div className="sw-template-action">Configure in editor <ArrowRight size={14}/></div>
    </button>;})}
      {!filtered.length && <div className="sw-no-results">No templates match your search. Try another term, or choose Configure Strategy.</div>}
    </div>
    <section className="sw-process"><div className="sw-section-heading"><div><h3>Professional strategy workflow</h3>
      <p>Design is only the first step. Validation, backtesting, and release approval are separate gates.</p></div></div>
      <div className="sw-process-grid">{onboarding.map(([index,title,detail],i)=><div key={index} className="sw-process-item">
        <span>{index}</span><strong>{title}</strong><small>{detail}</small>{i===4||i===5?<em>Planned</em>:null}
      </div>)}</div>
    </section>
    <div className="sw-footer-note"><ShieldCheck size={16}/> All strategies will ultimately pass a separate global risk gateway—editor settings can never override it.
      <span><FileJson size={15}/> JSON import/export available inside Configure Strategy.</span></div>
  </div>;
}
