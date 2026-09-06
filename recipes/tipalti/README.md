# Tipalti

Emulates the Tipalti API (v1), for local development and tests.

**14 conformance cases, 4 checked against the live API on 2026-09-01.**

## What this Recipe found

Its payee is **provisional from birth** -- isPayable false
in its very first create response, and read-only on every model that could set it.

## Sources

- Documentation: https://documentation.tipalti.com/docs/getting-started
- Machine-readable description: https://dash.readme.com/api/v1/api-registry/z8pm3amsomq9lm, last checked 2026-09-01
  `cauldron drift tipalti` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve tipalti     # run it
cauldron verify tipalti -v # check every claim
```
