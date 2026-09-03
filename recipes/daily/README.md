# Daily

Emulates the Daily API (v1), for local development and tests.

**15 conformance cases, 8 checked against the live API on 2026-09-02.**

## What this Recipe found

Its **expired room becomes a zombie** -- its own word, for a
room that is not deleted and does not exist.

## Sources

- Documentation: https://docs.daily.co/docs/rest-api
- Machine-readable description: https://docs.daily.co/openapi.json, last checked 2026-09-02
  `cauldron drift daily` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve daily     # run it
cauldron verify daily -v # check every claim
```
