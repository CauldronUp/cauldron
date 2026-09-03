# Rippling

Emulates the Rippling API (v1), for local development and tests.

**6 conformance cases, all of them checked against the live API on 2026-09-01.**

## What this Recipe found

Its **renamed surface did not move**: `/platform/api/workers`
-- the term its own current documentation uses -- 404s, while
`/platform/api/employees` still answers 401.

## Sources

- Documentation: https://developer.rippling.com/documentation/rest-api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve rippling     # run it
cauldron verify rippling -v # check every claim
```
