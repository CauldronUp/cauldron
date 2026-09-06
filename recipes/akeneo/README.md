# Akeneo

Emulates the Akeneo API (v1), for local development and tests.

**15 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A product in Akeneo has no name -- it has a list of names. Every attribute's value is an array tagged by locale and channel, so product.name is undefined, product.values.name is an array, and even product.values.name.data is undefined; reading "the name" means finding the one entry matching both a locale and a channel, and there may not be one because attributes are only present in combinations somebody filled in. A client that reads values.name[0].data just gets whichever locale happened to be first.

The envelope is HAL: the array of results sits two levels down at _embedded.items, so response.items and response.data are both undefined. Paging is a link to follow (_links.next), not a page number to construct -- the document says explicitly the page and search_after parameters "should never be set manually". And the count is off by default (with_count=false) with a documented performance warning, the same opt-in-counting shape as Shopware and commercetools.

Prices are a list of decimal-as-string amounts with no fixed scale ("15.5" and "15" appear in the same product), so "15" and "15.00" are the same price and not the same string, and sorting by that field sorts lexically.

## Sources

- Documentation: https://api.akeneo.com/api-reference-index.html
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve akeneo     # run it
cauldron verify akeneo -v # check every claim
```
