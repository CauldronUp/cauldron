# Squarespace

Emulates the Squarespace API (1.0), for local development and tests.

**12 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

`REFUNDED` doesn't mean an order was reversed -- Squarespace's own documentation says a partial refund of any amount sets an order to `REFUNDED`, the same status a fully cancelled-and-returned order gets. A two-pound goodwill refund on a four-hundred-and-eighty-pound order lands in the identical state as a total loss, and the field every integration switches on gives no hint which one happened; the arithmetic that tells them apart, `refundedTotal` against `grandTotal`, is available but nothing nudges anyone toward doing it.

Paging and filtering are also mutually exclusive: a cursor "cannot be used with other parameters," full stop, so an ordinary pager that keeps its filters and appends a cursor gets a 400 on page two. And the default order listing quietly excludes `PARTIALLY_PAID` -- exactly the state a Payment Plan order sits in from the deposit until the last instalment, for as long as months.

## Sources

- Documentation: https://developers.squarespace.com/commerce-apis/orders
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve squarespace     # run it
cauldron verify squarespace -v # check every claim
```
