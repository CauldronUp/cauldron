# MusicBrainz

Emulates the MusicBrainz API (ws/2), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-08-28.**

## What this Recipe found

It answers XML unless you ask for JSON and explains
that in XML. The default on `/ws/2/artist/{mbid}` is
`<?xml version="1.0" ...>`; `?fmt=json` gets a JSON document; and `?fmt=yaml`
gets **406 in application/xml**, carrying the message that the recognised types
are application/json and application/xml. The one answer that would tell you how
to ask for JSON is the one you cannot parse without already having solved the
problem. **And an identifier that does not exist is a 400 about the identifier's
format**: a perfectly well-formed UUID that nothing in the database uses answers
`{"error": "Invalid mbid."}`, exactly as a string that is not a UUID at all does,
so "no such artist" and "you sent nonsense" are one answer and there is no 404 in
the shape. The two bodies arrive with their keys in opposite orders.

A request without a `User-Agent` is **403 Forbidden with a message about being
throttled** -- the status says you are not allowed, the text says you are going
too fast, and the cause is neither -- and that body is `application/json` with no
charset where every other JSON answer declares one. The field names are
hyphenated (`sort-name`, `type-id`, `begin-area`, `iso-3166-1-codes`), so in
JavaScript every one has to be reached with brackets. `life-span` carries
`begin: "1987"` beside `end: "1994-04-05"` -- two precisions, both strings,
neither saying which. And there are four spellings of "not applicable" in one
document: `gender` is null, `end-area` is null, `ipis` is `[]`, and inside `area`
the `disambiguation` is `""` while `type` is null.

## Sources

- Documentation: https://musicbrainz.org/doc/MusicBrainz_API
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve musicbrainz     # run it
cauldron verify musicbrainz -v # check every claim
```
