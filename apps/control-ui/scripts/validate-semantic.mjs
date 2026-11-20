import { strict as assert } from "node:assert";
import { readFileSync } from "node:fs";
import { mkdtemp, rm } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join, resolve, dirname } from "node:path";
import { pathToFileURL, fileURLToPath } from "node:url";
import { build } from "esbuild";

const root = resolve(dirname(fileURLToPath(import.meta.url)), "../../..");
const temporary = await mkdtemp(join(tmpdir(), "delta-semantic-"));
try {
  const output = join(temporary, "validation.mjs");
  await build({
    entryPoints: [resolve(root, "apps/control-ui/src/strategy/validation.ts")],
    outfile: output,
    bundle: true,
    platform: "node",
    format: "esm",
    target: "node20",
    logLevel: "silent"
  });
  const { validateStrategy } = await import(pathToFileURL(output).href);
  const draft = JSON.parse(readFileSync(resolve(root, "examples/mtf-futures-draft.json"), "utf8"));
  const firstComparison = draft.signals.long_entry.children.find(node => node.type === "compare");
  assert.ok(firstComparison, "Fixture requires a comparison");
  const operandSide = firstComparison.left.kind === "indicator" ? "left" : "right";
  assert.equal(firstComparison[operandSide].kind, "indicator", "Fixture requires an indicator");
  const run = (mutate, expectedPath) => {
    const candidate = structuredClone(draft);
    const comparison = candidate.signals.long_entry.children.find(node => node.type === "compare");
    mutate(comparison[operandSide]);
    const issues = validateStrategy(candidate);
    assert.ok(issues.some(issue => issue.path.endsWith(expectedPath)), JSON.stringify(issues));
  };
  assert.deepEqual(validateStrategy(draft), [], "Valid fixture should pass semantic checks");
  run(operand => { operand.indicator = "not_a_real_indicator"; }, "/indicator");
  run(operand => { operand.indicator = "constructor"; }, "/indicator");
  run(operand => { operand.indicator = "__proto__"; }, "/indicator");
  run(operand => { operand.indicator = "toString"; }, "/indicator");
  run(operand => { operand.output = "not_an_output"; }, "/output");
  run(operand => { operand.parameters.period = 0; }, "/parameters/period");
  run(operand => { operand.parameters.period = 10001; }, "/parameters/period");
  run(operand => { operand.parameters.period = 1.5; }, "/parameters/period");
  run(operand => { operand.parameters.period = -1; }, "/parameters/period");
  console.log("Offline semantic validation: 1 valid fixture + 9 rejection cases passed.");
} finally {
  await rm(temporary, { recursive: true, force: true });
}
