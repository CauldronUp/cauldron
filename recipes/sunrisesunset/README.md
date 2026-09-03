# Sunrise-Sunset

Emulates the Sunrise-Sunset API (json), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-08-28.**

## What this Recipe found

The sunset is before the sunrise. A request
for Toronto answers `"sunrise": "10:35:37 AM"` and `"sunset": "12:01:46 AM"`:
both UTC, both formatted as a twelve-hour wall clock with no date and no zone.
Toronto's sunset that evening is a minute past midnight UTC -- the next calendar
day -- so parsed against the date that was asked for, the sunset lands sixteen
hours **before** the sunrise, and nothing in the response says the day rolled
over. **And the times are UTC while looking like local time**: `"10:35:37 AM"`
reads as half past ten in the morning, Toronto's sunrise was 06:35, and the only
thing that says otherwise is `tzid: "UTC"`, a sibling of `results` rather than a
field inside it.

A latitude of 999 answers `status: "OK"` with a full set of times -- there is no
validation of the coordinate at all. A failure sets `results` to the **empty
string**, where a success has it as an object, so `typeof results` differs
between them. `day_length` is `"13:26:09"`, a duration in the same
colon-separated shape as the times beside it and the only field in the object
that is not a clock reading. And at 78 degrees north on the solstice every time
is `"12:00:01 AM"` and `day_length` is `"00:00:00"`, with `status: "OK"` -- the
sun does not set there in June, and the API reports a zero-length day.

## Sources

- Documentation: https://sunrise-sunset.org/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve sunrisesunset     # run it
cauldron verify sunrisesunset -v # check every claim
```
