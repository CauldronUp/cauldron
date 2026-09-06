# Printify

Emulates the Printify API (v1), for local development and tests.

**15 conformance cases, 3 checked against the live API.**

Everything past the credential and routing checks still cites documentation rather than an observation, because reaching it needs a real shop. Those checks were verified directly against api.printify.com, unauthenticated, on 2026-09-05.

## What this Recipe found

Checked live: no Authorization header at all and a fictitious bearer answer the identical `{"error":"Unauthenticated","request_id":"..."}` -- no code field anywhere in it, which this file had wrong (it declared a numeric 401 code nothing on the wire carries). Routing runs ahead of the credential too: an unrouted path needs nothing sent at all and answers `{"error":"Not found"}`. A wrong method on a real path is its own story, checked but not encoded as a case: the response spells out the allowed methods by name and leaks an internal "public/" route prefix that never appears in the URL a caller sends -- a sentence built per route and per method, too specific to generalise onto a route this Recipe actually models.

Printify splits wholesale and retail differently from Printful's side-by-side fields: the retail total (`total_price`) lives on the order, but what you actually owe lives only on each line item, as `cost` and `shipping_cost` -- there is no `total_cost` anywhere on the order itself. So the margin is not a subtraction between two fields on one record, it is a subtraction between one field and the sum of a nested array, and an order split across two print providers has two costs to add up before the question can even be asked.

Page size is silently capped at ten regardless of what is requested -- asking for a hundred just gets ten back with nothing said about it, so a loop that stops when it receives fewer results than requested stops on the very first page. Timestamps use a space instead of the ISO 8601 `T` (`"2017-04-18 13:24:28+00:00"`), so a strict RFC 3339 parser rejects a string that otherwise looks exactly like a date. Money here is integer cents, while Printful, the other print-on-demand provider in this collection, sends the same kind of quantity as decimal strings, so a shared money helper between the two needs to be told which is which.

## Sources

- Documentation: https://developers.printify.com/
- Machine-readable description: https://developers.printify.com/openapi.json, last checked 2026-08-31
  `cauldron drift printify` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve printify     # run it
cauldron verify printify -v # check every claim
```
