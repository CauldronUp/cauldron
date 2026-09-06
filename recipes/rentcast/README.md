# RentCast

Emulates the RentCast API (v1), for local development and tests.

**13 conformance cases, 5 checked against the live API on 2026-09-02.**

## What this Recipe found

Its **two limits fail in two different ways**: the
monthly one bills as overage and never refuses, the per-second one does.

## Sources

- Documentation: https://developers.rentcast.io/reference/introduction
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve rentcast     # run it
cauldron verify rentcast -v # check every claim
```
