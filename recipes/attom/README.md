# ATTOM

Emulates the ATTOM API (v1.0.0), for local development and tests.

**10 conformance cases, 7 checked against the live API on 2026-09-02.**

## What this Recipe found

Its **own table pairs a failure status with the word
Success** -- a failure-shaped success, where this collection usually finds the
inverse.

## Sources

- Documentation: https://api.developer.attomdata.com/docs
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve attom     # run it
cauldron verify attom -v # check every claim
```
