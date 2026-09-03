# Diffbot

Emulates the Diffbot API (v3), for local development and tests.

**8 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

Its own client has **a class named for succeeding and
failing at once** -- `ExtractionError`, "The Diffbot API returned a 200 but
reported an extraction failure" -- and whose parser only inspects `errorCode`
after confirming the response was 2xx.

## Sources

- Documentation: https://www.diffbot.com/docs/extract/article
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve diffbot     # run it
cauldron verify diffbot -v # check every claim
```
