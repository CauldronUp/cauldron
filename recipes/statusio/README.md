# Status.io

Emulates the Status.io API (1.0), for local development and tests.

**10 conformance cases, all of them checked against the live API on 2026-08-31.**

## What this Recipe found

**A malformed id and an unknown one are the same
200**, both carrying an `error` key. Drop the id entirely and AWS API Gateway
answers instead, `{"message": "Missing Authentication Token"}`, on an API that
needs no authentication. Components use `id` and incidents use `_id`, and three
numeric scales reuse 100/200/300/400 for different meanings.

## Sources

- Documentation: https://kb.status.io/developers/public-status-api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve statusio     # run it
cauldron verify statusio -v # check every claim
```
