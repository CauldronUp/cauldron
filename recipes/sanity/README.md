# Sanity

Emulates the Sanity API (v2021-06-07), for local development and tests.

**11 conformance cases, 1 checked against the live API.**

Everything about drafts, queries, and datasets still cites documentation rather than an observation, because reaching it needs a real project subdomain. One credential shape was checked directly against api.sanity.io, the generic host, unauthenticated, on 2026-09-05.

## What this Recipe found

Checked live: a request with no Authorization header at all against the generic host does not reach an auth check at all -- it answers 404, "Use project hostname for data requests", because Sanity's data API is served per-project and the generic host refuses to even attempt the query. A syntactically fine but fictitious bearer changes the code path entirely, answering a flat 401 shape, `{"error":"Unauthorized","errorCode":"SIO-401-ANF","message":"Session not found"}`, which contradicts the nested shape this file had assumed for its default authentication error. Whether a real project-scoped host answers the same way is not settled by this probe.

A draft is the same document with a `drafts.` prefix on its id -- the published version and the draft are two separate documents, and a query that doesn't filter finds both, so a site can render an unpublished draft alongside its live self. System fields are underscore-prefixed (`_id`, `_type`, `_rev`), so `document.id` finds nothing and the field an optimistic write actually needs is `_rev`.

## Sources

- Documentation: https://www.sanity.io/docs/http-api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve sanity     # run it
cauldron verify sanity -v # check every claim
```
