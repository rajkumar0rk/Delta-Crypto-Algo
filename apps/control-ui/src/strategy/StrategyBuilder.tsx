import { useMemo, useRef, useState } from "react";
import type { ChangeEvent, ReactNode } from "react";
import { AlertCircle, ArrowLeft, ArrowRight, Braces, CheckCircle2, Clock3, Download, FileJson, GitBranch, Layers3, LockKeyhole, Plus, ShieldCheck, SlidersHorizontal, Trash2, Upload, X } from "lucide-react";
import {
  INDICATORS, INDICATOR_OUTPUTS, OPERATORS, PRICE_FIELDS, PRODUCT_TYPES, TIMEFRAMES, WEEKDAYS,
  createDraft, editNode, makeComparison, makeGroup, makeOperand
} from "./model";
import type { Comparison, Operand, RuleGroup, RuleNode, SignalSide, StrategyDefinition, Timeframe } from "./model";
import { validateStrategy } from "./validation";
import "./StrategyBuilder.css";

type Step = "identity" | "universe" | "time" | "entry" | "exit" | "execution" | "risk" | "review";
const steps: { id: Step; label: string; subtitle: string }[] = [
  { id: "identity", label: "Identity", subtitle: "Name and version" },
  { id: "universe", label: "Instrument universe", subtitle: "All supported markets" },
  { id: "time", label: "Timeframes & sessions", subtitle: "Timing and schedule" },
  { id: "entry", label: "Entry conditions", subtitle: "Long and short rules" },
  { id: "exit", label: "Exit conditions", subtitle: "Long and short rules" },
  { id: "execution", label: "Execution & sizing", subtitle: "Orders and contracts" },
  { id: "risk", label: "Risk controls", subtitle: "Protective limits" },
  { id: "review", label: "Review & export", subtitle: "Validate your draft" }
];
const choices = {
  operator: { gt: ">", gte: "≥", lt: "<", lte: "≤", eq: "=", neq: "≠", crosses_above: "Crosses above", crosses_below: "Crosses below" } as Record<string,string>,
  kind: { constant: "Number", price: "Price / market field", indicator: "Indicator" } as Record<string,string>,
  protection: { off: "Disabled", percent: "Percentage", price_ticks: "Price ticks", atr_multiple: "ATR multiple", signal_candle: "Signal candle high / low" } as Record<string,string>
};
function Field({ label, hint, children }: { label: string; hint?: string; children: ReactNode }) {
  return <label className="sb-field"><span className="sb-field-name">{label}</span>{children}{hint && <small>{hint}</small>}</label>;
}
function SelectField({ label, value, options, onChange, hint }: {
  label: string; value: string; options: { value: string; label: string }[]; onChange: (value: string) => void; hint?: string
}) {
  return <Field label={label} hint={hint}><select value={value} onChange={event => onChange(event.target.value)}>
    {options.map(option => <option key={option.value} value={option.value}>{option.label}</option>)}
  </select></Field>;
}
function choicesFrom(values: readonly string[]) {
  return values.map(value => ({ value, label: value.replaceAll("_", " ") }));
}
function NumberField({ label, value, onChange, min, max, step, hint }: {
  label: string; value: number; onChange: (next: number) => void; min?: number; max?: number; step?: number; hint?: string
}) {
  return <Field label={label} hint={hint}><input type="number" value={value} min={min} max={max} step={step ?? 1}
    onChange={event => onChange(event.target.value === "" ? 0 : Number(event.target.value))}/></Field>;
}
function DecimalField({ label, value, onChange, hint }: {
  label: string; value: string; onChange: (next: string) => void; hint?: string
}) {
  return <Field label={label} hint={hint}><input type="text" inputMode="decimal" value={value} onChange={event => onChange(event.target.value)} placeholder="0"/></Field>;
}
function CheckRow({ label, checked, onChange, hint }: {
  label: string; checked: boolean; onChange: (value: boolean) => void; hint?: string
}) {
  return <label className="sb-check"><input type="checkbox" checked={checked} onChange={event => onChange(event.target.checked)}/>
    <span><strong>{label}</strong>{hint && <small>{hint}</small>}</span></label>;
}
function OperandEditor({ label, value, onChange, defaultTimeframe }: {
  label: string; value: Operand; onChange: (v: Operand) => void; defaultTimeframe: Timeframe
}) {
  const detail = value.kind === "indicator" ? INDICATORS.find(x => x.id === value.indicator) : undefined;
  return <div className="sb-operand">
    <div className="sb-operand-heading">{label}</div>
    <div className="sb-operand-row">
      <select aria-label={label + " type"} value={value.kind}
        onChange={event => onChange(makeOperand(event.target.value as Operand["kind"], defaultTimeframe))}>
        {Object.entries(choices.kind).map(([key, text]) => <option key={key} value={key}>{text}</option>)}
      </select>
      {value.kind === "constant" ? <input aria-label={label + " value"} inputMode="decimal" placeholder="Numeric value"
        value={value.value} onChange={event => onChange({ ...value, value: event.target.value })}/> : (
        <>
          {value.kind === "price" ? <select aria-label={label + " price field"} value={value.field}
            onChange={event => onChange({ ...value, field: event.target.value })}>
            {PRICE_FIELDS.map(field => <option key={field} value={field}>{field}</option>)}
          </select> : <select aria-label={label + " indicator"} value={value.indicator}
            onChange={event => {
              const next = INDICATORS.find(x => x.id === event.target.value);
              if (next) onChange({ ...value, indicator: next.id, output: next.output, parameters: next.period ? { period: next.period } : {} });
            }}>{INDICATORS.map(indicator => <option key={indicator.id} value={indicator.id}>{indicator.label}</option>)}</select>}
          <select className="sb-tf-select" aria-label={label + " timeframe"} value={value.timeframe}
            onChange={event => onChange({ ...value, timeframe: event.target.value as Timeframe })}>
            {TIMEFRAMES.map(tf => <option key={tf} value={tf}>{tf}</option>)}
          </select>
          <input className="sb-offset-input" type="number" min="0" max="500" aria-label={label + " previous candles"} title="0 = current confirmed candle, 1 = previous confirmed candle"
            value={value.offset} onChange={event => onChange({ ...value, offset: Number(event.target.value) })}/>
        </>
      )}
    </div>
    {value.kind === "indicator" && <div className="sb-operand-row sb-operand-secondary">
      <label>Output<select value={value.output} onChange={event => onChange({ ...value, output: event.target.value })}>
        {(INDICATOR_OUTPUTS[value.indicator] ?? [value.output]).map(output => <option key={output} value={output}>{output}</option>)}
      </select></label>
      {detail && detail.period > 0 && <label>Period<input type="number" min="1" max="10000"
        value={Number(value.parameters.period ?? detail.period)} onChange={event => onChange({
          ...value, parameters: { ...value.parameters, period: Number(event.target.value) }
        })}/></label>}
      {value.indicator === "supertrend" && <label>Factor<input type="number" min="0.1" max="100" step="0.1"
        value={Number(value.parameters.factor ?? 3)} onChange={event => onChange({
          ...value, parameters: { ...value.parameters, factor: Number(event.target.value) }
        })}/></label>}
    </div>}
  </div>;
}
function ConditionGroupEditor({ group, onChange, defaultTimeframe, depth = 0 }: {
  group: RuleGroup; onChange: (value: RuleGroup) => void; defaultTimeframe: Timeframe; depth?: number
}) {
  const patchNode = (id: string, modify: (node: RuleNode) => RuleNode | null) => onChange(editNode(group, id, modify));
  return <div className={"sb-rule-group " + (depth ? "sb-nested-group" : "")}>
    <div className="sb-group-top">
      <div className="sb-group-title"><GitBranch size={15}/><span>{depth ? "Nested rules" : "Conditions"}</span><span className="sb-counter">{group.children.length}</span></div>
      <label className="sb-join-label">Match <select aria-label="Match group rules" value={group.join} onChange={event =>
        onChange({ ...group, join: event.target.value as "all" | "any" })}>
        <option value="all">ALL (AND)</option><option value="any">ANY (OR)</option>
      </select></label>
    </div>
    {group.children.length === 0 && <div className="sb-rule-empty">No conditions. This signal stays disabled until a condition is added.</div>}
    <div className="sb-rule-stack">{group.children.map((node, index) => node.type === "group" ? (
      <div key={node.id} className="sb-child-group">
        <div className="sb-remove-row"><strong>Group {index + 1}</strong><button className="sb-icon-action" aria-label="Remove nested group"
          onClick={() => patchNode(node.id, () => null)}><Trash2 size={15}/></button></div>
        <ConditionGroupEditor group={node} depth={depth + 1} defaultTimeframe={defaultTimeframe}
          onChange={next => patchNode(node.id, () => next)}/>
      </div>
    ) : (
      <div className="sb-rule-card" key={node.id}>
        <div className="sb-remove-row"><strong>Condition {index + 1}</strong><button className="sb-icon-action" aria-label="Remove condition"
          onClick={() => patchNode(node.id, () => null)}><Trash2 size={15}/></button></div>
        <OperandEditor label="Compare" value={node.left} defaultTimeframe={defaultTimeframe}
          onChange={operand => patchNode(node.id, curr => ({ ...(curr as Comparison), left: operand }))}/>
        <div className="sb-operator">
          <label>Operator<select value={node.operator} onChange={event =>
            patchNode(node.id, curr => ({ ...(curr as Comparison), operator: event.target.value as Comparison["operator"] }))}>
            {OPERATORS.map(op => <option value={op} key={op}>{choices.operator[op]}</option>)}
          </select></label>
        </div>
        <OperandEditor label="Against" value={node.right} defaultTimeframe={defaultTimeframe}
          onChange={operand => patchNode(node.id, curr => ({ ...(curr as Comparison), right: operand }))}/>
      </div>
    ))}</div>
    <div className="sb-add-actions">
      <button type="button" onClick={() => onChange({ ...group, children: [...group.children, makeComparison(defaultTimeframe)] })}
        disabled={group.children.length >= 64}><Plus size={15}/> Add condition</button>
      <button type="button" onClick={() => onChange({ ...group, children: [...group.children, makeGroup("all")] })}
        disabled={depth >= 4 || group.children.length >= 64}><Layers3 size={15}/> Add AND/OR group</button>
    </div>
  </div>;
}

export default function StrategyBuilder() {
  const [draft, setDraft] = useState<StrategyDefinition>(() => createDraft());
  const [step, setStep] = useState<Step>("identity");
  const [notice, setNotice] = useState("");
  const [tagsInput, setTagsInput] = useState("draft");
  const [symbolsInput, setSymbolsInput] = useState("");
  const [idsInput, setIdsInput] = useState("");
  const fileInput = useRef<HTMLInputElement>(null);
  const issues = useMemo(() => validateStrategy(draft), [draft]);
  const serialized = useMemo(() => JSON.stringify(draft, null, 2), [draft]);
  const stepIndex = steps.findIndex(item => item.id === step);
  function identity(update: Partial<StrategyDefinition["metadata"]>) {
    setDraft(old => ({ ...old, metadata: { ...old.metadata, ...update } }));
  }
  function universe(update: Partial<StrategyDefinition["universe"]>) {
    setDraft(old => ({ ...old, universe: { ...old.universe, ...update } }));
  }
  function clock(update: Partial<StrategyDefinition["clock"]>) {
    setDraft(old => ({ ...old, clock: { ...old.clock, ...update } }));
  }
  function session(update: Partial<StrategyDefinition["session"]>) {
    setDraft(old => ({ ...old, session: { ...old.session, ...update } }));
  }
  function execution(update: Partial<StrategyDefinition["execution"]>) {
    setDraft(old => ({ ...old, execution: { ...old.execution, ...update } }));
  }
  function risk(update: Partial<StrategyDefinition["risk"]>) {
    setDraft(old => ({ ...old, risk: { ...old.risk, ...update } }));
  }
  function setSignal(side: SignalSide, group: RuleGroup) {
    setDraft(old => ({ ...old, signals: { ...old.signals, [side]: group } }));
  }
  function saveJSON() {
    const blob = new Blob([serialized], { type: "application/json" });
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = (draft.id.replace(/[^a-z0-9_-]/gi, "_") || "strategy") + "-draft.json";
    document.body.appendChild(link);
    link.click();
    link.remove();
    URL.revokeObjectURL(url);
    setNotice(issues.length ? "Draft downloaded with validation issues. It cannot be activated." : "Validated draft downloaded. This is not a live deployment.");
  }
  async function importJSON(event: ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0];
    event.target.value = "";
    if (!file) return;
    if (file.size > 512_000) {
      setNotice("Import rejected: maximum JSON file size is 500 KB.");
      return;
    }
    try {
      const content: unknown = JSON.parse(await file.text());
      const errors = validateStrategy(content);
      if (errors.length) {
        setNotice("Import rejected: " + errors[0].path + " — " + errors[0].message);
        return;
      }
      const loaded = content as StrategyDefinition;
      setDraft(loaded);
      setTagsInput(loaded.metadata.tags.join(", "));
      setSymbolsInput(loaded.universe.underlying_symbols.join(", "));
      setIdsInput(loaded.universe.product_ids.join(", "));
      setStep("review");
      setNotice("JSON imported as an offline draft. No account or trading engine was affected.");
    } catch {
      setNotice("Import rejected: invalid JSON.");
    }
  }
  async function copyJSON() {
    try {
      await navigator.clipboard.writeText(serialized);
      setNotice("Strategy JSON copied to clipboard.");
    } catch {
      setNotice("Clipboard access unavailable. Use Download JSON instead.");
    }
  }
  const hasOptions = draft.universe.product_types.some(type => type.endsWith("_option"));
  function editSignal(side: SignalSide, label: string) {
    return <section className="sb-signal-panel" key={side}>
      <div className="sb-signal-heading"><div><h3>{label}</h3><p>Evaluate using explicitly assigned timeframes. An empty group means disabled.</p></div><span className="sb-draft-chip">{side.includes("long") ? "LONG" : "SHORT"}</span></div>
      <ConditionGroupEditor group={draft.signals[side]} onChange={next => setSignal(side, next)} defaultTimeframe={draft.clock.signal_timeframe}/>
    </section>;
  }
  const section = (() => {
    switch(step) {
      case "identity": return <>
        <div className="sb-form-heading"><h2>Strategy identity</h2><p>Create a descriptive, versioned strategy. Editing this draft never modifies live systems.</p></div>
        <div className="sb-grid two">
          <Field label="Strategy name"><input maxLength={120} value={draft.metadata.name} onChange={event => identity({ name: event.target.value })}/></Field>
          <Field label="Strategy ID" hint="Stable lowercase ID used for revisions and exports"><input maxLength={64} value={draft.id} onChange={event =>
            setDraft(old => ({ ...old, id: event.target.value }))}/></Field>
          <NumberField label="Draft revision" value={draft.revision} min={1} onChange={value => setDraft(old => ({ ...old, revision: value }))}/>
          <Field label="Tags" hint="Separate tags with commas"><input value={tagsInput} onChange={event => {
            setTagsInput(event.target.value);
            identity({ tags: event.target.value.split(",").map(t => t.trim()).filter(Boolean) });
          }}/></Field>
        </div>
        <Field label="Description"><textarea value={draft.metadata.description} rows={5} maxLength={2000} placeholder="What market behavior should this strategy capture?"
          onChange={event => identity({ description: event.target.value })}/></Field>
        <div className="sb-info"><LockKeyhole size={16}/> This is a local, unsaved draft until you download its JSON. No live strategy is activated.</div>
      </>;
      case "universe": return <>
        <div className="sb-form-heading"><h2>Choose markets</h2><p>Select the instrument universe. Exact tradability will be checked against the exchange catalog and account limits in a later milestone.</p></div>
        <SelectField label="Product selection" value={draft.universe.selection} onChange={value =>
          universe({ selection: value as StrategyDefinition["universe"]["selection"] })} options={[{ value: "dynamic_filter", label: "Dynamic catalog filters" }, { value: "product_ids", label: "Exact product IDs" }]}/>
        <div className="sb-section-label">Eligible contract types</div>
        <div className="sb-checkbox-grid">{PRODUCT_TYPES.map(type => <CheckRow key={type} checked={draft.universe.product_types.includes(type)} label={type.replaceAll("_", " ")}
          onChange={checked => universe({ product_types: checked
            ? [...draft.universe.product_types, type]
            : draft.universe.product_types.filter(value => value !== type) })}/>)}</div>
        <div className="sb-grid two">
          {draft.universe.selection === "dynamic_filter" ? <Field label="Underlying symbols" hint="Optional comma-separated list; empty means all eligible symbols">
            <input value={symbolsInput} placeholder="All, or BTC, ETH, ..." onChange={event => {
              setSymbolsInput(event.target.value);
              universe({ underlying_symbols: event.target.value.split(",").map(v => v.trim().toUpperCase()).filter(Boolean) });
            }}/></Field> : <Field label="Delta product IDs" hint="Comma-separated positive numeric exchange IDs">
            <input value={idsInput} placeholder="12345, 98765" onChange={event => {
              setIdsInput(event.target.value);
              universe({ product_ids: event.target.value.split(",").map(v => v.trim()).filter(Boolean).map(Number) });
            }}/></Field>}
          <NumberField label="Maximum instruments" value={draft.universe.max_products} min={1} max={1000}
            hint="Bounded subscription and execution workload" onChange={value => universe({ max_products: value })}/>
        </div>
        <div className="sb-info"><Layers3 size={16}/> Covers crypto and RWA-linked products listed by Delta India. Product eligibility is not guaranteed by a filter match.</div>
      </>;
      case "time": return <>
        <div className="sb-form-heading"><h2>Timeframes and sessions</h2><p>Signal candles, individual indicator timeframes and protective monitoring operate independently.</p></div>
        <div className="sb-grid two">
          <SelectField label="Default signal timeframe" value={draft.clock.signal_timeframe} options={choicesFrom(TIMEFRAMES)}
            onChange={value => clock({ signal_timeframe: value as Timeframe })}/>
          <SelectField label="Protective execution timeframe" value={draft.clock.execution_timeframe} options={choicesFrom(TIMEFRAMES)}
            onChange={value => clock({ execution_timeframe: value as Timeframe })}/>
          <NumberField label="Maximum market data age (ms)" value={draft.clock.max_price_age_ms} min={100} max={120000}
            onChange={value => clock({ max_price_age_ms: value })}/>
          <Field label="Session timezone" hint="IANA timezone; validated by the backend when execution is added"><input value={draft.session.timezone}
            onChange={event => session({ timezone: event.target.value })}/></Field>
        </div>
        <div className="sb-checkbox-grid"><CheckRow checked={draft.clock.confirm_on_close} label="Confirm signal on candle close" hint="Required to prevent lookahead"
          onChange={() => {}}/><CheckRow checked={draft.clock.intrabar_protective_exits} label="Intrabar protective exits"
          hint="Monitored separately from signal closes" onChange={value => clock({ intrabar_protective_exits: value })}/></div>
        <div className="sb-section-label">Trading weekdays</div>
        <div className="sb-day-row">{WEEKDAYS.map((day, index) => <label key={day} className={"sb-day " + (draft.session.weekdays.includes(index + 1) ? "on" : "")}>
          <input type="checkbox" checked={draft.session.weekdays.includes(index + 1)} onChange={event => session({
            weekdays: event.target.checked
              ? [...draft.session.weekdays, index + 1].sort()
              : draft.session.weekdays.filter(value => value !== index + 1)
          })}/>{day}</label>)}</div>
        <div className="sb-section-label sb-section-heading-row">Session windows
          <button onClick={() => session({ windows: [...draft.session.windows, { start: "09:00", end: "17:00" }] })}
            disabled={draft.session.windows.length >= 12}><Plus size={14}/> Add window</button>
        </div>
        <div className="sb-windows">{draft.session.windows.map((window, index) => <div className="sb-window" key={index}>
          <Field label="From"><input type="time" value={window.start} onChange={event =>
            session({ windows: draft.session.windows.map((v, i) => i === index ? { ...v, start: event.target.value } : v) })}/></Field>
          <Field label="To"><input type="time" value={window.end} onChange={event =>
            session({ windows: draft.session.windows.map((v, i) => i === index ? { ...v, end: event.target.value } : v) })}/></Field>
          <button aria-label="Remove session window" className="sb-icon-action" disabled={draft.session.windows.length < 2}
            onClick={() => session({ windows: draft.session.windows.filter((_, i) => i !== index) })}><Trash2 size={16}/></button>
        </div>)}</div>
        <div className="sb-info"><Clock3 size={16}/> Every rule operand sets its own timeframe and candle offset. Historical aggregation must use only completed candles.</div>
      </>;
      case "entry": return <>
        <div className="sb-form-heading"><h2>Entry rules</h2><p>Configure long and short signals separately. Nested groups allow ALL / ANY combinations without writing code.</p></div>
        {editSignal("long_entry", "Enter long")}
        {editSignal("short_entry", "Enter short")}
        <SelectField label="If both entries trigger at once" value={draft.signals.conflict_policy}
          options={[{value:"skip_both",label:"Skip both — safest"},{value:"prefer_long",label:"Prefer long"},{value:"prefer_short",label:"Prefer short"}]}
          onChange={value => setDraft(old => ({ ...old, signals: { ...old.signals, conflict_policy: value as StrategyDefinition["signals"]["conflict_policy"] } }))}/>
      </>;
      case "exit": return <>
        <div className="sb-form-heading"><h2>Discretionary exit rules</h2><p>These rules complement mandatory stop-loss, trailing and global risk controls. Exits are reduce-only.</p></div>
        {editSignal("long_exit", "Exit long")}
        {editSignal("short_exit", "Exit short")}
        <div className="sb-info"><ShieldCheck size={16}/> Protective stops and portfolio-level controls always apply, even when these exit groups are empty.</div>
      </>;
      case "execution": return <>
        <div className="sb-form-heading"><h2>Orders and position sizing</h2><p>Define order preferences and allocation. Exchange-specific support and actual margin are validated later.</p></div>
        <div className="sb-grid two">
          <SelectField label="Order type" value={draft.execution.order_type} options={choicesFrom(["limit","market"])}
            onChange={value => execution({ order_type: value as "limit" | "market", post_only: value === "market" ? false : draft.execution.post_only })}/>
          <SelectField label="Time in force" value={draft.execution.time_in_force} options={choicesFrom(["gtc","ioc","fok"])}
            onChange={value => execution({ time_in_force: value as "gtc" | "ioc" | "fok", post_only: value === "gtc" ? draft.execution.post_only : false })}/>
          <NumberField label="Maximum slippage (bps)" value={draft.execution.max_slippage_bps} min={0} max={1000}
            onChange={value => execution({ max_slippage_bps: value })}/>
          <SelectField label="Sizing method" value={draft.execution.sizing.mode} options={[
            {value:"fixed_lots",label:"Fixed lots"},{value:"fixed_quote_notional",label:"Quote currency notional"},{value:"risk_pct_equity",label:"Risk % of equity"}
          ]} onChange={value => execution({ sizing: { ...draft.execution.sizing, mode: value as StrategyDefinition["execution"]["sizing"]["mode"] } })}/>
          <DecimalField label="Order size / budget" value={draft.execution.sizing.value} hint="Decimal input; interpreted according to sizing method"
            onChange={value => execution({ sizing: { ...draft.execution.sizing, value } })}/>
        </div>
        <div className="sb-checkbox-grid">
          <CheckRow label="Post-only (maker preference)" checked={draft.execution.post_only} hint="Requires limit + GTC; may reject on crossing"
            onChange={value => execution({ post_only: value, order_type: value ? "limit" : draft.execution.order_type,
              time_in_force: value ? "gtc" : draft.execution.time_in_force })}/>
          <CheckRow label="Reduce-only exits" checked={draft.execution.reduce_only_exits} hint="Always required" onChange={() => {}}/>
        </div>
        {hasOptions && <div className="sb-subpanel"><div className="sb-form-heading"><h3>Option contract selection</h3><p>Applies to call/put options; no fixed strike or expiry assumptions.</p></div>
          <div className="sb-grid two">
            <SelectField label="Expiry selection" value={draft.execution.options.expiry_selection}
              options={choicesFrom(["nearest_eligible","fixed_expiry","furthest_eligible"])}
              onChange={value => execution({ options: { ...draft.execution.options, expiry_selection: value as StrategyDefinition["execution"]["options"]["expiry_selection"] } })}/>
            {draft.execution.options.expiry_selection === "fixed_expiry" && <Field label="Fixed expiry (UTC)" hint="ISO-8601 date-time"><input type="datetime-local"
              value={draft.execution.options.expiry_utc?.slice(0,16) ?? ""}
              onChange={event => execution({ options: { ...draft.execution.options, expiry_utc: event.target.value ? new Date(event.target.value + "Z").toISOString() : undefined } })}/></Field>}
            <NumberField label="Minimum DTE (hours)" value={draft.execution.options.min_dte_hours} min={0}
              onChange={value => execution({ options: { ...draft.execution.options, min_dte_hours: value } })}/>
            <NumberField label="Maximum DTE (hours)" value={draft.execution.options.max_dte_hours} min={1}
              onChange={value => execution({ options: { ...draft.execution.options, max_dte_hours: value } })}/>
            <SelectField label="Strike choice" value={draft.execution.options.strike_selection}
              options={choicesFrom(["atm","offset_strikes","target_delta"])}
              onChange={value => execution({ options: { ...draft.execution.options, strike_selection: value as StrategyDefinition["execution"]["options"]["strike_selection"] } })}/>
            {draft.execution.options.strike_selection === "offset_strikes" && <NumberField label="Strike offset" value={draft.execution.options.strike_offset} min={-100} max={100}
              onChange={value => execution({ options: { ...draft.execution.options, strike_offset: value } })}/>}
            {draft.execution.options.strike_selection === "target_delta" && <DecimalField label="Target absolute delta" value={draft.execution.options.target_delta ?? "0.5"}
              onChange={value => execution({ options: { ...draft.execution.options, target_delta: value } })}/>}
            <SelectField label="Call / put policy" value={draft.execution.options.option_side}
              options={choicesFrom(["follow_signal","call_only","put_only"])}
              onChange={value => execution({ options: { ...draft.execution.options, option_side: value as StrategyDefinition["execution"]["options"]["option_side"] } })}/>
          </div>
        </div>}
      </>;
      case "risk": return <>
        <div className="sb-form-heading"><h2>Risk controls</h2><p>Every limit here is a strategy ceiling. The future global risk gateway can only apply stricter limits.</p></div>
        <div className="sb-grid two">
          <NumberField label="Max open positions" value={draft.risk.max_open_positions} min={1} max={1000}
            onChange={value => risk({ max_open_positions: value })}/>
          <NumberField label="Max orders per minute" value={draft.risk.max_orders_per_minute} min={1} max={10000}
            onChange={value => risk({ max_orders_per_minute: value })}/>
          <NumberField label="Max daily loss (%)" value={draft.risk.max_daily_loss_pct} min={0} max={100} step={0.1}
            onChange={value => risk({ max_daily_loss_pct: value })}/>
          <NumberField label="Max drawdown (%)" value={draft.risk.max_drawdown_pct} min={0} max={100} step={0.1}
            onChange={value => risk({ max_drawdown_pct: value })}/>
          <NumberField label="Max consecutive losses" value={draft.risk.max_consecutive_losses} min={1} max={1000}
            onChange={value => risk({ max_consecutive_losses: value })}/>
          <DecimalField label="Max quote notional" value={draft.risk.max_notional_quote} hint="Quote-currency units; account exposure enforced separately"
            onChange={value => risk({ max_notional_quote: value })}/>
          <NumberField label="Minimum margin buffer (%)" value={draft.risk.min_margin_buffer_pct} min={0} max={100} step={0.1}
            onChange={value => risk({ min_margin_buffer_pct: value })}/>
        </div>
        <div className="sb-subpanel"><div className="sb-form-heading"><h3>Stop-loss and take-profit</h3><p>Percent basis and tick value depend on product metadata and must be validated before any trade.</p></div>
          <div className="sb-grid two">
            <SelectField label="Stop-loss mode" value={draft.risk.stop_loss.mode} options={Object.entries(choices.protection).map(([value,label]) => ({value,label}))}
              onChange={value => risk({ stop_loss: { ...draft.risk.stop_loss, mode: value as StrategyDefinition["risk"]["stop_loss"]["mode"] } })}/>
            <DecimalField label="Stop-loss value" value={draft.risk.stop_loss.value} hint="0 for disabled or signal-candle mode"
              onChange={value => risk({ stop_loss: { ...draft.risk.stop_loss, value } })}/>
            <SelectField label="Take-profit mode" value={draft.risk.take_profit.mode} options={Object.entries(choices.protection).map(([value,label]) => ({value,label}))}
              onChange={value => risk({ take_profit: { ...draft.risk.take_profit, mode: value as StrategyDefinition["risk"]["take_profit"]["mode"] } })}/>
            <DecimalField label="Take-profit value" value={draft.risk.take_profit.value}
              onChange={value => risk({ take_profit: { ...draft.risk.take_profit, value } })}/>
          </div>
        </div>
        <div className="sb-subpanel"><CheckRow label="Enable trailing protective stop" checked={draft.risk.trailing_stop.enabled}
          onChange={value => risk({ trailing_stop: { ...draft.risk.trailing_stop, enabled: value } })}/>
          {draft.risk.trailing_stop.enabled && <div className="sb-grid three">
            <SelectField label="Trail unit" value={draft.risk.trailing_stop.unit} options={choicesFrom(["percent","price_ticks","atr_multiple"])}
              onChange={value => risk({ trailing_stop: { ...draft.risk.trailing_stop, unit: value as StrategyDefinition["risk"]["trailing_stop"]["unit"] } })}/>
            <DecimalField label="Activation" value={draft.risk.trailing_stop.activation}
              onChange={value => risk({ trailing_stop: { ...draft.risk.trailing_stop, activation: value } })}/>
            <DecimalField label="Pullback distance" value={draft.risk.trailing_stop.distance}
              onChange={value => risk({ trailing_stop: { ...draft.risk.trailing_stop, distance: value } })}/>
          </div>}
        </div>
        <div className="sb-checkbox-grid">
          <CheckRow checked={draft.risk.halt_on_stale_data} label="Halt on stale data" hint="Mandatory, cannot disable" onChange={() => {}}/>
          <CheckRow checked={draft.risk.halt_on_reconciliation_gap} label="Halt on account reconciliation gap" hint="Mandatory, cannot disable" onChange={() => {}}/>
        </div>
      </>;
      case "review": return <>
        <div className="sb-form-heading"><h2>Review your strategy</h2><p>Inspect validation messages and export this strategy as a versioned JSON draft. No strategy runner is connected.</p></div>
        <div className={"sb-validation-status " + (issues.length ? "bad" : "good")}>
          {issues.length ? <AlertCircle size={23}/> : <CheckCircle2 size={23}/>}
          <div><strong>{issues.length ? issues.length + " issue(s) to resolve" : "Draft passes local validation"}</strong>
            <p>{issues.length ? "Correct these before using the draft as a specification." : "Syntax and basic semantic checks passed. Live readiness and exchange compatibility are NOT verified."}</p></div>
        </div>
        {issues.length > 0 && <div className="sb-issues">{issues.map((issue,index) =>
          <div key={index}><code>{issue.path}</code><span>{issue.message}</span></div>)}</div>}
        <div className="sb-subpanel sb-review-detail"><h3>Execution summary</h3>
          <dl><dt>Universe</dt><dd>{draft.universe.product_types.join(", ") || "none"}</dd>
            <dt>Signal / protective</dt><dd>{draft.clock.signal_timeframe} / {draft.clock.execution_timeframe}</dd>
            <dt>Order sizing</dt><dd>{draft.execution.sizing.mode.replaceAll("_"," ")} · {draft.execution.sizing.value}</dd>
            <dt>Stop-loss</dt><dd>{draft.risk.stop_loss.mode} · {draft.risk.stop_loss.value}</dd>
            <dt>Live trading</dt><dd>Disabled / unavailable</dd></dl>
        </div>
        <div className="sb-review-actions"><button className="sb-btn-primary" onClick={saveJSON}><Download size={16}/> Download draft JSON</button>
          <button className="sb-btn-muted" onClick={copyJSON}><Braces size={16}/> Copy JSON</button></div>
      </>;
    }
  })();
  return <div className="sb-root">
    <div className="sb-topline">
      <div><div className="sb-eyebrow">STRATEGIES / NO-CODE EDITOR</div><h2>Strategy Builder</h2><p>Build multi-asset, multi-timeframe rules visually. Export for review; no live execution.</p></div>
      <div className="sb-top-buttons"><input type="file" ref={fileInput} accept=".json,application/json" onChange={importJSON} hidden/>
        <button className="sb-btn-muted" onClick={() => fileInput.current?.click()}><Upload size={16}/> Import JSON</button>
        <button className="sb-btn-primary" onClick={saveJSON}><Download size={16}/> Export draft</button>
      </div>
    </div>
    {notice && <div className="sb-notice" role="status"><span>{notice}</span><button aria-label="Dismiss notification" onClick={() => setNotice("")}><X size={15}/></button></div>}
    <div className="sb-layout">
      <nav className="sb-step-nav" aria-label="Strategy configuration">
        <div className="sb-step-nav-title">CONFIGURATION</div>
        {steps.map((item,index) => <button key={item.id} className={"sb-step " + (step === item.id ? "active" : "")}
          onClick={() => setStep(item.id)} aria-current={step === item.id ? "step" : undefined}>
          <span className="sb-step-index">{String(index + 1).padStart(2,"0")}</span>
          <span><strong>{item.label}</strong><small>{item.subtitle}</small></span>
        </button>)}
        <div className="sb-nav-safe"><LockKeyhole size={16}/><span>Offline draft only<br/>No live trading controls</span></div>
      </nav>
      <section className="sb-editor" aria-label="Configure strategy">
        <div className="sb-editor-body">{section}</div>
        <div className="sb-editor-footer"><button className="sb-btn-muted" disabled={stepIndex === 0} onClick={() => setStep(steps[stepIndex-1].id)}>
          <ArrowLeft size={16}/> Previous</button>
          <span>Step {stepIndex+1} of {steps.length}</span>
          <button className="sb-btn-primary" disabled={stepIndex === steps.length-1} onClick={() => setStep(steps[stepIndex+1].id)}>
            Next <ArrowRight size={16}/></button></div>
      </section>
      <aside className="sb-review-rail" aria-label="Draft preview">
        <div className="sb-rail-title"><span><FileJson size={17}/> Live draft preview</span><span className="sb-draft-chip">v{draft.schema_version}</span></div>
        <div className="sb-rail-meta">
          <strong>{draft.metadata.name || "Untitled draft"}</strong>
          <span>{draft.universe.product_types.length} contract types · {draft.clock.signal_timeframe} signal</span>
          <div className={"sb-rail-indicator " + (issues.length ? "bad" : "good")}>{issues.length ? <AlertCircle size={14}/> : <CheckCircle2 size={14}/>}
            {issues.length ? issues.length + " validation issues" : "Local checks passed"}</div>
        </div>
        {issues.length > 0 && <div className="sb-rail-error">{issues[0].path}: {issues[0].message}</div>}
        <div className="sb-rail-code"><pre>{serialized.length > 8500 ? serialized.slice(0,8500) + "\n... preview truncated — export for full JSON" : serialized}</pre></div>
        <div className="sb-rail-actions"><button onClick={copyJSON}><Braces size={15}/> Copy full JSON</button>
          <button onClick={() => setStep("review")}><SlidersHorizontal size={15}/> Review</button></div>
        <div className="sb-rail-safety"><ShieldCheck size={15}/> Draft validity ≠ verified trading safety.</div>
      </aside>
    </div>
  </div>;
}
