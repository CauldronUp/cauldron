# Open-Meteo

Emulates the Open-Meteo API (v1), for local development and tests.

**11 conformance cases, all of them checked against the live API on 2026-08-28.**

## What this Recipe found

The timestamps carry no offset and the same
string means two different hours. Ask for Toronto with a timezone and the hourly
times read `2026-08-28T00:00`; ask for the same coordinates without one and they
read `2026-08-28T00:00` again, four hours earlier, because the default timezone
is GMT. The first is 17.5 °C and the second is 21.7 °C. There is no `Z`, no
`+00:00` and no seconds, so `new Date()` on either parses as local time in
whatever zone the reader is in -- a third answer -- and the only way to be right
is to read `utc_offset_seconds` from the other side of the document and apply it
by hand. **And the coordinates come back moved**: a request for 43.65, -79.38
answers with 43.646603, -79.38272, the nearest cell of the weather model,
silently substituted.

The rest is the same quiet. The hourly data is **parallel arrays** -- `time`,
`temperature_2m` and `precipitation` as three lists whose only relationship is
the index -- with the units in a different object, where `hourly_units.time` is
`"iso8601"`, a serialisation format described as a unit for a field that has
none. Adding a second coordinate changes the top-level type from an object to an
array, and `location_id` is then absent on the first element rather than 0.
`timezone_abbreviation` is `"GMT-4"` rather than `EDT`. And a failure names a
field `error` that is only ever `true`, beside a `reason` that is a Swift
decoder's stringified error, complete with the service's internal generic type
signature -- while a forgotten `latitude` is reported as `Parameter 'latitude'
and 'longitude' must have the same number of elements`.

## Sources

- Documentation: https://open-meteo.com/en/docs
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve openmeteo     # run it
cauldron verify openmeteo -v # check every claim
```
