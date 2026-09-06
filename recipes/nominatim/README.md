# Nominatim

Emulates the Nominatim API (nominatim), for local development and tests.

**10 conformance cases, 9 checked against the live API on 2026-08-28.** The unchecked one is the paging case: it sends the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

Nothing found is an empty array on one path and an
error object on the other, and both are 200. `GET /search?q=...` with no matches
answers `[]`; `GET /reverse?lat=0&lon=0` with no match answers
`{"error": "Unable to geocode"}`. Same service, same `format=jsonv2`, same
status, and the reverse one is the expensive shape: an object with an error key
where an object of results should be, so `response.lat` is undefined and
`response.ok` is true. **And search answers with an array while reverse answers
with a bare object** -- not an array of one, an object -- so the two endpoints of
one geocoder cannot share a parser.

The rest is what OpenStreetMap's data model does to a JSON API. `lat` and `lon`
are **strings**, so arithmetic concatenates and comparison sorts lexically.
`boundingbox` is four strings in the order south, north, west, east -- neither
GeoJSON's `[west, south, east, north]` nor a pair of coordinates in sequence.
The keys of `address` depend on the place, and one of them carries an
administrative level inside its own name: `ISO3166-2-lvl4`, where the 4 is the
depth at which that country keeps its subdivisions and is 3 or 6 elsewhere.
`importance` arrives in scientific notation. `place_id` is local to the
installation and changes when the database is rebuilt, so the stable identity is
`osm_type` and `osm_id`, two fields. And a request without a `User-Agent` is
refused in **plain text** -- 403 and one line pointing at the usage policy --
from an endpoint that was asked for JSON.

## Sources

- Documentation: https://nominatim.org/release-docs/latest/api/Overview/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve nominatim     # run it
cauldron verify nominatim -v # check every claim
```
