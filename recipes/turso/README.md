# Turso

Emulates the Turso API (v1), for local development and tests.

**16 conformance cases, 6 checked against the live API on 2026-09-01.**

## What this Recipe found

It **answers a missing header with a JWT parsing error** --
"token contains an invalid number of segments" when no token was sent at all.
Its create returns a truncated projection, three fields where a read gives
eight, omitting the one the request was required to supply.

## Sources

- Documentation: https://docs.turso.tech/api-reference
- Machine-readable description: https://turso.tech/openapi.json, last checked 2026-09-01
  `cauldron drift turso` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve turso     # run it
cauldron verify turso -v # check every claim
```
