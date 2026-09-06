# Radar

Emulates the Radar API (v1), for local development and tests.

**10 conformance cases, 8 checked against the live API on 2026-09-01.**

## What this Recipe found

It **tells you to fetch the wrong kind of key** -- asking
for the publishable one on an endpoint documented as needing the secret.

## Sources

- Documentation: https://docs.radar.com/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve radar     # run it
cauldron verify radar -v # check every claim
```
