# Flagsmith

Emulates the Flagsmith API (v1), for local development and tests.

**10 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

It **answers three failures in three formats**: a
printed Python tuple served as JSON, a zero-byte 404, and a bare array holding
one object. No field in common.

## Sources

- Documentation: https://docs.flagsmith.com/clients/rest
- Machine-readable description: https://api.flagsmith.com/api/v1/swagger.json, last checked 2026-09-01
  `cauldron drift flagsmith` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve flagsmith     # run it
cauldron verify flagsmith -v # check every claim
```
