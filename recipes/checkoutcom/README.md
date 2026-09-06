# Checkout.com

Emulates the Checkout.com API (v1), for local development and tests.

**17 conformance cases, 5 checked against the live API on 2026-09-02.**

## What this Recipe found

**Whose refusal is zero bytes** -- no body, no
Content-Type, across four credential shapes and five methods.

## Sources

- Documentation: https://api-reference.checkout.com
- Machine-readable description: https://api-reference.checkout.com/v1/swagger.yaml, last checked 2026-09-02
  `cauldron drift checkoutcom` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve checkoutcom     # run it
cauldron verify checkoutcom -v # check every claim
```
