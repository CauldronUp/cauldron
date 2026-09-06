# Lever

Emulates the Lever API (v1), for local development and tests.

**18 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

It **answers four different credential complaints**,
including one naming exactly which half of a Basic credential you put the key in.

## Sources

- Documentation: https://hire.lever.co/developer/documentation
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve lever     # run it
cauldron verify lever -v # check every claim
```
