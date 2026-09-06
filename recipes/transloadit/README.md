# Transloadit

Emulates the Transloadit API (v2), for local development and tests.

**11 conformance cases, 7 checked against the live API on 2026-09-01.**

## What this Recipe found

Its assembly **is finished receiving and not
finished** -- its own two examples show bytes received equal to bytes expected
while the status is still executing.

## Sources

- Documentation: https://transloadit.com/docs/api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve transloadit     # run it
cauldron verify transloadit -v # check every claim
```
