// M1 guard: direct dependencies must be exact; a committed lockfile is still required.
import { readFileSync } from 'node:fs';

const manifest = JSON.parse(readFileSync(new URL('../package.json', import.meta.url), 'utf8'));
const errors = [];
for (const section of ['dependencies', 'devDependencies', 'optionalDependencies']) {
  for (const [name, version] of Object.entries(manifest[section] ?? {})) {
    if (typeof version !== 'string' || !/^\d+\.\d+\.\d+(?:-[0-9A-Za-z.-]+)?(?:\+[0-9A-Za-z.-]+)?$/.test(version)) {
      errors.push(section + '.' + name + ': expected an exact semver version, got ' + JSON.stringify(version));
    }
  }
}
if (errors.length) {
  console.error('Non-reproducible dependency specifications:\n' + errors.join('\n'));
  process.exitCode = 1;
} else {
  console.log('All direct UI dependency versions are pinned exactly (transitive lockfile still required).');
}
