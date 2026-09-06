# Workable

Emulates the Workable API (v3), for local development and tests.

**18 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

**The same person applying twice is two people** --
job-scoped candidate profiles, correlated only by matching the raw email by eye,
with no merge.

## Sources

- Documentation: https://workable.readme.io/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve workable     # run it
cauldron verify workable -v # check every claim
```
