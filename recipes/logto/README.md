# Logto

Emulates the Logto API (Cloud), for local development and tests.

**5 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## Sources

- Documentation: https://openapi.logto.io/
- Machine-readable description: https://openapi.logto.io/source.json, last checked 2026-09-01
  `cauldron drift logto` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve logto     # run it
cauldron verify logto -v # check every claim
```
