# Magento

Emulates the Magento API (V1), for local development and tests.

**14 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A product is addressed by SKU, not by the id in its own body: `GET /V1/products/ABC-1` works and `GET /V1/products/42` does not, even for a product whose record says `"id": 42`. Code that lists products and then fetches one by the id it just read gets 404 on every one of them. Orders have the opposite problem in reverse -- they are addressed by an internal `entity_id`, while the number printed on the invoice and read down the phone is a separate field, `increment_id`, a zero-padded string like `"000000123"`.

Search is nested brackets in the query string (`searchCriteria[filterGroups][0][filters][0][field]`) rather than anything resembling a normal query parameter, and a client that sends ordinary `limit`/`page` params is not refused -- it is silently answered with page one of everything, which looks like working code until the catalogue grows. A listing wraps its records in `{items, search_criteria, total_count}`; a single record is bare, so `response.items` exists on the list and not on the thing it contained.

Failure messages are unresolved templates -- `"The %fieldName value is required."` with `{"fieldName": "sku"}` sitting beside it in `parameters` -- so showing `message` straight to a user shows them a percent sign and a variable name instead of a sentence.

## Sources

- Documentation: https://developer.adobe.com/commerce/webapi/rest/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve magento     # run it
cauldron verify magento -v # check every claim
```
