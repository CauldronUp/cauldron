# HERE

Emulates the HERE API (v1), for local development and tests.

**8 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

**The path is checked before the key that is not there**.
An unrouted path 404s and a wrong method 405s with no credential at all, and only
a well-formed request reaches the credential -- where HERE does name which of the
two mistakes you made.

## Sources

- Documentation: https://docs.here.com/geocoding-and-search/reference/get_geocode
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve here     # run it
cauldron verify here -v # check every claim
```
