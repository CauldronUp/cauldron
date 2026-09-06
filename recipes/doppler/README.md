# Doppler

Emulates the Doppler API (v3), for local development and tests.

**13 conformance cases, 8 checked against the live API on 2026-09-03.**

## What this Recipe found

It **refuses before it routes, without exception**, so
it never tells a caller their path was wrong.

## Sources

- Documentation: https://docs.doppler.com/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve doppler     # run it
cauldron verify doppler -v # check every claim
```
