# Zenvia

Emulates the Zenvia API (v2), for local development and tests.

**7 conformance cases, 4 checked against the live API on 2026-09-01.**

## What this Recipe found

Its send returns an id **nothing takes back** -- its own
description has no GET for a message or a batch anywhere.

## Sources

- Documentation: https://zenvia.github.io/zenvia-openapi-spec/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve zenvia     # run it
cauldron verify zenvia -v # check every claim
```
