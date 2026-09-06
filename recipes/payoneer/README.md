# Payoneer

Emulates the Payoneer API (v4), for local development and tests.

**13 conformance cases, 6 checked against the live API on 2026-09-01.**

## What this Recipe found

**A payment is cancelled by doing nothing**: a debit
expires five minutes after it is created unless committed, and the object's shape
changes across that -- pending carries no status field at all.

## Sources

- Documentation: https://github.com/payoneer/chargeaccount-integration-example
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve payoneer     # run it
cauldron verify payoneer -v # check every claim
```
