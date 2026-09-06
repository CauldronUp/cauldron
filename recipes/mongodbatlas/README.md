# MongoDB Atlas

Emulates the MongoDB Atlas API (v2), for local development and tests.

**9 conformance cases, 2 checked against the live API.**

Struck live 2026-09-05 against cloud.mongodb.com, no account and no key -- and found two different failures where this file's own case expected one shape for both. No Authorization header at all gets a real Digest challenge with a JSON body whose `detail` matched this file exactly (the `errorCode` did not: the real one is `UNEXPECTED_ERROR`, not `NOT_ORG_GROUP_CREATOR`); a Bearer-shaped credential nobody issued gets no body at all, empty. This file's existing case asserted the JSON shape against the second scenario and had never been checked.

## What this Recipe found

Which version of this API you are talking to is a date hidden inside a content type, not a path segment or a header of its own. The same URL, `GET /clusters`, documents three different response schemas under three different `Accept` values (`application/vnd.atlas.2023-01-01+json`, `2023-02-01+json`, `2024-08-05+json`), and the difference is not cosmetic: the legacy view carries `mongoURI`, the actual connection string, and the newest view carries none of it. A client that upgrades its Accept header loses the one field it needed to connect to the database, on the same URL, with the same credentials.

Failure is also documented far more thoroughly than success across this spec: 401 and 403 responses are described on 540 of 333 operations each, more than the 420 operations that describe a 200. The credential scheme is HTTP Digest, a two-round-trip challenge-response mechanism, sitting beside a modern OAuth flow in a cloud API written in the 2020s -- every other provider in this collection authenticates with a single header.

Digest is not served here since Cauldron's auth model expects a header, query parameter, basic or bearer credential in one request rather than a challenge-response handshake; the Recipe uses the bearer-token scheme the same spec also declares. Only two of the three dated cluster views are served, the oldest and newest, since the middle one would just repeat the same claim, and no create route exists, so nothing here actually mints a cluster.

## Sources

- Documentation: https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve mongodbatlas     # run it
cauldron verify mongodbatlas -v # check every claim
```
