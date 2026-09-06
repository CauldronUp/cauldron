# Aha

Emulates the Aha API (v1), for local development and tests.

**17 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

Its **custom-field schema disagrees with its own values** about
what type a field is, and mints two ids for one field name on two record types.

## Sources

- Documentation: https://www.aha.io/api
- Machine-readable description: https://www.aha.io/openapi.json, last checked 2026-09-01
  `cauldron drift aha` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve aha     # run it
cauldron verify aha -v # check every claim
```
