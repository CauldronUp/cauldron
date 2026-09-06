# Unsplash

Emulates the Unsplash API (v1), for local development and tests.

**9 conformance cases, 4 checked against the live API on 2026-09-01.**

## What this Recipe found

It **names its rate-limit headers without sending
them.** Every 401 carries an `access-control-expose-headers` line listing
`X-RateLimit-Limit` and `X-RateLimit-Remaining`, and neither header is present.
The case asserts their absence, so if that changes the suite notices.

## Sources

- Documentation: https://unsplash.com/documentation
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve unsplash     # run it
cauldron verify unsplash -v # check every claim
```
