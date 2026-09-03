# Adyen

Emulates the Adyen API (v71), for local development and tests.

**9 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## Sources

- Documentation: https://docs.adyen.com/api-explorer/Checkout/71/overview
- Machine-readable description: https://raw.githubusercontent.com/Adyen/adyen-openapi/main/json/CheckoutService-v71.json, last checked 2026-08-30
  `cauldron drift adyen` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve adyen     # run it
cauldron verify adyen -v # check every claim
```
