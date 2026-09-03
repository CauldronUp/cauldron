# Duo

Emulates the Duo API (v2), for local development and tests.

**11 conformance cases, 4 checked against the live API on 2026-09-01.**

## What this Recipe found

Its decision **does not exist when the call returns**: an
asynchronous authentication hands back a txid and the verdict appears a poll
later.

## Sources

- Documentation: https://duo.com/docs/authapi
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve duo     # run it
cauldron verify duo -v # check every claim
```
