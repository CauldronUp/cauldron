# Watchmode

Emulates the Watchmode API (v1), for local development and tests.

**10 conformance cases, 8 checked against the live API on 2026-09-03.**

## What this Recipe found

Its **description guarantees nothing at all** -- an
empty required array, and every foreign identifier nullable.

## Sources

- Documentation: https://api.watchmode.com/docs/
- Machine-readable description: https://api.watchmode.com/openapi.json, last checked 2026-09-03
  `cauldron drift watchmode` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve watchmode     # run it
cauldron verify watchmode -v # check every claim
```
