# Apideck

Emulates the Apideck API (unify), for local development and tests.

**15 conformance cases, 2 checked against the live API.**

Two were struck live against unify.apideck.com on 2026-09-05, and corrected a detail this file had guessed at: an absent Authorization header does not get "Verify your Api Key is being set correctly in the authorization header." -- it gets a "detail" naming the missing header directly. Sending any Authorization value at all changes the failure entirely: Apideck refuses for the application id instead, because no app id this project can invent is ever a real, registered one, and that is checked before the token.

## What this Recipe found

A 200 from Apideck's unified ecommerce API can contain somebody else's 429. The envelope has a meta.warnings array that exists only when a downstream connector partially failed -- a shop somewhere got rate-limited, and part of the array is quietly missing -- and every check a client normally makes (status, no error object, data present) passes anyway, because the field that would tell you only appears when something's already wrong.

A valid request can also just be unsupported by whichever connector answers it: the same call with the same parameters works against one shop and is refused for another (UnsupportedFiltersError, PaginationNotSupportedError), and a working integration can start failing the moment a customer connects a second store, because the header selecting which connector to use is optional until there's more than one. And a unified field can be a fact from one provider and an inference from another -- Walmart's payment_status is inferred from order line status because Walmart's own API exposes no such field, while Shopify's is reported directly. Both look identical on the wire.

The credential is three required parts (a bearer token plus two id headers), one of which identifies the end customer rather than the caller -- one more part than VTEX or Voucherify need.

## Sources

- Documentation: https://developers.apideck.com/apis/ecommerce/reference
- Machine-readable description: https://apideck.com/openapi.json, last checked 2026-09-05
  `cauldron drift apideck` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve apideck     # run it
cauldron verify apideck -v # check every claim
```
