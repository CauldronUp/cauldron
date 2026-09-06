# MX

Emulates the MX API (v1), for local development and tests.

**16 conformance cases, 5 checked against the live API on 2026-09-02.**

## What this Recipe found

It **has three fields for one question and no answer**: the one
that would say a connection needs re-authentication is never true in any
published example.

## Sources

- Documentation: https://docs.mx.com/api
- Machine-readable description: https://raw.githubusercontent.com/mxenabled/openapi/main/openapi/mx_platform_api.yml, last checked 2026-09-02
  `cauldron drift mx` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve mx     # run it
cauldron verify mx -v # check every claim
```
