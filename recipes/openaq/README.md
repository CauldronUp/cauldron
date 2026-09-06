# OpenAQ

Emulates the OpenAQ API (v3), for local development and tests.

**13 conformance cases, 8 checked against the live API on 2026-08-31.**

## What this Recipe found

**The retired version is the one that still answers.**
v3 needs a key for everything; v2 needs none and returns 410 saying it is gone.
Its two v3 failures come from different layers and disagree, and the more nearly
correct the credential, the less the error resembles the documented one.

## Sources

- Documentation: https://docs.openaq.org
- Machine-readable description: https://api.openaq.org/openapi.json, last checked 2026-08-31
  `cauldron drift openaq` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve openaq     # run it
cauldron verify openaq -v # check every claim
```
