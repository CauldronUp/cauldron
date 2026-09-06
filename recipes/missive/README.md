# Missive

Emulates the Missive API (v1), for local development and tests.

**15 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

Its attachment URL is **signed and undated** -- landing
between PDFMonkey's stated hour and Api2Pdf's stated day.

## Sources

- Documentation: https://missiveapp.com/docs/developers/rest-api/endpoints
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve missive     # run it
cauldron verify missive -v # check every claim
```
