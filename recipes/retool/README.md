# Retool

Emulates the Retool API (v2), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-09-01.**

## What this Recipe found

It **has no data-table API to model**, and whose
credential check is absolute and three-tiered.

## Sources

- Documentation: https://docs.retool.com/org-users/concepts/retool-api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve retool     # run it
cauldron verify retool -v # check every claim
```
