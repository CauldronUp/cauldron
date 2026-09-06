# SavvyCal

Emulates the SavvyCal API (v1), for local development and tests.

**20 conformance cases, 10 checked against the live API on 2026-09-01.**

## What this Recipe found

It **remembers what its neighbours forget**: eight states
where Calendly has two, and a rescheduled meeting that keeps the time it used to
be at.

## Sources

- Documentation: https://developers.savvycal.com/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve savvycal     # run it
cauldron verify savvycal -v # check every claim
```
