# TomTom

Emulates the TomTom API (2), for local development and tests.

**8 conformance cases, 4 checked against the live API on 2026-09-01.**

## What this Recipe found

Its **Content-Type says text and whose body is JSON**, on
both a genuine RFC 7807 404 and a 401 that is nested under `detailedError` and is
not.

## Sources

- Documentation: https://docs.tomtom.com/geocoding-api/documentation/geocode
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve tomtom     # run it
cauldron verify tomtom -v # check every claim
```
