# Airbrake

Emulates the Airbrake API (v4), for local development and tests.

**13 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

Its **own message hedges between two refusals**: one
eighty-one-byte sentence saying not found *or* access denied, for a missing key,
a garbage key, and absurd project ids alike.

## Sources

- Documentation: https://docs.airbrake.io/docs/devops-tools/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve airbrake     # run it
cauldron verify airbrake -v # check every claim
```
