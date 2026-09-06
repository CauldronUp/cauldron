# Onfleet

Emulates the Onfleet API (v2), for local development and tests.

**11 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

It has **no 404 reachable anywhere**: every unrouted path
answers 405 with Allow: OPTIONS from a CORS catch-all, the exact mirror of Turso's
no-405-anywhere in this same collection.

## Sources

- Documentation: https://docs.onfleet.com/reference/tasks
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve onfleet     # run it
cauldron verify onfleet -v # check every claim
```
