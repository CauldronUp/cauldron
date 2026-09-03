# Porter

Emulates the Porter API (undocumented), for local development and tests.

**9 conformance cases, all of them checked against the live API on 2026-09-01.**

## What this Recipe found

It **has no HTTP API to model**. `api.porter.run` is dead,
the documentation publishes no REST reference, and the only live surface is the
dashboard's own backend -- so the Recipe is entirely failures, labelled as such.

## Sources

- Documentation: https://docs.porter.run/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve porter     # run it
cauldron verify porter -v # check every claim
```
