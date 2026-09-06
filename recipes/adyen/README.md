# Adyen

Emulates the Adyen API (v71), for local development and tests.

**11 conformance cases, 2 checked against the live API.**

Most cases here cite documentation rather than an observation. Two were struck live against checkout-test.adyen.com on 2026-09-05, and they corrected one: this file had asserted errorType "authentication_error" for an unauthorised call, which was its own key name leaking through as a default rather than anything Adyen sends. The live host answers errorType "security", and answers it identically whether the X-API-Key header is wrong or missing entirely.

## What this Recipe found

Adyen spells its result the British way: resultCode is "Authorised", not "authorized", so code ported from a provider using the other spelling compares against a string Adyen never sends and every payment reads as failed. The field has twelve values and only one means success -- "Received" and "Pending" mean the outcome arrives later by webhook, and "RedirectShopper" means the shopper hasn't finished -- so branching on != "Authorised" declines payments that are still going to settle.

A refusal carries prose and a code, and only the code is stable: refusalReason is written for a human and its wording changes, while refusalReasonCode -- a string of digits, not a number -- is what should be branched on. Amounts are minor units (1000 EUR is ten euros), except for zero-decimal currencies, which aren't scaled at all.

The webhook notification's own success field is the string "true" rather than a boolean, so a plain truth test on it is also true for "false" -- every failed payment reads as successful to code that writes the obvious thing.

## Sources

- Documentation: https://docs.adyen.com/api-explorer/Checkout/71/overview
- Machine-readable description: https://raw.githubusercontent.com/Adyen/adyen-openapi/main/json/CheckoutService-v71.json, last checked 2026-08-30
  `cauldron drift adyen` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve adyen     # run it
cauldron verify adyen -v # check every claim
```
