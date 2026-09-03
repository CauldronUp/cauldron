# apify

Emulates the apify API (v2), for local development and tests.

**11 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## Sources

- Documentation: https://docs.apify.com/api/v2
- Machine-readable description: https://docs.apify.com/api/openapi.json, last checked 2026-08-31
  `cauldron drift apify` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve apify     # run it
cauldron verify apify -v # check every claim
```
