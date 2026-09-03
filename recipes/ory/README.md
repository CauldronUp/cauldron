# Ory

Emulates the Ory API (Cloud), for local development and tests.

**4 conformance cases, 2 checked against the live API on 2026-09-01.**

## What this Recipe found

Its hosted playground **answers one 404 to everything** -- the
same 155 bytes across four distinct paths, with a rule identifier that never
varies. Everything else there is behind a bot challenge, so Ory's own software
is never reached.

## Sources

- Documentation: https://www.ory.com/docs/reference/api
- Machine-readable description: https://raw.githubusercontent.com/ory/docs/master/docs/reference/api.json, last checked 2026-09-01
  `cauldron drift ory` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve ory     # run it
cauldron verify ory -v # check every claim
```
