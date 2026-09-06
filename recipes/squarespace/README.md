# Squarespace

Emulates the Squarespace API (1.0), for local development and tests.

**15 conformance cases, 2 checked against the live API on 2026-09-05.**

Most of this Recipe still cites documentation, since orders and refunds need a real site. The credential shape does not, and checking it live found this Recipe's own error claim wrong.

## What this Recipe found

`REFUNDED` doesn't mean an order was reversed -- Squarespace's own documentation says a partial refund of any amount sets an order to `REFUNDED`, the same status a fully cancelled-and-returned order gets. A two-pound goodwill refund on a four-hundred-and-eighty-pound order lands in the identical state as a total loss, and the field every integration switches on gives no hint which one happened; the arithmetic that tells them apart, `refundedTotal` against `grandTotal`, is available but nothing nudges anyone toward doing it.

Paging and filtering are also mutually exclusive: a cursor "cannot be used with other parameters," full stop, so an ordinary pager that keeps its filters and appends a cursor gets a 400 on page two. And the default order listing quietly excludes `PARTIALLY_PAID` -- exactly the state a Payment Plan order sits in from the deposit until the last instalment, for as long as months.

## What checking it live found

An absent `Authorization` header and a key the site does not know both answer "You are not authorized to do that." with `subtype` present and `null`. This Recipe had read `OAUTH_TOKEN_REQUIRED` off the schema's own enumerated values for exactly this failure and never run the case; live, that value is not sent. The required-`User-Agent` claim could not be checked the same way, and the reason is its own finding: a genuinely blank header gets a 403 HTML page from Squarespace's edge before any application code runs, not the documented 400 -- a bot defence that every ordinary client's default User-Agent already clears, so the two layers never visibly disagree.

## Sources

- Documentation: https://developers.squarespace.com/commerce-apis/orders
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve squarespace     # run it
cauldron verify squarespace -v # check every claim
```
