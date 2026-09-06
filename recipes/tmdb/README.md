# TMDB

Emulates the TMDB API (v3), for local development and tests.

**11 conformance cases, 3 checked against the live API on 2026-08-31.**

## What this Recipe found

**Six different mistakes are one error.** No key, a fake
key, a fake key on a missing film, a fake bearer token, the wrong verb and a bad
path all answer the same 103 bytes with `status_code: 7`. The key is checked
before anything else, so five of the six are misdiagnosed.

## Sources

- Documentation: https://developer.themoviedb.org/reference/intro/getting-started
- Machine-readable description: https://developer.themoviedb.org/openapi/tmdb-api.json, last checked 2026-08-31
  `cauldron drift tmdb` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve tmdb     # run it
cauldron verify tmdb -v # check every claim
```
