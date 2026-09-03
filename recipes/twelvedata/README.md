# Twelve Data

Emulates the Twelve Data API (twelvedata), for local development and tests.

**9 conformance cases, 6 checked against the live API on 2026-09-01.**

## What this Recipe found

Data's, written to test whether it repeats Alpha Vantage's
every-failure-is-a-200 and finding that **it does not**. Real statuses
throughout. It does have two 404s sharing nothing but the number: a declared
route sends JSON, an unmatched path sends eighteen bytes of the Go router's
plain text.

## Sources

- Documentation: https://twelvedata.com/docs
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve twelvedata     # run it
cauldron verify twelvedata -v # check every claim
```
