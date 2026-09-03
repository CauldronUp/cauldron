# SWAPI

Emulates the SWAPI API (swapi), for local development and tests.

**6 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

Every number is a string and one of them has a unit in
it. `rotation_period` is `"23"`, `diameter` is `"10465"`, `population` is
`"200000"` -- and `gravity` is `"1 standard"`. Sorting planets by diameter puts
`"10465"` before `"9000"` because `"1"` comes before `"9"`, and `parseFloat` on
gravity gives 1 while the string it came from means one standard gravity, so the
field that looks most like a number is the one least safely read as one.

**And the absence of data is spelled two ways in one record.** The planet at
`/planets/28` is named `"unknown"`, and its six numeric fields split
three-three: `rotation_period`, `orbital_period` and `diameter` are `"0"`, while
`surface_water`, `population` and `gravity` are `"unknown"`. The same missing
information, two encodings, in the same object -- so a client treating `"0"` as
a real zero charts a planet with no diameter, and one treating `"unknown"` as
the sentinel misses the three not spelled that way. Related records are URLs
rather than identifiers, and a planet that does not exist answers the static
site's own HTML shell with a 200.

## Sources

- Documentation: https://swapi.info/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve swapi     # run it
cauldron verify swapi -v # check every claim
```
