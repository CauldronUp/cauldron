# Moov

Emulates the Moov API (latest), for local development and tests.

**18 conformance cases, 3 checked against the live API on 2026-09-01.**

## What this Recipe found

A Moov transfer's source and destination can each be a bank account, a card, a wallet, or several other payment-method shapes, and Moov's own schema handles this not with a discriminated union but by declaring every possible shape as an optional sibling field on one object, with a `paymentMethodType` enum saying which one is actually populated. Nothing in the schema enforces that the populated field agrees with the discriminator -- a `paymentMethodType: card-payment` record with a `bankAccount` object populated instead would be accepted by both Moov's own OpenAPI description and by this Recipe, which reproduces the same declare-every-shape approach field by field. Every fixture here keeps the two in agreement by discipline, not by any mechanism that checks.

Checked live, a missing Authorization header and a syntactically fine but fake bearer token come back completely identical: 401, zero response bytes, no Content-Type header at all -- there is no way for a client to tell "you forgot your key" from "your key is wrong." An address Moov does not recognize (wrong method, unregistered path, even a trailing slash on a real collection) answers 403 rather than 404 or 405, confirmed against four different live probes.

No transfer is actually created and no money moves; creation, cancellations, refunds, and reversals are not modelled. Every error observed live carries zero bytes and no Content-Type, which Cauldron's error path cannot fully reproduce since it always answers some JSON body -- this Recipe gets as close as possible with every field suppressed down to `{}`, but that is a real, permanent gap from the true zero-byte response, not a rounding error.

## Sources

- Documentation: https://docs.moov.io/api/authentication/access-tokens/
- Machine-readable description: https://raw.githubusercontent.com/moovfinancial/moov-api-public/main/latest/openapi.yaml, last checked 2026-09-01
  `cauldron drift moov` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve moov     # run it
cauldron verify moov -v # check every claim
```
