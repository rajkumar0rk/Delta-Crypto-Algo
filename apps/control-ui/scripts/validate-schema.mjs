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
  throw new Error("Sample strategy invalid: " + JSON.stringify(validate.errors));
}
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
console.log("Strategy Schema validation passed: valid sample + 4 rejection cases.");
