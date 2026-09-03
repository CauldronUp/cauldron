# Nile

Emulates the Nile API (v2), for local development and tests.

**7 conformance cases, all of them checked against the live API on 2026-09-01.**

## What this Recipe found

**Existence is unknowable before you authenticate** -- one
empty 401 whether the database and tenant are real or invented.

## Sources

- Documentation: https://thenile.dev/docs/auth/api-reference/tenants/get-a-tenant
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve nile     # run it
cauldron verify nile -v # check every claim
```
