# OpenRouteService

Emulates the OpenRouteService API (v2), for local development and tests.

**11 conformance cases, 6 checked against the live API on 2026-09-01.**

## What this Recipe found

**Whose own source comment gets that order
backwards**: its tests and its bounding-box serialiser are longitude-first, and
the schema example in its GeoJSON response class shows the same box
latitude-first.

## Sources

- Documentation: https://giscience.github.io/openrouteservice/api-reference/endpoints/directions/requests-and-return-types
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve openrouteservice     # run it
cauldron verify openrouteservice -v # check every claim
```
