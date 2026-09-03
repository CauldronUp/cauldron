# Svix

Emulates the Svix API (v1), for local development and tests.

**13 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## Sources

- Documentation: https://api.svix.com/docs
- Machine-readable description: https://api.svix.com/api/v1/openapi.json, last checked 2026-08-31
  `cauldron drift svix` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve svix     # run it
cauldron verify svix -v # check every claim
```
