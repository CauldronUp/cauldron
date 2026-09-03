# Cronitor

Emulates the Cronitor API (telemetry+monitors), for local development and tests.

**13 conformance cases, all of them checked against the live API on 2026-09-01.**

## What this Recipe found

**Whose success declares JSON and sends nothing** --
200, `Content-Type: application/json`, zero bytes, on the ping path rather than
a failure. Its management API distinguishes missing from wrong in *four* ways,
because two authentication schemes stack on one route.

## Sources

- Documentation: https://cronitor.io/docs
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve cronitor     # run it
cauldron verify cronitor -v # check every claim
```
