# Cloudflare Stream

Emulates the Cloudflare Stream API (v4), for local development and tests.

**12 conformance cases, 5 checked against the live API on 2026-09-03.**

## What this Recipe found

Stream's, whose **own example contradicts its own field**: a
video in progress carrying an error reason code.

## Sources

- Documentation: https://developers.cloudflare.com/stream/
- Machine-readable description: https://raw.githubusercontent.com/cloudflare/api-schemas/main/openapi.json, last checked 2026-09-03
  `cauldron drift cloudflarestream` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve cloudflarestream     # run it
cauldron verify cloudflarestream -v # check every claim
```
