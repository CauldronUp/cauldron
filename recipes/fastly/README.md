# Fastly

Emulates the Fastly API (v1), for local development and tests.

**21 conformance cases, 11 checked against the live API on 2026-09-03.**

## What this Recipe found

**Locked and active are independent and only one is
provable**: no response anywhere shows a version active.

## Sources

- Documentation: https://www.fastly.com/documentation/reference/api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve fastly     # run it
cauldron verify fastly -v # check every claim
```
