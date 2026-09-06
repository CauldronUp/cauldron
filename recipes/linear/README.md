# Linear

Emulates the Linear API (graphql), for local development and tests.

**13 conformance cases, 1 checked against the live API.**

Struck live 2026-09-05 against api.linear.app, no account and no key -- and found a declared error that nothing could ever reach. This file's authentication failure was named "authentication" and declared status 400; the runtime's own credential check falls back to the name "authentication_error" on a real refusal, so every actual unauthenticated request was serving a generic placeholder instead of this entry. Renamed and fixed to the real status, 401 -- the sentence itself was already right.

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
