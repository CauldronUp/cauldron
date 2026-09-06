# Raygun

Emulates the Raygun API (v3), for local development and tests.

**14 conformance cases, 3 checked against the live API on 2026-09-01.**

## What this Recipe found

It **announces JSON and sends fifteen bytes of prose**:
`Invalid API key`, no quotes, no braces, under `Content-Type: application/json`.

## Sources

- Documentation: https://api.raygun.io/v3/swagger/index.html
- Machine-readable description: https://api.raygun.io/v3/raygun-openapi-spec.json, last checked 2026-09-01
  `cauldron drift raygun` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve raygun     # run it
cauldron verify raygun -v # check every claim
```
