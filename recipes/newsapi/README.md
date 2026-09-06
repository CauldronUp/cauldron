# NewsAPI

Emulates the NewsAPI API (v2), for local development and tests.

**8 conformance cases, 3 checked against the live API on 2026-09-01.**

## What this Recipe found

A **fifth spelling of the status inside the body** --
lowercase `"ok"` and `"error"`, matching none of EPSS's restated integer,
Paystack's boolean, Flutterwave's string or Wise's quoted numeral.

## Sources

- Documentation: https://newsapi.org/docs
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve newsapi     # run it
cauldron verify newsapi -v # check every claim
```
