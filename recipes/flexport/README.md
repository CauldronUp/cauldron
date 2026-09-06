# Flexport

Emulates the Flexport API (2023-07-01), for local development and tests.

**10 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

Its container carries **four date pairs and no status** --
the status enum lives one level up, on the shipment.

## Sources

- Documentation: https://apidocs.flexport.com/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve flexport     # run it
cauldron verify flexport -v # check every claim
```
