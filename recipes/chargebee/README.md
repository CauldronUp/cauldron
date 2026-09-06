# Chargebee

Emulates the Chargebee API (v2), for local development and tests.

**16 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Every item in a Chargebee list is double-wrapped under its own resource name -- a subscription list is {"list": [{"subscription": {...}}]} -- so a client has to read list[0].subscription.id, and anyone indexing straight into the item finds nothing. Nothing else in this collection wraps a list item that way. A subscription in "non_renewing" is also still actively serving: it won't renew again, but it hasn't stopped, so code that revokes access on anything other than "active" cuts off a customer who's already paid through the current period.

Chargebee's test sites drift from live the moment plans or coupons are edited, and the states worth testing -- a renewal, a dunning cycle running to its end, a trial converting -- are the ones only time produces, so a test site that can't fast-forward leaves them untested.

## Sources

- Documentation: https://apidocs.chargebee.com/docs/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve chargebee     # run it
cauldron verify chargebee -v # check every claim
```
