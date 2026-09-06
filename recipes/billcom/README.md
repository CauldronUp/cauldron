# Bill.com

Emulates the Bill.com API (v2), for local development and tests.

**16 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Failures in Bill.com arrive with HTTP 200 -- the envelope carries response_status 1 and the actual reason under response_data, so a client that checks the status code treats every failure as a success. isActive is a string, "1" or "2", where 2 means inactive; it's not a boolean and not an enum, and both values are truthy in any language that just checks for truthiness. A paid invoice also isn't a settled one -- paymentStatus moves to paid once the payment is scheduled, and the money doesn't actually leave until the process date days later, during which it can still be stopped.

Every id carries a type prefix that identifies the object kind (a vendor id starts 0056, an invoice 00e, a payment 00o), so passing one id where another belongs fails as an invalid id rather than a wrong-type error. This is one of two Recipes here (with Gusto) where a mistake in a test moves real money rather than burning credits; nothing here actually settles a payment, and whether one has cleared is whatever the fixture says.

## Sources

- Documentation: https://developer.bill.com/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve billcom     # run it
cauldron verify billcom -v # check every claim
```
