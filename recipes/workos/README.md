# WorkOS

Emulates the WorkOS API (v1), for local development and tests.

**9 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

An inactive directory user is still returned by the ordinary listing. Deprovisioning someone in the IdP sets their `state` to `inactive`; it does not remove them, so an integration that syncs everything it's handed re-creates accounts for people who have already left -- the exact failure mode a fixture full of only-active users would never surface. `emails` is also an array with a `primary` flag rather than a string, and the primary one isn't necessarily first, so `emails[0].value` is right most of the time, which is the worst way for a field to be wrong.

## Sources

- Documentation: https://workos.com/docs/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve workos     # run it
cauldron verify workos -v # check every claim
```
