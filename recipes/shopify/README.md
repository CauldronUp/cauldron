# Shopify

Emulates the Shopify API (2026-01), for local development and tests.

**16 conformance cases, 1 checked against the live API.**

Everything past the credential check still cites documentation rather than an observation, because a real, running store needs an account. One credential shape was checked directly against a shop subdomain Shopify has never issued, on 2026-09-05.

## What writing this Recipe changed

An early round found that routes like `/orders/{id}.json` matched nothing at
all, so every single-object route here answered 404. The suite found it; no
amount of reading the code had.

## What checking it live found

A GET against a nonexistent shop, with no access token or a fictitious one, answers the identical `{"errors":"Not Found"}` either way -- the store's existence is resolved before the credential is. A DELETE on the same path, which this admin resource does not support at all, answers differently depending on whether anything was sent: with nothing at all it answers 401 with the exact sentence this file already had, `[API] Invalid API key or access token (unrecognized login or wrong password)`, now confirmed; with a fictitious token present it falls back to the same 404 the GET requests get. Which check runs first depends on both the method and whether a credential was presented, a three-way split with no single setting in this project's routing mechanism, so most of it is recorded rather than encoded.

## Sources

- Documentation: https://shopify.dev/docs/api/admin-rest
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve shopify     # run it
cauldron verify shopify -v # check every claim
```
