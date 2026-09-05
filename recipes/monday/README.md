# Monday

Emulates the Monday API (2026-01), for local development and tests.

**9 conformance cases, 1 checked against the live API.**

Struck live 2026-09-05 against api.monday.com, no account and no key -- and found this file had never declared an auth failure at all, so a real refusal fell back to the runtime's own generic placeholder. Added now: a missing credential and a wrong one both answer `{"errors":[{"message":"Not authenticated","extensions":{"code":"NOT_AUTHENTICATED"}}]}`.

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
