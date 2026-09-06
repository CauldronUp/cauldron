# Cloudflare

Emulates the Cloudflare API (v4), for local development and tests.

**18 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What writing this Recipe changed

It puts every payload under `result`, success or failure alike, so a client
branching on the HTTP status is checking the wrong thing.

## Sources

- Documentation: https://developers.cloudflare.com/api
- Machine-readable description: https://developers.cloudflare.com/openapi.json, last checked 2026-08-31
  `cauldron drift cloudflare` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve cloudflare     # run it
cauldron verify cloudflare -v # check every claim
```
