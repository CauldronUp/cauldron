# Dynatrace

Emulates the Dynatrace API (v2), for local development and tests.

**6 conformance cases, 2 checked against the live API on 2026-09-01.**

## What this Recipe found

Its **per-tenant hostname never reaches Dynatrace**: a
fabricated environment id answers 23 bytes of AWS gateway default across six
probes, so the ordinary questions are unanswerable without an account.

## Sources

- Documentation: https://docs.dynatrace.com/docs/dynatrace-api/environment-api/entity-v2
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve dynatrace     # run it
cauldron verify dynatrace -v # check every claim
```
