# EIA

Emulates the EIA API (v2), for local development and tests.

**11 conformance cases, 6 checked against the live API on 2026-09-02.**

## What this Recipe found

Its **two layers disagree about what an error is** while using
the same key for it: an object from the gateway, a bare sentence from the
application.

## Sources

- Documentation: https://www.eia.gov/opendata/documentation.php
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve eia     # run it
cauldron verify eia -v # check every claim
```
