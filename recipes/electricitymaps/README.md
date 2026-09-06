# Electricity Maps

Emulates the Electricity Maps API (v3), for local development and tests.

**11 conformance cases, 6 checked against the live API on 2026-09-02.**

## What this Recipe found

Maps', where **forecast-ness lives in the URL**: the field
saying whether a value was estimated is absent from exactly the endpoint that
forecasts.

## Sources

- Documentation: https://app.electricitymaps.com/docs/reference/carbon-intensity/latest
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve electricitymaps     # run it
cauldron verify electricitymaps -v # check every claim
```
