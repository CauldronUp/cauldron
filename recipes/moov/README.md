# Moov

Emulates the Moov API (latest), for local development and tests.

**10 conformance cases, 3 checked against the live API on 2026-09-01.**

## Sources

- Documentation: https://docs.moov.io/api/authentication/access-tokens/
- Machine-readable description: https://raw.githubusercontent.com/moovfinancial/moov-api-public/main/latest/openapi.yaml, last checked 2026-09-01
  `cauldron drift moov` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve moov     # run it
cauldron verify moov -v # check every claim
```
