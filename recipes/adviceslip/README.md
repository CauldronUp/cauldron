# Advice Slip

Emulates the Advice Slip API (adviceslip), for local development and tests.

**6 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

Slip's, where every response carries two `Cache-Control` headers
and they contradict. One says `max-age=3600`; the other says
`max-age=600, private, must-revalidate`. Both are on the wire, on the same 200,
on every request -- and a recipient that follows RFC 9110 joins repeated field
lines with commas, producing one directive list with two `max-age`s in it, which
is not a thing the grammar has an answer for. `fetch`'s `headers.get` returns
exactly that string, so a client reading its own cache policy has to decide
which half to believe.

**And three outcomes carry three different top-level keys.** A slip by id
answers `{"slip": {...}}`, a search answers `{"slips": [...]}`, and a failure
answers `{"message": {...}}` -- singular, plural, and neither, so there is no key
a client can read first to find out what happened. `total_results` is the string
`"1"` beside an `id` that is the number 1. The same slip has two field sets:
`/advice/1` answers `{id, advice}` and the search answers `{id, advice, date}`
for that slip, so a date exists that one route never shows. And a miss and an
empty search are both 200, distinguished only by one word -- `"error"` against
`"notice"` -- in a field called `type` inside an object called `message`.

## Sources

- Documentation: https://api.adviceslip.com/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve adviceslip     # run it
cauldron verify adviceslip -v # check every claim
```
