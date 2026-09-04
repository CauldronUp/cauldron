# Logto

Emulates the Logto API (Cloud), for local development and tests.

**5 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Logto has no reachable public entry point to probe: `logto.app` alone does not resolve, each tenant lives at its own subdomain, and no publicly documented demo tenant exists. Everything here comes from Logto's own OpenAPI description and, where that spec is silent on errors (it documents success shapes for all 240 paths and no failure shape for any of them), from reading Logto's own MIT-licensed source directly.

The finding worth keeping is that "no credential" is three branches, not two: a missing Authorization header, a header present but not `Bearer`, and a `Bearer` token that fails to verify all produce different error codes. The middle case is the interesting one -- it is the only one of the three whose body carries a `data` field at all (`{supportedTypes: ["Bearer"]}`); the other two never set it, so it is not null on those, it is absent from the JSON entirely, and a client checking for the literal key `data` gets a different answer across three failures that never differ in HTTP status.

No route models a 200 -- inventing an example response for any Management API object would be exactly the fabrication this Recipe avoids -- so every route here is error-only.

## Sources

- Documentation: https://openapi.logto.io/
- Machine-readable description: https://openapi.logto.io/source.json, last checked 2026-09-01
  `cauldron drift logto` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve logto     # run it
cauldron verify logto -v # check every claim
```
