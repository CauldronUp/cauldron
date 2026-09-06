# SurrealDB

Emulates the SurrealDB API (v0), for local development and tests.

**12 conformance cases, 5 checked against the live API on 2026-09-02.**

## What this Recipe found

It **answers the deleting question twice, differently**:
no enum on the instance, a six-value one on the organisation, with the second
value spelled wrong.

## Sources

- Documentation: https://surrealdb.com/docs/cloud
- Machine-readable description: https://api.surrealdb.com/openapi/public.yaml, last checked 2026-09-02
  `cauldron drift surrealdb` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve surrealdb     # run it
cauldron verify surrealdb -v # check every claim
```
