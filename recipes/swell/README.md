# Swell

Emulates the Swell API (backend), for local development and tests.

**14 conformance cases, 6 checked against the live API on 2026-09-01.**

## What this Recipe found

**A cart and an order point at each other** in two
collections, which is the friendliest answer in its group.

## Sources

- Documentation: https://developers.swell.is/backend-api/introduction
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve swell     # run it
cauldron verify swell -v # check every claim
```
