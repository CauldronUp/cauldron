# Mezmo

Emulates the Mezmo API (v1), for local development and tests.

**2 conformance cases, all of them checked against the live API on 2026-09-01.**

## What this Recipe found

**Whose two surfaces check things in opposite orders** --
ingest resolves the credential first, management resolves the route first.

## Sources

- Documentation: https://docs.mezmo.com/log-analysis-api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve mezmo     # run it
cauldron verify mezmo -v # check every claim
```
