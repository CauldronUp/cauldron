# Northflank

Emulates the Northflank API (v1), for local development and tests.

**15 conformance cases, 8 checked against the live API on 2026-09-01.**

## What this Recipe found

It **validates the body before the credential** on its
nested build route and after it on its shallow ones -- so whether you are told
your credential is missing depends on whether your JSON parsed.

## Sources

- Documentation: https://northflank.com/docs/v1/api/introduction
- Machine-readable description: https://api.northflank.com/v1/swagger-json, last checked 2026-09-01
  `cauldron drift northflank` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve northflank     # run it
cauldron verify northflank -v # check every claim
```
