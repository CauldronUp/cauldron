# Clerk

Emulates the Clerk API (2024-10-01), for local development and tests.

**11 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## Sources

- Documentation: https://clerk.com/docs/reference/backend-api
- Machine-readable description: https://clerk.com/openapi.json, last checked 2026-08-31
  `cauldron drift clerk` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve clerk     # run it
cauldron verify clerk -v # check every claim
```
