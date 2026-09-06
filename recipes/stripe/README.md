# Stripe

Emulates the Stripe API (2026-06-30), for local development and tests.

**24 conformance cases, 4 checked against the live API on 2026-09-05.**

The customer and payment intent behaviour still cites documentation, since a real charge needs a real account. The credential and routing shapes do not, and checking them live found two sentences this Recipe had merged into one.

## What writing this Recipe changed

Two of the runtime's early bugs were about this provider being treated as the
default. Errors came back in Stripe's nested shape for providers that send a
flat one, and the list envelope carried a `next_cursor` field Stripe does not
send.

That second one is the worst kind of infidelity, because code written against it
works locally and breaks in production.

## What checking it live found

No API key at all and a present-but-wrong one are different sentences, not one: absence gets a long explain-yourself message naming the header and the auth scheme, and a wrong key gets a short one that echoes the key back (masked on the real API; this Recipe echoes it unmasked, since masking isn't a scheme this project implements). `auth.absent_error` now keeps them apart. And neither is checked first on a path or method the API does not have: a path nothing declares and a method a real path does not support both answer `"Unrecognized request URL"` with a 404, no key needed to produce it -- so `after_routing: true` reflects an ordering this Recipe had not stated either way.

## Sources

- Documentation: https://docs.stripe.com/api
- Machine-readable description: https://raw.githubusercontent.com/stripe/openapi/master/openapi/spec3.json, last checked 2026-08-30
  `cauldron drift stripe` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve stripe     # run it
cauldron verify stripe -v # check every claim
```
