# Worldpay

Emulates the Worldpay API (v6), for local development and tests.

**9 conformance cases, 7 checked against the live API on 2026-09-01.**

## What this Recipe found

Its **content-type versioning cannot be reached**: the
media type decides which API you get, and Basic auth is resolved first for every
credential state.

## Sources

- Documentation: https://docs.worldpay.com/access/products/card-payments/v6/get-started
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve worldpay     # run it
cauldron verify worldpay -v # check every claim
```
