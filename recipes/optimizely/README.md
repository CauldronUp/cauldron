# Optimizely

Emulates the Optimizely API (v2), for local development and tests.

**15 conformance cases, all of them checked against the live API on 2026-09-01.**

## What this Recipe found

Its **assignment is computed on the client** --
MurmurHash3, seed 1, from its own open-source SDK -- so there is no decide
endpoint to ask.

## Sources

- Documentation: https://docs.developers.optimizely.com/feature-experimentation/reference
- Machine-readable description: https://api.optimizely.com/openapi.json, last checked 2026-09-05
  `cauldron drift optimizely` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve optimizely     # run it
cauldron verify optimizely -v # check every claim
```
