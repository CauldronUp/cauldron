# WorkOS

Emulates the WorkOS API (v1), for local development and tests.

**12 conformance cases, 3 checked against the live API on 2026-09-05.**

Most of this Recipe still cites documentation, since directory sync needs a real IdP connection. The credential and routing shapes needed no account at all, and checking them live found this Recipe's own error shape wrong.

## What this Recipe found

An inactive directory user is still returned by the ordinary listing. Deprovisioning someone in the IdP sets their `state` to `inactive`; it does not remove them, so an integration that syncs everything it's handed re-creates accounts for people who have already left -- the exact failure mode a fixture full of only-active users would never surface. `emails` is also an array with a `primary` flag rather than a string, and the primary one isn't necessarily first, so `emails[0].value` is right most of the time, which is the worst way for a field to be wrong.

## What checking it live found

No credential at all and a present, wrong bearer token both answer a bare `{"message":"Unauthorized"}` with no `code` field at all, where this Recipe had invented one. A path nothing declares answers a third shape again -- `{"message":"The requested resource was not found","error":"Not Found"}`, also with no code -- resolved before any credential is judged.

## Sources

- Documentation: https://workos.com/docs/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve workos     # run it
cauldron verify workos -v # check every claim
```
