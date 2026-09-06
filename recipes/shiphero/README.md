# ShipHero

Emulates the ShipHero API (2020-04), for local development and tests.

**13 conformance cases, 2 checked against the live API.**

Everything past the credential check still cites documentation rather than an observation, because reaching it needs a username and a password. The credential check itself was verified directly against public-api.shiphero.com on 2026-09-05.

## What writing this Recipe changed

It ships on the mechanism written for Linear: a route that matches on the field
a GraphQL query names, rather than on the path alone.

## What checking it live found

No Authorization header at all and a fictitious bearer are different sentences -- `{"message":"Authorization header is expected"}` versus `{"message":"Bad token"}` -- and neither had been declared here before this check, so an unauthenticated request served this project's own generic default instead. Both are flat, not GraphQL's errors array: they are checked ahead of the GraphQL engine entirely, by whatever middleware guards it. A path this Recipe does not route at all answers a bare 403 Forbidden HTML page from the load balancer in front of the application, needing no credential and never reaching either shape above -- recorded in the file header rather than encoded, since it belongs to infrastructure this format has no representation for at all.

## Sources

- Documentation: https://developer.shiphero.com/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve shiphero     # run it
cauldron verify shiphero -v # check every claim
```
