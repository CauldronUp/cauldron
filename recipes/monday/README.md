# Monday

Emulates the Monday API (2026-01), for local development and tests.

**8 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What writing this Recipe changed

It ships on the mechanism written for Linear: a route that matches on the field
a GraphQL query names, rather than on the path alone. Three Recipes came out of
that one change.

## Sources

- Documentation: https://developer.monday.com/api-reference/docs
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve monday     # run it
cauldron verify monday -v # check every claim
```
