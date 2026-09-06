# Google Maps Platform

Emulates the Google Maps Platform API (v1), for local development and tests.

**9 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

Maps', whose **status line is honest about routing and silent
about credentials**: a routing miss gets a real 404, and an absent or wrong key
gets 200 with REQUEST_DENIED. The wrong-key sentence ends in a literal trailing
space.

## Sources

- Documentation: https://developers.google.com/maps/documentation/geocoding/overview
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve googlemaps     # run it
cauldron verify googlemaps -v # check every claim
```
