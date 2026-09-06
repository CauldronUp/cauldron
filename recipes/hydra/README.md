# hydra

Emulates the hydra API (v2), for local development and tests.

**10 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

An inactive token is a 200 with exactly one field in it. Introspecting a token that's expired, been revoked, or never existed answers `{"active": false}` -- not an error, not a 401, and not the token object with a flag on it: everything else (`sub`, `scope`, `client_id`, `exp`) is simply gone. `active` is the only required field in Hydra's own schema, which is the specification's way of saying exactly this. Code that checks the status code, then reads `result.sub` without checking `active` first, gets `undefined` and passes it straight through to whatever comes next.

`scope` is a space-separated string, not an array, so `scope.includes("admin")` returns true for a token holding only `"administrator"` -- any substring test on it is a privilege check that can pass for scopes nobody actually granted. `aud`, sitting right beside it, actually is an array, so a client that learned the shape of one field gets the other wrong. And `exp`, `iat` and `nbf` are all RFC 7519 seconds, not the millisecond timestamps most JavaScript expects, so a date built from one without multiplying by 1000 lands in January 1970.

A live check was attempted on 2026-09-05: Ory Hydra is self-hosted with no vendor-run public instance, and the one candidate found, Ory's own multi-tenant playground, resolves but answers a Cloudflare bot challenge on every Admin API path this Recipe models. The domain is real; nothing behind it was reachable.

## Sources

- Documentation: https://www.ory.sh/docs/hydra/reference/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve hydra     # run it
cauldron verify hydra -v # check every claim
```
