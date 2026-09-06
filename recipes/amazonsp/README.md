# Amazon SP-API

Emulates the Amazon SP-API API (v0), for local development and tests.

**13 conformance cases, 2 checked against the live API.**

Two were struck live against sellingpartnerapi-na.amazon.com on 2026-09-05: an entirely absent token and a present, invented one both answer 403 code Unauthorized, and this file had given both the same empty "details" string. They are not the same sentence live -- an absent token says "Access token is missing in the request header." and a present, wrong one says "The access token you provided is revoked, malformed or invalid."

## What this Recipe found

The quietest failure in the whole Recipe collection lives here: buyer name, email and shipping address require a second Restricted Data Token, and without it the call does not fail -- it returns 200 with an order that looks complete, personal fields just silently absent. An integration built and tested in a seller's sandbox works perfectly and then can't print a shipping label in production, with nothing in any log explaining why.

Everything is also double-wrapped and double-cased: response.payload.Orders is the array (response.Orders is undefined), and the envelope is lower case while everything inside it is PascalCase (AmazonOrderId, OrderStatus). The rate-limit header is a rate, not a remaining count -- 0.0055 means roughly one request every three minutes, and reading it as "requests left" is wrong by three orders of magnitude in the direction that gets an app throttled. A Pending order has no OrderTotal and no items at all, since Amazon creates the order before payment is confirmed, so summing totals across a day silently produces NaN or under-reports.

Only Orders is modelled; the token exchange, restricted-data-token scoping/expiry, and everything else in the SP-API (Reports, Feeds, FBA, Catalog) are absent.

## Sources

- Documentation: https://developer-docs.amazon.com/sp-api/docs/orders-api-v0-reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve amazonsp     # run it
cauldron verify amazonsp -v # check every claim
```
