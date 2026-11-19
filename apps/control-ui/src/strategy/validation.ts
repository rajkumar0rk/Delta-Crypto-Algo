import Ajv2020 from "ajv/dist/2020.js";
import addFormats from "ajv-formats";
import schema from "../../../../schemas/strategy.v1.schema.json";
import type { RuleNode, StrategyDefinition } from "./model";

export type ValidationIssue = { path: string; message: string };
const ajv = new Ajv2020({ allErrors: true, strict: false });
addFormats(ajv);
const structural = ajv.compile(schema);

export function validateStrategy(input: unknown): ValidationIssue[] {
  const valid = structural(input);
  if (!valid) {
    return (structural.errors ?? []).map(err => ({
      path: err.instancePath || "/",
      message: err.message ?? "Invalid field"
    })).slice(0, 60);
  }
  const draft = input as StrategyDefinition;
  const issues: ValidationIssue[] = [];
  const add = (path: string, message: string) => issues.push({ path, message });
  if (draft.universe.selection === "product_ids" && draft.universe.product_ids.length === 0) {
    add("/universe/product_ids", "Enter at least one product ID");
  }
  if (draft.clock.execution_timeframe === "1w" && draft.clock.intrabar_protective_exits) {
    add("/clock/execution_timeframe", "Weekly timeframe cannot provide prompt intrabar protection");
  }
  const enabled = draft.signals.long_entry.children.length || draft.signals.short_entry.children.length;
  if (!enabled) add("/signals", "At least one entry side must contain a rule");
  const ids = new Set<string>();
  let total = 0;
  function examine(node: RuleNode, path: string, depth: number) {
    total++;
    if (ids.has(node.id)) add(path + "/id", "Duplicate condition ID");
    ids.add(node.id);
    if (depth > 5) add(path, "Rule nesting may not exceed 5 levels");
    if (node.type === "group") {
      node.children.forEach((child, index) => examine(child, path + "/children/" + index, depth + 1));
    } else {
      if (node.left.kind === "constant" && node.right.kind === "constant") {
        add(path, "At least one side must reference market data");
      }
      for (const [side, operand] of [["left", node.left], ["right", node.right]] as const) {
        if (operand.kind === "indicator") {
          const period = operand.parameters.period;
          if (period !== undefined && (typeof period !== "number" || !Number.isInteger(period) || period < 1 || period > 10000)) {
            add(path + "/" + side + "/parameters/period", "Period must be 1–10,000");
          }
        }
      }
    }
  }
  (["long_entry", "short_entry", "long_exit", "short_exit"] as const).forEach(side => examine(draft.signals[side], "/signals/" + side, 0));
  if (total > 120) add("/signals", "Too many rule nodes; maximum 120");
  const optionsEnabled = draft.universe.product_types.some(kind => kind.endsWith("_option"));
  if (draft.risk.stop_loss.mode === "percent" && draft.risk.stop_loss.basis === "option_entry_premium" && !optionsEnabled) {
    add("/risk/stop_loss/basis", "Premium basis requires an option contract type");
  }
  if (draft.risk.take_profit.mode === "percent" && draft.risk.take_profit.basis === "option_entry_premium" && !optionsEnabled) {
    add("/risk/take_profit/basis", "Premium basis requires an option contract type");
  }
  if (draft.execution.options.min_dte_hours > draft.execution.options.max_dte_hours) {
    add("/execution/options", "Maximum DTE must not be below minimum DTE");
  }
  if (draft.execution.options.strike_selection === "target_delta" &&
      (!draft.execution.options.target_delta || Number(draft.execution.options.target_delta) <= 0 || Number(draft.execution.options.target_delta) > 1)) {
    add("/execution/options/target_delta", "Target delta must be > 0 and <= 1");
  }
  if (draft.execution.post_only && draft.execution.order_type !== "limit") add("/execution/post_only", "Post-only requires a limit order");
  if (draft.execution.sizing.mode === "risk_pct_equity" && draft.risk.stop_loss.mode === "off") {
    add("/risk/stop_loss", "Risk-based sizing requires a protective stop");
  }
  if (Number(draft.execution.sizing.value) <= 0) add("/execution/sizing/value", "Order size must be positive");
  if (Number(draft.risk.max_notional_quote) <= 0) add("/risk/max_notional_quote", "Maximum notional must be positive");
  for (const [name, protection] of [["stop_loss", draft.risk.stop_loss], ["take_profit", draft.risk.take_profit]] as const) {
    if (protection.mode !== "off" && protection.mode !== "signal_candle" && Number(protection.value) <= 0) {
      add("/risk/" + name + "/value", "Enabled protection requires a positive value");
    }
  }
  if (draft.risk.trailing_stop.enabled &&
      (Number(draft.risk.trailing_stop.activation) <= 0 || Number(draft.risk.trailing_stop.distance) <= 0)) {
    add("/risk/trailing_stop", "Trailing activation and distance must be positive");
  }
  for (const [index, window] of draft.session.windows.entries()) {
    if (window.start === window.end) add("/session/windows/" + index, "Start and end cannot be equal");
  }
  return issues.slice(0, 60);
}
