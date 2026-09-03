# Contentful

Emulates the Contentful API (v1), for local development and tests.

**8 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

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
