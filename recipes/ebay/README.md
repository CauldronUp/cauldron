# eBay

Emulates the eBay API (v1), for local development and tests.

**12 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

eBay is really a dozen APIs under one host built by different teams, and they disagree about basic response shape: Fulfillment's collection key is orders, Browse's is itemSummaries, Inventory's is inventoryItems -- same host, same token, same envelope otherwise -- and the field carrying the total differs too (total on Fulfillment and Browse, size on Inventory), so one paging helper written against one API reads undefined and silently stops after the first page against another. Failures are always an array ({"errors": [...]}) since one request can fail several ways at once, so response.message and response.error are both undefined until a client indexes into the array first.

An inventory item also has no id of its own anywhere in the body -- it's addressed purely by sku in the path -- so there's nothing to key a cache by except the sku you already had. And Content-Language is required on an inventory write and nowhere else, so an integration that only ever reads in its tests meets the requirement for the first time in production, with a failure message about a language rather than about the product being listed.

Application tokens (client-credentials) and user tokens (authorization-code) are scoped differently on the real API -- only a user token reaches Sell endpoints, and sending the wrong kind is one of the most common eBay integration mistakes -- but this format authenticates a request rather than scoping a credential, so one token here reaches everything and that particular mistake can't be reproduced.

## Sources

- Documentation: https://developer.ebay.com/api-docs/static/ebay-rest-landing.html
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve ebay     # run it
cauldron verify ebay -v # check every claim
```
