# NocoDB

Emulates the NocoDB API (v2), for local development and tests.

**15 conformance cases, 9 checked against the live API on 2026-09-01.**

## What this Recipe found

Its **auth order depends on whether the route names an
id** -- a route with one resolves it first and 404s with the id echoed back,
credential or not.

## Sources

- Documentation: https://docs.nocodb.com/developer-resources/rest-apis
- Machine-readable description: https://nocodb.com/apis/v2/swagger-v2.json, last checked 2026-09-01
  `cauldron drift nocodb` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve nocodb     # run it
cauldron verify nocodb -v # check every claim
```
