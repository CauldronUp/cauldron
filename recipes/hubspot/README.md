# HubSpot

Emulates the HubSpot API (v3), for local development and tests.

**23 conformance cases, 4 checked against the live API on 2026-09-05.**

The resource cases cite documentation rather than an observation on a real portal; the refusal cases were struck live, unauthenticated, against api.hubapi.com.

## What writing this Recipe changed

It nests every business attribute under `properties`, so a client reading
`contact.email` finds nothing and has to know to look one level down.

The live probe found this file's error envelope wrong in two ways: the category value lives under a field literally called `category`, not `code`, and every failure carries a constant `"status":"error"` beside it, neither of which this file modelled. The authentication message was also truncated -- the real sentence continues past the first period to explain OAuth 2.0 -- and a missing credential and an invented one get the identical refusal. An unrouted path answers HubSpot's own branded HTML 404 before authentication is ever consulted, and a wrong method on a real path is a genuine empty 405.

## Sources

- Documentation: https://developers.hubspot.com/docs/api/crm/contacts
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve hubspot     # run it
cauldron verify hubspot -v # check every claim
```
