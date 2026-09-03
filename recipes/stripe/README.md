# Stripe

Emulates the Stripe API (2026-06-30), for local development and tests.

**10 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What writing this Recipe changed

Two of the runtime's early bugs were about this provider being treated as the
default. Errors came back in Stripe's nested shape for providers that send a
flat one, and the list envelope carried a `next_cursor` field Stripe does not
send.

That second one is the worst kind of infidelity, because code written against it
works locally and breaks in production.

## Sources

- Documentation: https://docs.stripe.com/api
- Machine-readable description: https://raw.githubusercontent.com/stripe/openapi/master/openapi/spec3.json, last checked 2026-08-30
  `cauldron drift stripe` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve stripe     # run it
cauldron verify stripe -v # check every claim
```
