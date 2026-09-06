# FRED

Emulates the FRED API (fred), for local development and tests.

**11 conformance cases, 8 checked against the live API on 2026-09-01.**

## What this Recipe found

It **says three different things about a bad key**: absent,
misshapen, or well-formed but unregistered, all under one status, so only the
prose separates them. Its **XML default is real** -- omit `file_type` and a JSON
API answers `text/xml`, on successes and failures alike.

## Sources

- Documentation: https://fred.stlouisfed.org/docs/api/fred/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve fred     # run it
cauldron verify fred -v # check every claim
```
