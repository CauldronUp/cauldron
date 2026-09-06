# Shopware

Emulates the Shopware API (store-api 6.7), for local development and tests.

**30 conformance cases, 5 checked against the live API on 2026-09-05.**

Most of this Recipe was read out of Shopware's own source rather than a running shop. `demo-frontends.shopware.store` needs no account to check the shape of a refusal, so the credential and routing claims were struck live there.

## What this Recipe found

Shopware answers products from two endpoints that don't share a paging rule. The generic entity route defaults to a hundred per page and rejects a request for 250; the storefront's category listing defaults to 24, silently truncates a request for 250 to 100, and 404s past the last page instead of returning an empty one. Worse, the entity route's `total` field is, by default, just the size of the current page -- counting rows costs Shopware a second query it doesn't run unless asked -- so a shop with four hundred products can report a ten-record page with `total: 10`, and nothing about the response says which counting mode produced that number.

A wrong access key answers 412 Precondition Failed, not 401 or 403, and no client's retry logic has a branch for that. Money is a float in the shop's own currency, and whether a total is gross or net is decided by a `taxStatus` field sitting beside it that the caller didn't choose.

## What checking it live changed

Two guesses from source turned out right and two turned out fixable. A missing `sw-access-key` really is a 401 rendered from Symfony, separate from the 412 a wrong key gets -- this Recipe used to say Cauldron could not tell the two apart because it checked the credential before required headers; `auth.absent_error` now does, and both got their own case. The same probe found routing resolved before the credential on every request, not after: a path nothing declares is a 404 and an unsupported method is a 405, with no access key sent on either -- `after_routing: true` now says so. A fourth answer, a key of the wrong *shape* (not `SWSC`-prefixed) answering its own 403, was first recorded as unmodellable -- `auth.pattern` is the acceptance test itself, so turning it on would have made any `SWSC`-shaped key authenticate rather than being checked against the fixture. A later engine change added `auth.shape`, which sits in front of the key comparison instead of replacing it, and closed the gap: the wrong-shape 403 is now its own verified case too.

## Sources

- Documentation: https://shopware.stoplight.io/docs/store-api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve shopware     # run it
cauldron verify shopware -v # check every claim
```
