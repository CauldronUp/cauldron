# Deezer

Emulates the Deezer API (v1), for local development and tests.

**12 conformance cases, 11 checked against the live API on 2026-08-30.** The unchecked one is the paging case: it sends the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

**Every failure is HTTP 200.** A record that does not
exist, a path that does not exist and an empty parameter all answer 200 with an
`error` object, so a client checking `response.ok` -- or `status < 400`, or
axios's default -- sees success on every failure this API has. **And one carries
a 500 inside the 200**: the code for an empty parameter is `500`, three digits
that look exactly like a status, beside neighbours of `600` and `800` that are
not HTTP statuses at all.

**And the number of matching records depends on how many you ask for.** The same
query answers `total: 179` at `limit=1` and `limit=5`, `177` at `limit=25`, and
`172` at `limit=100` -- asking for a bigger page makes the collection smaller,
which is backwards from every guess anyone would make. A `limit=0` is ignored
and returns the default twenty-five. The share link comes with Deezer's own
`utm_source`, `utm_medium`, `utm_content` and a `utm_term`, so an application
rendering a share button propagates their attribution without deciding to. Of
five picture fields, the first is an API path where its four siblings are CDN
JPEGs. And explicitness is sent three times in three encodings in one record:
`explicit_lyrics: true`, `explicit_content_lyrics: 1`, `explicit_content_cover:
0`.

## Sources

- Documentation: https://developers.deezer.com/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve deezer     # run it
cauldron verify deezer -v # check every claim
```
