# DigitalOcean

Emulates the DigitalOcean API (v2), for local development and tests.

**9 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## Sources

- Documentation: https://docs.digitalocean.com/reference/api/api-reference/
- Machine-readable description: https://raw.githubusercontent.com/digitalocean/openapi/main/specification/DigitalOcean-public.v2.yaml, last checked 2026-08-30
  `cauldron drift digitalocean` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve digitalocean     # run it
cauldron verify digitalocean -v # check every claim
```
