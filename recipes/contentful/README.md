# Contentful

Emulates the Contentful API (v1), for local development and tests.

**24 conformance cases, 2 checked against the live API.**

Two were struck live on 2026-09-05 against a real, standing example Contentful publishes in its own docs: space cfexampleapi with the read-only token b4c0n73n7fu1. A request with that token succeeds with real data in the exact shape modelled here; a request with no token, and one with an invented token, are each refused under the same AccessTokenInvalid id with two different sentences -- neither of which this file had modelled for the no-token case before now.

## What writing this Recipe changed

It keeps the identifier at `sys.id` rather than at the top level.

## Sources

- Documentation: https://www.contentful.com/developers/docs/references/content-delivery-api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve contentful     # run it
cauldron verify contentful -v # check every claim
```
