# Printful

Emulates the Printful API (v1), for local development and tests.

**12 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## Sources

- Documentation: https://developers.printful.com/docs/
- Machine-readable description: https://developers.printful.com/docs/openapi.json, last checked 2026-08-31
  `cauldron drift printful` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve printful     # run it
cauldron verify printful -v # check every claim
```
