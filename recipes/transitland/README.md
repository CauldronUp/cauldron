# Transitland

Emulates the Transitland API (v2), for local development and tests.

**5 conformance cases, all of them checked against the live API on 2026-09-03.**

## What this Recipe found

Its **two credential failures differ only by a
header**, and which dates every source feed separately -- the only honest answer
available to an aggregator.

## Sources

- Documentation: https://www.transit.land/documentation/rest-api/
- Machine-readable description: https://transit.land/api/v2/rest/openapi.json, last checked 2026-09-03
  `cauldron drift transitland` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve transitland     # run it
cauldron verify transitland -v # check every claim
```
