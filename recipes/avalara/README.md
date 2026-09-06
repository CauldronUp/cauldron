# Avalara

Emulates the Avalara API (v2), for local development and tests.

**12 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A tax quote and a tax record are the same call with one word different. POST /transactions/create with type SalesOrder works out the tax and records nothing; the same call with SalesInvoice records a transaction Avalara will file. Both return an identical-looking totalTax and both succeed, so an integration that quotes with SalesOrder and never sends the invoice collects tax from customers all year and files nothing -- discovered by a state government, not by a test. There's a second switch besides the type: a SalesInvoice still isn't filed unless committed, and commit defaults to false, so getting the type right and the commit wrong looks exactly like getting both right.

Tax itself is a breakdown, not a rate -- one address can span a state, county, city and special district each with its own rate, and totalTax is their sum, so reconciling against any single rate never comes out even. The address that comes back is also not the address you sent (Avalara resolves it -- corrected city, ZIP+4), so storing the response as the customer's address silently rewrites what they typed. No tax is actually calculated here; the jurisdiction breakdown is a fixture, and the Recipe models the shape and the two switches, not the arithmetic.

## Sources

- Documentation: https://developer.avalara.com/api-reference/avatax/rest/v2/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve avalara     # run it
cauldron verify avalara -v # check every claim
```
