export const TIMEFRAMES = ["1m", "3m", "5m", "10m", "15m", "30m", "1h", "2h", "4h", "6h", "12h", "1d", "1w"] as const;
export type Timeframe = (typeof TIMEFRAMES)[number];
export const PRODUCT_TYPES = ["perpetual_future", "dated_future", "call_option", "put_option"] as const;
export type ProductType = (typeof PRODUCT_TYPES)[number];
export const WEEKDAYS = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"] as const;
export const PRICE_FIELDS = ["open", "high", "low", "close", "volume", "bid", "ask", "mark", "last", "open_interest", "funding_rate"] as const;
export const OPERATORS = ["gt", "gte", "lt", "lte", "eq", "neq", "crosses_above", "crosses_below"] as const;
export type Operator = (typeof OPERATORS)[number];

export const INDICATORS = [
  { id: "ema", label: "EMA", output: "value", period: 20 },
  { id: "sma", label: "SMA", output: "value", period: 50 },
  { id: "rsi", label: "RSI", output: "value", period: 14 },
  { id: "atr", label: "ATR", output: "value", period: 14 },
  { id: "vwap", label: "VWAP", output: "value", period: 0 },
  { id: "supertrend", label: "Supertrend", output: "trend", period: 10 },
  { id: "macd", label: "MACD", output: "line", period: 12 },
  { id: "stoch_rsi", label: "Stoch RSI", output: "k", period: 14 },
  { id: "adx", label: "ADX", output: "adx", period: 14 },
  { id: "volume_ratio", label: "Volume ratio", output: "value", period: 20 },
  { id: "bollinger", label: "Bollinger Bands", output: "middle", period: 20 },
  { id: "donchian", label: "Donchian", output: "upper", period: 20 }
] as const;
export const INDICATOR_OUTPUTS: Record<string, string[]> = {
  ema: ["value"], sma: ["value"], rsi: ["value"], atr: ["value"], vwap: ["value"],
  supertrend: ["trend", "band"], macd: ["line", "signal", "histogram"],
  stoch_rsi: ["k", "d"], adx: ["adx", "plus_di", "minus_di"],
  volume_ratio: ["value"], bollinger: ["upper", "middle", "lower"],
  donchian: ["upper", "lower", "middle"]
};

export type Operand =
  | { kind: "constant"; value: string }
  | { kind: "price"; field: string; timeframe: Timeframe; offset: number }
  | { kind: "indicator"; indicator: string; parameters: Record<string, string | number | boolean>; output: string; timeframe: Timeframe; offset: number };
export type Comparison = { type: "compare"; id: string; left: Operand; operator: Operator; right: Operand };
export type RuleGroup = { type: "group"; id: string; join: "all" | "any"; children: RuleNode[] };
export type RuleNode = Comparison | RuleGroup;
export type SignalSide = "long_entry" | "short_entry" | "long_exit" | "short_exit";

export type StrategyDefinition = {
  schema_version: "1.0.0";
  id: string;
  revision: number;
  metadata: { name: string; description: string; tags: string[] };
  universe: {
    selection: "dynamic_filter" | "product_ids"; product_types: ProductType[];
    underlying_symbols: string[]; product_ids: number[]; max_products: number;
  };
  clock: {
    signal_timeframe: Timeframe; execution_timeframe: Timeframe; confirm_on_close: true;
    intrabar_protective_exits: boolean; max_price_age_ms: number;
  };
  session: { timezone: string; weekdays: number[]; windows: { start: string; end: string }[] };
  signals: {
    long_entry: RuleGroup; short_entry: RuleGroup; long_exit: RuleGroup; short_exit: RuleGroup;
    conflict_policy: "skip_both" | "prefer_long" | "prefer_short";
  };
  execution: {
    order_type: "limit" | "market"; time_in_force: "gtc" | "ioc" | "fok";
    post_only: boolean; reduce_only_exits: true; max_slippage_bps: number;
    sizing: { mode: "fixed_lots" | "fixed_quote_notional" | "risk_pct_equity"; value: string };
    options: {
      expiry_selection: "nearest_eligible" | "fixed_expiry" | "furthest_eligible";
      expiry_utc?: string; min_dte_hours: number; max_dte_hours: number;
      strike_selection: "atm" | "offset_strikes" | "target_delta";
      strike_offset: number; target_delta?: string;
      direction_map: { long: OptionOrderMapping; short: OptionOrderMapping };
    };
  };
  risk: {
    max_open_positions: number; max_orders_per_minute: number; max_daily_loss_pct: number;
    max_drawdown_pct: number; max_consecutive_losses: number; max_notional_quote: string;
    min_margin_buffer_pct: number; stop_loss: Protection; take_profit: Protection;
    trailing_stop: { enabled: boolean; activation: string; distance: string; unit: "price_ticks" | "percent" | "atr_multiple" };
    halt_on_stale_data: true; halt_on_reconciliation_gap: true;
  };
};
export type Protection = { mode: "off" | "percent" | "price_ticks" | "atr_multiple" | "signal_candle"; value: string; basis?: "contract_entry_price" | "underlying_entry_price" | "option_entry_premium" | "account_equity" };
export type OptionOrderMapping = { option_type: "call" | "put"; action: "buy_to_open" | "sell_to_open" };

export const ruleId = (): string => "rule_" + crypto.randomUUID().replaceAll("-", "").slice(0, 16);
export function makeOperand(kind: Operand["kind"], timeframe: Timeframe): Operand {
  if (kind === "constant") return { kind: "constant", value: "50" };
  if (kind === "price") return { kind: "price", field: "close", timeframe, offset: 0 };
  return { kind: "indicator", indicator: "ema", parameters: { period: 20 }, output: "value", timeframe, offset: 0 };
}
export function makeComparison(timeframe: Timeframe): Comparison {
  return { type: "compare", id: ruleId(), left: makeOperand("price", timeframe), operator: "gt", right: makeOperand("indicator", timeframe) };
}
export function makeGroup(join: "all" | "any" = "all"): RuleGroup {
  return { type: "group", id: ruleId(), join, children: [] };
}
export function createDraft(): StrategyDefinition {
  return {
    schema_version: "1.0.0",
    id: "new_strategy_001", revision: 1,
    metadata: { name: "New strategy draft", description: "", tags: ["draft"] },
    universe: { selection: "dynamic_filter", product_types: ["perpetual_future", "dated_future"], underlying_symbols: [], product_ids: [], max_products: 10 },
    clock: { signal_timeframe: "15m", execution_timeframe: "1m", confirm_on_close: true, intrabar_protective_exits: true, max_price_age_ms: 2000 },
    session: { timezone: "Asia/Kolkata", weekdays: [1, 2, 3, 4, 5], windows: [{ start: "00:00", end: "23:59" }] },
    signals: {
      long_entry: { ...makeGroup(), children: [makeComparison("15m")] },
      short_entry: makeGroup(), long_exit: makeGroup("any"), short_exit: makeGroup("any"),
      conflict_policy: "skip_both"
    },
    execution: {
      order_type: "limit", time_in_force: "gtc", post_only: false, reduce_only_exits: true, max_slippage_bps: 20,
      sizing: { mode: "fixed_lots", value: "1" },
      options: { expiry_selection: "nearest_eligible", min_dte_hours: 24, max_dte_hours: 336, strike_selection: "atm", strike_offset: 0, direction_map: { long: { option_type: "call", action: "buy_to_open" }, short: { option_type: "put", action: "buy_to_open" } } }
    },
    risk: {
      max_open_positions: 1, max_orders_per_minute: 6, max_daily_loss_pct: 2, max_drawdown_pct: 10,
      max_consecutive_losses: 4, max_notional_quote: "2000", min_margin_buffer_pct: 30,
      stop_loss: { mode: "percent", value: "1", basis: "contract_entry_price" }, take_profit: { mode: "off", value: "0" },
      trailing_stop: { enabled: false, activation: "1", distance: "0.5", unit: "percent" },
      halt_on_stale_data: true, halt_on_reconciliation_gap: true
    }
  };
}
export function editNode(group: RuleGroup, id: string, modify: (node: RuleNode) => RuleNode | null): RuleGroup {
  function map(current: RuleGroup): RuleGroup {
    return { ...current, children: current.children.flatMap(node => {
      if (node.id === id) {
        const next = modify(node);
        return next ? [next] : [];
      }
      return [node.type === "group" ? map(node) : node];
    }) };
  }
  return map(group);
}
export function countNodes(node: RuleNode): number {
  return node.type === "compare" ? 1 : node.children.reduce((n, child) => n + countNodes(child), 0);
}
