# ShipBob

Emulates the ShipBob API (1.0), for local development and tests.

**6 conformance cases, all of them checked against the live API on 2026-09-01.**

## What this Recipe found

Every one of them is a failure, because **the host has been
switched off**. Its own versioning page said 1.0 would answer 410 Gone from 29
August, and it does, to every credentialed request.

## Sources

- Documentation: https://developer.shipbob.com/v1.0/api/orders/create-order
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve shipbob     # run it
cauldron verify shipbob -v # check every claim
```
