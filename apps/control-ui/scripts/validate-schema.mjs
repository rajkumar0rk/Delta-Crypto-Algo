import { readFileSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import Ajv2020 from "ajv/dist/2020.js";
import addFormats from "ajv-formats";

const root = resolve(dirname(fileURLToPath(import.meta.url)), "../../..");
const readJSON = path => JSON.parse(readFileSync(resolve(root, path), "utf8"));
const schema = readJSON("schemas/strategy.v1.schema.json");
const ajv = new Ajv2020({ allErrors: true, strict: false });
addFormats(ajv);
const validate = ajv.compile(schema);

const example = readJSON("examples/mtf-futures-draft.json");
if (!validate(example)) {
  throw new Error("Futures draft invalid: " + JSON.stringify(validate.errors));
}
const optionsDraft = readJSON("examples/stoch-rsi-options-draft.json");
if (!validate(optionsDraft)) {
  throw new Error("Options draft invalid: " + JSON.stringify(validate.errors));
}
const missingMapping = structuredClone(optionsDraft);
delete missingMapping.execution.options.direction_map;
if (validate(missingMapping)) throw new Error("Option action mapping was not required");
const missingPremiumBasis = structuredClone(optionsDraft);
delete missingPremiumBasis.risk.stop_loss.basis;
if (validate(missingPremiumBasis)) throw new Error("Percent SL basis was not required");
const unknown = structuredClone(example);
unknown.untrusted_code = "process.exit(0)";
if (validate(unknown)) throw new Error("Unexpected property was accepted");
const missing = structuredClone(example);
delete missing.risk;
if (validate(missing)) throw new Error("Missing risk rules were accepted");
const liveOverride = structuredClone(example);
liveOverride.clock.confirm_on_close = false;
if (validate(liveOverride)) throw new Error("Unconfirmed candle execution was accepted");
const numericPrice = structuredClone(example);
numericPrice.execution.sizing.value = 0.5;
if (validate(numericPrice)) throw new Error("Binary float money sizing was accepted");
console.log("Strategy Schema validation passed: 2 valid samples + 6 rejection cases.");
