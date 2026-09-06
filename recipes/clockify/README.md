# Clockify

Emulates the Clockify API (v1), for local development and tests.

**9 conformance cases, 4 checked against the live API on 2026-09-03.**

## What this Recipe found

**Every region is the same region** -- four
documented hosts and an invented one all answer identically.

## Sources

- Documentation: https://docs.clockify.me/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve clockify     # run it
cauldron verify clockify -v # check every claim
```
