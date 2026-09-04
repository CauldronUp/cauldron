# Shopware

Emulates the Shopware API (store-api 6.7), for local development and tests.

**16 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Shopware answers products from two endpoints that don't share a paging rule. The generic entity route defaults to a hundred per page and rejects a request for 250; the storefront's category listing defaults to 24, silently truncates a request for 250 to 100, and 404s past the last page instead of returning an empty one. Worse, the entity route's `total` field is, by default, just the size of the current page -- counting rows costs Shopware a second query it doesn't run unless asked -- so a shop with four hundred products can report a ten-record page with `total: 10`, and nothing about the response says which counting mode produced that number.

A wrong access key answers 412 Precondition Failed, not 401 or 403, and no client's retry logic has a branch for that. Money is a float in the shop's own currency, and whether a total is gross or net is decided by a `taxStatus` field sitting beside it that the caller didn't choose.

## Sources

- Documentation: https://shopware.stoplight.io/docs/store-api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve shopware     # run it
cauldron verify shopware -v # check every claim
```
