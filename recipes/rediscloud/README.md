# Redis Cloud

Emulates the Redis Cloud API (v1), for local development and tests.

**16 conformance cases, 7 checked against the live API on 2026-09-02.**

## What this Recipe found

Cloud's, where **half a credential is worse than none** -- either
header alone answers a bare nginx 500 and only both together reach a real 401.

## Sources

- Documentation: https://redis.io/docs/latest/operate/rc/api/get-started/
- Machine-readable description: https://api.redislabs.com/v1/cloud-api-docs/capi, last checked 2026-09-02
  `cauldron drift rediscloud` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve rediscloud     # run it
cauldron verify rediscloud -v # check every claim
```
