# Budibase

Emulates the Budibase API (v1), for local development and tests.

**7 conformance cases, 4 checked against the live API on 2026-09-01.**

## What this Recipe found

Its **documented API host is not an API host** --
`api.budibase.app` is read as a tenant name and answers "Tenant not found".

## Sources

- Documentation: https://docs.budibase.com/docs/public-api
- Machine-readable description: https://raw.githubusercontent.com/Budibase/budibase/master/packages/server/specs/openapi.yaml, last checked 2026-09-01
  `cauldron drift budibase` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve budibase     # run it
cauldron verify budibase -v # check every claim
```
