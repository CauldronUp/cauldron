# NOAA

Emulates the NOAA API (v1), for local development and tests.

**13 conformance cases, 11 checked against the live API on 2026-08-31.** The 2 unchecked ones are the paging cases: they send the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

**An unknown station is a 200 with an empty
array.** Four bytes, no 404, and exactly what a real station with no data in
range also returns -- so nothing tells a typo from a genuine absence, on an
endpoint where every other mistake answers a structured 400 naming the field.
The values are strings all the way down, right-justified and fixed-width
(`"     17.92"`), copied out of the fixed-width files underneath, so sorting by
`TMAX` as text puts 9.5 above 17.92. Its search service disagrees with itself
about counts: `count` said 10000 whatever the limit, `results.length` matched
the limit, and `totalCount` held the real 86045.

## Sources

- Documentation: https://www.ncei.noaa.gov/support/access-data-service-api-user-documentation
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve noaa     # run it
cauldron verify noaa -v # check every claim
```
