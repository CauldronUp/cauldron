# WattTime

Emulates the WattTime API (v3), for local development and tests.

**5 conformance cases, 4 checked against the live API on 2026-09-02.**

## What this Recipe found

Its **refusals are bare words under a JSON label** -- a
single unquoted word, no braces, announced as JSON.

## Sources

- Documentation: https://docs.watttime.org/
- Machine-readable description: https://docs.watttime.org/openapi.json, last checked 2026-09-02
  `cauldron drift watttime` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve watttime     # run it
cauldron verify watttime -v # check every claim
```
