# Papertrail

Emulates the Papertrail API (v1), for local development and tests.

**10 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

Its **405 body is one byte** -- a single space, served
under a Content-Type claiming HTML.

## Sources

- Documentation: https://www.papertrail.com/help/http-api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve papertrail     # run it
cauldron verify papertrail -v # check every claim
```
