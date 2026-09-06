# Drip

Emulates the Drip API (v2), for local development and tests.

**13 conformance cases, 6 checked against the live API on 2026-09-01.**

## What this Recipe found

Its documentation **contradicts itself on the same page** about
whether an inactive person can be modified.

## Sources

- Documentation: https://developer.drip.com/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve drip     # run it
cauldron verify drip -v # check every claim
```
