# football-data

Emulates the football-data API (v4), for local development and tests.

**10 conformance cases, 7 checked against the live API on 2026-09-02.**

## What this Recipe found

Its **open route breaks when you authenticate** --
the competitions list works with no credential and answers 400 with a wrong
one.

## Sources

- Documentation: https://docs.football-data.org/general/v4/index.html
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve footballdata     # run it
cauldron verify footballdata -v # check every claim
```
