# Mapbox

Emulates the Mapbox API (v5), for local development and tests.

**10 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

**Longitude comes first and the response says so** --
consistently, in `center`, `geometry.coordinates` and a reverse lookup's path
parameter. A client that transposes them gets the wrong hemisphere and a
perfectly successful response.

## Sources

- Documentation: https://docs.mapbox.com/api/search/geocoding-v5/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve mapbox     # run it
cauldron verify mapbox -v # check every claim
```
