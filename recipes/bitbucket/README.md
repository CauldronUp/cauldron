# Bitbucket

Emulates the Bitbucket API (2.0), for local development and tests.

**9 conformance cases, 4 checked against the live API on 2026-08-22.**

## What this Recipe found

Its public workspaces answer without a token. Its
listing envelope is exactly `values`, `size`, `pagelen`, `page` and `next`,
with no `data` and no `items` -- the two names a client reaches for first,
both finding undefined rather than an error -- and `size` was 407 against a
`pagelen` of 2, so it counts everything rather than what arrived. A repository
carries `is_private` as a boolean and no `visibility` key at all. A missing
repository answers two keys, `type` and `error`, and `error` holds nothing but
`message`.

The fourth is the one worth reading twice. `?pagelen=1` answers one value.
`?limit=1` answers ten, with `pagelen: 10` in the body: the wrong name is not
refused, it is ignored, and a full page comes back looking like a successful
request for one. That is the failure this README keeps describing, watched
happening.

Checking Bitbucket also corrected a message this project had invented. Its 404
had a generic "Resource not found", and Bitbucket's own words are that the
repository may not exist *or* may not be visible to you, and it declines to
say which. Code treating a 404 as a deletion -- dropping a row, ending a sync,
marking a connection dead -- is acting on a message that explicitly refuses to
support it, and the case it gets wrong is a repository somebody made private
this morning.

## Sources

- Documentation: https://developer.atlassian.com/cloud/bitbucket/rest/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve bitbucket     # run it
cauldron verify bitbucket -v # check every claim
```
