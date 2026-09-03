# VWO

Emulates the VWO API (v2), for local development and tests.

**13 conformance cases, all of them checked against the live API on 2026-09-01.**

## What this Recipe found

Its **routing failure wears a success status**: any
unsupported API version answers HTTP 200 with the error in the body.

## Sources

- Documentation: https://developers.vwo.com/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve vwo     # run it
cauldron verify vwo -v # check every claim
```
