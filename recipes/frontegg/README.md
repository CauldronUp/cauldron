# Frontegg

Emulates the Frontegg API (1.0), for local development and tests.

**9 conformance cases, 6 checked against the live API on 2026-09-01.**

## What this Recipe found

It **names the host you got wrong**: a fabricated
hostname answers `Failed to find vendor for host:` and quotes it back, where the
Make Recipe records a vendor that cannot manage the same thing.

## Sources

- Documentation: https://developers.frontegg.com/guides/getting-started/home
- Machine-readable description: https://raw.githubusercontent.com/frontegg/openapi-public/master/identity.json, last checked 2026-09-01
  `cauldron drift frontegg` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve frontegg     # run it
cauldron verify frontegg -v # check every claim
```
