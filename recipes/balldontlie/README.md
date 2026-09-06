# balldontlie

Emulates the balldontlie API (v1), for local development and tests.

**10 conformance cases, 4 checked against the live API on 2026-09-02.**

## What this Recipe found

**An API that used to be open and is not**: missing
credential, wrong credential, unrouted path and wrong method all answer the same
401 with the same entity tag.

## Sources

- Documentation: https://docs.balldontlie.io/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve balldontlie     # run it
cauldron verify balldontlie -v # check every claim
```
