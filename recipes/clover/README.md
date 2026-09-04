# Clover

Emulates the Clover API (v3), for local development and tests.

**11 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A Clover order comes back with no items in it -- lineItems is simply absent unless the request explicitly asks with ?expand=lineItems, and a request that forgets isn't refused, it just gets an order that looks empty, which is a perfectly ordinary thing for a real order to be. The practical failure is a sales report that says a busy Saturday sold nothing in particular, with the correct total sitting right there.

Money is integer cents with no currency field on the order at all (it's implicitly the merchant's, and the merchant is a path segment), and dates are Unix milliseconds, so a client that reads a timestamp as seconds lands in 1970. A deleted order is also still there -- Clover soft-deletes, so a total summed without filtering the deleted flag includes voided sales.

Only lineItems and payments are modelled as expandable, and one request reaches one route rather than composing several, so ?expand=lineItems,payments only brings back one of them here, where real Clover sends both -- reproducing the trap of forgetting expand rather than the full expand mechanism.

## Sources

- Documentation: https://docs.clover.com/dev/reference/order
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve clover     # run it
cauldron verify clover -v # check every claim
```
