# Supabase

Emulates the Supabase API (v1), for local development and tests.

**8 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A project has two identifiers and only one of them works: `ref` is what every path takes, and `id` is documented "Deprecated: use `ref` instead" while still being sent beside it -- code that stores `id` and interpolates it into a path addresses nothing and gets a plain 404. Status also has fifteen values and none of them is `ACTIVE`; a healthy project reports `ACTIVE_HEALTHY`, so `status === 'ACTIVE'` is never true, and `status.startsWith('ACTIVE')` is true for `ACTIVE_UNHEALTHY`, which is a project that is not working.

Creating a project answers with no database object at all, because there's nothing to connect to yet -- reading `project.database.host` straight off the create response reads undefined.

## Sources

- Documentation: https://api.supabase.com/api/v1-json
- Machine-readable description: https://supabase.com/openapi.json, last checked 2026-08-31
  `cauldron drift supabase` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve supabase     # run it
cauldron verify supabase -v # check every claim
```
