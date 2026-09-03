# Ayrshare

Emulates the Ayrshare API (v1), for local development and tests.

**12 conformance cases, 8 checked against the live API on 2026-09-01.**

## What this Recipe found

Its **partial failure has no word of its own**: a boolean
status, so a mostly-successful batch and a total failure read alike.

## Sources

- Documentation: https://www.ayrshare.com/docs/apis/overview
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve ayrshare     # run it
cauldron verify ayrshare -v # check every claim
```
