# AviationStack

Emulates the AviationStack API (v1), for local development and tests.

**12 conformance cases, 6 checked against the live API on 2026-09-03.**

## What this Recipe found

Its **own published schema does not match its wire**:
four declared fields, none of which the live error carries.

## Sources

- Documentation: https://docs.apilayer.com/aviationstack/docs/api-documentation
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve aviationstack     # run it
cauldron verify aviationstack -v # check every claim
```
