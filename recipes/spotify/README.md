# Spotify

Emulates the Spotify API (v1), for local development and tests.

**13 conformance cases, 4 checked against the live API on 2026-08-31.**

## What this Recipe found

**The body cannot tell two failures apart and the
header can.** A missing and a bogus token give identical JSON naming three
possible causes, while `WWW-Authenticate` carries `missing_token` against
`invalid_token`. An unknown path answers 410 Gone, a status meaning the resource
was deliberately removed, for one that never existed.

## Sources

- Documentation: https://developer.spotify.com/documentation/web-api
- Machine-readable description: https://developer.spotify.com/reference/web-api/open-api-schema.yaml, last checked 2026-09-01
  `cauldron drift spotify` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve spotify     # run it
cauldron verify spotify -v # check every claim
```
