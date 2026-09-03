# Linear

Emulates the Linear API (graphql), for local development and tests.

**8 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What writing this Recipe changed

This Recipe was on a list of things the format could not do, for being GraphQL:
a single endpoint where the request body decides the response shape.

It ships because a route learned to match on the field a query names. Monday and
ShipHero shipped on the same mechanism. The entry that excluded it had sat
unchallenged because nothing rereads a list of things that cannot be done, and a
reason that has quietly expired reads exactly like one that still holds.

## Sources

- Documentation: https://linear.app/developers/graphql
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve linear     # run it
cauldron verify linear -v # check every claim
```
