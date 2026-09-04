# Fauna

Emulates the Fauna API (v10), for local development and tests.

**4 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Fauna's own hosted service no longer exists. `db.fauna.com`, `api.fauna.com`, `dashboard.fauna.com` and `status.fauna.com` all resolved to nothing as of a check on 2026-09-01 -- Fauna, Inc. shut the product down in May 2025 and deleted every account. There is no live host to strike, so every case here cites Fauna's own open-source drivers (`fauna-js`, `fauna-go`) instead of an observation, because a driver has to get the wire shape right or it stops parsing its own vendor's responses.

The larger finding is structural: Fauna is a query API, not a resource API. There is exactly one route, `POST /query/1`, and everything -- create, get, list, delete -- happens as a side effect of the FQL text in the body rather than as a verb on a path. A create and a get return the identical `Document` shape (`id`, `coll`, `ts`, plus whatever fields the document holds), because FQL's create expression evaluates to the document it just made the same way any other expression evaluates to its result. This Recipe can't address a document by id the way Cauldron normally expects -- path segment, query param, or named body field -- because the id isn't a discrete value in the request at all, just a substring inside the one opaque field the body has.

## Sources

- Documentation: https://docs.fauna.com/fauna/current/reference/http/reference/query/post
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve fauna     # run it
cauldron verify fauna -v # check every claim
```
