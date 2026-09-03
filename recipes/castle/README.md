# Castle

Emulates the Castle API (v1), for local development and tests.

**8 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

Its **edge tier crashes on a header with no whitespace** --
500, and an openresty HTML page where every other failure is JSON.

## Sources

- Documentation: https://docs.castle.io/docs/risk-scoring
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve castle     # run it
cauldron verify castle -v # check every claim
```
