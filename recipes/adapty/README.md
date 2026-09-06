# Adapty

Emulates the Adapty API (v2), for local development and tests.

**9 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

Its **documented credential scheme is not the live one**:
`Api-Key` cannot be told from sending nothing, and the undocumented Bearer scheme
turns out to be the one its gateway parses.

## Sources

- Documentation: https://adapty.io/docs/api-adapty
- Machine-readable description: https://adapty.io/docs/api-specs/adapty-api.yaml, last checked 2026-09-01
  `cauldron drift adapty` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve adapty     # run it
cauldron verify adapty -v # check every claim
```
