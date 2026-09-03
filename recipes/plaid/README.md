# Plaid

Emulates the Plaid API (2020-09-14), for local development and tests.

**7 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## Sources

- Documentation: https://plaid.com/docs/api
- Machine-readable description: https://raw.githubusercontent.com/plaid/plaid-openapi/master/2020-09-14.yml, last checked 2026-08-30
  `cauldron drift plaid` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve plaid     # run it
cauldron verify plaid -v # check every claim
```
