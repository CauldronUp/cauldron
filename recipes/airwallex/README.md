# airwallex

Emulates the airwallex API (v1), for local development and tests.

**12 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Airwallex is the payments Recipe where a payment's currency and the currency it settles in can differ -- you charge in one currency, get settled in another, and the difference is a conversion rate fixed at a moment nobody in your code chose. Code that reconciles a settlement against an order total is comparing two numbers that were never meant to match. Balances are also per-currency and separate: having funds in one currency doesn't mean a payout in another currency will succeed, even though the account looks funded.

A payment intent's status runs REQUIRES_PAYMENT_METHOD, REQUIRES_CAPTURE, SUCCEEDED, CANCELLED -- and REQUIRES_CAPTURE is money held but not taken, the same trap Razorpay and Marqeta set. The captured amount and the intent amount are separate fields with no third field recording what a partial capture released, and every amount here is a decimal number rather than a minor-unit integer, the opposite convention from Stripe.

No conversion is actually performed and no rate is real -- the point is only that a settlement carries a different currency and rate at all, which a mock built from one example response never shows.

## Sources

- Documentation: https://www.airwallex.com/docs/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve airwallex     # run it
cauldron verify airwallex -v # check every claim
```
