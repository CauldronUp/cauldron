# Stripe

Emulates the Stripe API (2026-06-30), for local development and tests.

**10 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## Sources

- Documentation: https://docs.stripe.com/api
- Machine-readable description: https://raw.githubusercontent.com/stripe/openapi/master/openapi/spec3.json, last checked 2026-08-30
  `cauldron drift stripe` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve stripe     # run it
cauldron verify stripe -v # check every claim
```
