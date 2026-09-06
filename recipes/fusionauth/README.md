# FusionAuth

Emulates the FusionAuth API (v1), for local development and tests.

**14 conformance cases, 4 checked against the live API on 2026-09-05.**

The resource cases cite documentation rather than an observation on a real install. FusionAuth is self-hosted, but FusionAuth itself runs a public instance at demo.fusionauth.io -- apparently for exactly this kind of check -- and the refusal cases were struck live, unauthenticated, against it.

## What this Recipe found

A header the schema calls optional becomes mandatory the moment somebody else, in a different system, adds a second tenant. `X-FusionAuth-TenantId` is declared `required: false` on nearly every operation, with the real rule buried in prose: it's only needed once an installation has more than one tenant, or once an API key stops being scoped to just one. A single-tenant sandbox, which is what every developer's own environment is, can never surface this, because that's precisely the state in which the header really is optional. This Recipe seeds two tenants and refuses a request naming neither.

Two fields elsewhere hide a similar gap between what the type says and what the real state needs. A user account can go from usable to blocked without anyone touching it, because `breachedPasswordStatus` reacts to a breach-corpus check that runs on FusionAuth's own schedule, unconnected to any account activity. And whether an account works is split across a boolean (`active`) and a separate expiry date that don't reference each other -- Keycloak's `enabled` has the identical shape of problem one field over, and no provider in this collection solves it in a single field.

The live probe found this file's authentication error wrong: a missing or wrong API key answers zero bytes, not the JSON `generalErrors` shape this file declared. An unrouted path answers 404, also empty, before authentication is ever consulted. It also found a real, documented, credential-free health check (`GET /api/status`) that this Recipe does not model, because this format's required-header check has no way to exempt one path from a requirement every other route needs.

## Sources

- Documentation: https://fusionauth.io/docs/apis/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve fusionauth     # run it
cauldron verify fusionauth -v # check every claim
```
