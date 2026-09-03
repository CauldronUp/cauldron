# Amplitude

Emulates the Amplitude API (2), for local development and tests.

**7 conformance cases, 6 checked against the live API on 2026-09-03.**

## What this Recipe found

**The only one of three that checks the key on ingest**: an
unrecognised one is refused by name.

## Sources

- Documentation: https://amplitude.com/docs/apis/analytics/http-v2
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve amplitude     # run it
cauldron verify amplitude -v # check every claim
```
