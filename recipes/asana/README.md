# Asana

Emulates the Asana API (1.0), for local development and tests.

**12 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## Sources

- Documentation: https://developers.asana.com/reference/rest-api-reference
- Machine-readable description: https://raw.githubusercontent.com/Asana/openapi/master/defs/asana_oas.yaml, last checked 2026-08-30
  `cauldron drift asana` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve asana     # run it
cauldron verify asana -v # check every claim
```
