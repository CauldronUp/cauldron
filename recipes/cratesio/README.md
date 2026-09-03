# crates.io

Emulates the crates.io API (v1), for local development and tests.

**11 conformance cases, all of them checked against the live API on 2026-08-25.**

## What this Recipe found

Written the same way and for the same reason: no
OpenAPI, no credential, so every claim was checked live. Four fields say which
version is the latest -- `max_version`, `newest_version`, `max_stable_version`
and `default_version` -- and on `serde` all four read 1.0.229, which is why
nobody notices. On `rand` the field called `newest_version` is 0.8.8 while
`max_version` is 0.10.2: newest means most recently published and max means
highest by semver, and a patch to an old line went up on 25 August against a
0.10.2 from 2 July. A request with no `User-Agent` is a 403 that has nothing to
do with permission -- "we ask that your user agent actually identify your bot"
-- on a registry that needs no credential at all. And `id` is the crate's name
on a crate and an integer on a version, in one response.

## Sources

- Documentation: https://crates.io/data-access
- Machine-readable description: https://crates.io/api/openapi.json, last checked 2026-08-31
  `cauldron drift cratesio` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve cratesio     # run it
cauldron verify cratesio -v # check every claim
```
