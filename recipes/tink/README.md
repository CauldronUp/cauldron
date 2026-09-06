# Tink

Emulates the Tink API (v1), for local development and tests.

**20 conformance cases, 9 checked against the live API on 2026-08-31.**

## What this Recipe found

It **answers 401 in five different shapes on one host** --
a nested `error`, RFC 7807 `problem+json`, a grpc-gateway `code` of 16, a flat
`{status, reason}`, and another nested sentence. Three backend families behind
one domain.

## Sources

- Documentation: https://docs.tink.com/api-introduction
- Machine-readable description: https://api.tink.com/swagger.json, last checked 2026-08-31
  `cauldron drift tink` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve tink     # run it
cauldron verify tink -v # check every claim
```
