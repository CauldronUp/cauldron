# Estated

Emulates the Estated API (v4), for local development and tests.

**11 conformance cases, 7 checked against the live API on 2026-09-02.**

## What this Recipe found

**The only one in its group that reports its own
confidence**, and whose public gateway root answers today with a copyright line
naming the company that absorbed it.

## Sources

- Documentation: https://estated.com/developers/docs/v4/property/overview
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve estated     # run it
cauldron verify estated -v # check every claim
```
