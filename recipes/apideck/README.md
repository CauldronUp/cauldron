# Apideck

Emulates the Apideck API (unify), for local development and tests.

**11 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A 200 from Apideck's unified ecommerce API can contain somebody else's 429. The envelope has a meta.warnings array that exists only when a downstream connector partially failed -- a shop somewhere got rate-limited, and part of the array is quietly missing -- and every check a client normally makes (status, no error object, data present) passes anyway, because the field that would tell you only appears when something's already wrong.

A valid request can also just be unsupported by whichever connector answers it: the same call with the same parameters works against one shop and is refused for another (UnsupportedFiltersError, PaginationNotSupportedError), and a working integration can start failing the moment a customer connects a second store, because the header selecting which connector to use is optional until there's more than one. And a unified field can be a fact from one provider and an inference from another -- Walmart's payment_status is inferred from order line status because Walmart's own API exposes no such field, while Shopify's is reported directly. Both look identical on the wire.

The credential is three required parts (a bearer token plus two id headers), one of which identifies the end customer rather than the caller -- one more part than VTEX or Voucherify need.

## Sources

- Documentation: https://developers.apideck.com/apis/ecommerce/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve apideck     # run it
cauldron verify apideck -v # check every claim
```
