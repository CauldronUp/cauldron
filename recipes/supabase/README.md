# Supabase

Emulates the Supabase API (v1), for local development and tests.

**12 conformance cases, 2 checked against the live API on 2026-09-05.**

Most of this Recipe still cites documentation, since projects and secrets need a real account. The credential shape needed no account at all, and checking it live found something worse than a wrong claim.

## What this Recipe found

A project has two identifiers and only one of them works: `ref` is what every path takes, and `id` is documented "Deprecated: use `ref` instead" while still being sent beside it -- code that stores `id` and interpolates it into a path addresses nothing and gets a plain 404. Status also has fifteen values and none of them is `ACTIVE`; a healthy project reports `ACTIVE_HEALTHY`, so `status === 'ACTIVE'` is never true, and `status.startsWith('ACTIVE')` is true for `ACTIVE_UNHEALTHY`, which is a project that is not working.

Creating a project answers with no database object at all, because there's nothing to connect to yet -- reading `project.database.host` straight off the create response reads undefined.

## What checking it live found

`auth` declared `scheme: bearer` with no keys and no pattern, which this project's own rule reads as "accepts anything" -- so nothing was checking the credential at all, and every one of this file's own cases would have passed with their `Authorization` header removed or replaced with garbage. `keys` now names the fixture token those cases already sent. Once that was fixed, the live probe found a real distinction to model: no header at all answers `{"message":"Unauthorized"}`, and a well-formed bearer value that is not a real access token answers `{"message":"JWT could not be decoded"}` instead.

## Sources

- Documentation: https://api.supabase.com/api/v1-json
- Machine-readable description: https://supabase.com/openapi.json, last checked 2026-08-31
  `cauldron drift supabase` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve supabase     # run it
cauldron verify supabase -v # check every claim
```
