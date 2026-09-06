# Scout APM

Emulates the Scout APM API (0.1), for local development and tests.

**18 conformance cases, 4 checked against the live API on 2026-09-01.**

## What this Recipe found

APM's, whose **performance group has no id at all** -- an Endpoint
is addressed by name and a Trace, the ephemeral thing, has a plain integer. Its
error pair is the mirror: a group id that names no get-one route.

## Sources

- Documentation: https://scoutapm.com/docs/api/
- Machine-readable description: https://scoutapm.com/api/v0/openapi.yaml, last checked 2026-09-01
  `cauldron drift scoutapm` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve scoutapm     # run it
cauldron verify scoutapm -v # check every claim
```
