# razorpay

Emulates the razorpay API (v1), for local development and tests.

**13 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Razorpay catches integrators used to Stripe with one specific behavior: capturing an authorized payment is a separate call, and not making it is itself a decision with consequences. A payment that authorizes sits at `authorized`, and if it is not explicitly captured within the account's configured window, Razorpay automatically refunds it -- the customer sees a charge and then a refund, the order stays unpaid, and nothing in anyone's logs records a failure, because there was not one. Doing nothing was the choice that refunded the money.

Amounts are in paise, a currency's smallest unit that most code has never had to think about -- `50000` is five hundred rupees, and that applies to every amount on every object including fees and tax, so a client that divides by a hundred in one place and not everywhere reports some figure a hundredfold wrong. An order, a payment, and a capture are also three distinct things: one order can have several failed payment attempts against it, and the order's own status only ever says whether it is paid, not which attempt succeeded.

The checkout callback's signature, an HMAC over `order_id` and `payment_id` joined with a pipe, is the entire security boundary for confirming a payment client-side, and it is not computed or verified here on purpose: a fake that accepted any signature would teach exactly the wrong lesson, so this Recipe leaves out the verification endpoint entirely rather than offer one that always passes.

## Sources

- Documentation: https://razorpay.com/docs/api/
- Machine-readable description: https://razorpay.com/openapi.json, last checked 2026-08-31
  `cauldron drift razorpay` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve razorpay     # run it
cauldron verify razorpay -v # check every claim
```
