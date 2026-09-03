# Tigris

Emulates the Tigris API (v1), for local development and tests.

**5 conformance cases, all of them checked against the live API on 2026-09-01.**

## What this Recipe found

It **echoes your bucket name whatever it is**, so the
credential verdict is the only thing that varies.

## Sources

- Documentation: https://www.tigrisdata.com/docs/api/s3/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve tigris     # run it
cauldron verify tigris -v # check every claim
```
