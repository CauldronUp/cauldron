# Fitbit

Emulates the Fitbit API (1), for local development and tests.

**12 conformance cases, 7 checked against the live API on 2026-09-02.**

## What this Recipe found

**A synced empty day and an unsynced day are
identical** -- one sync time exists in the whole API, on the device listing,
attached to no date.

## Sources

- Documentation: https://dev.fitbit.com/build/reference/web-api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve fitbit     # run it
cauldron verify fitbit -v # check every claim
```
