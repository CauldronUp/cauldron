# Anvil

Emulates the Anvil API (graphql), for local development and tests.

**11 conformance cases, 10 checked against the live API on 2026-09-01.**

## What this Recipe found

It **answers an authentication failure four ways** across two
surfaces and two schemes, two of them with a 200.

## Sources

- Documentation: https://www.useanvil.com/docs/api/getting-started/
- Machine-readable description: https://app.useanvil.com/api/v1/openapi.json, last checked 2026-09-05
  `cauldron drift anvil` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve anvil     # run it
cauldron verify anvil -v # check every claim
```
