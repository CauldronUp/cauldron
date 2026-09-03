# Lithic

Emulates the Lithic API (v1), for local development and tests.

**14 conformance cases, 3 checked against the live API on 2026-09-01.**

## Sources

- Documentation: https://docs.lithic.com/reference/authentication
- Machine-readable description: https://raw.githubusercontent.com/lithic-com/lithic-openapi/main/lithic-openapi.yml, last checked 2026-09-01
  `cauldron drift lithic` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve lithic     # run it
cauldron verify lithic -v # check every claim
```
