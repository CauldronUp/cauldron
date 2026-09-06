# Rollbar

Emulates the Rollbar API (1), for local development and tests.

**16 conformance cases, 5 checked against the live API.**

Everything past the credential and routing checks still cites documentation rather than an observation, because reaching it needs a real project. Those checks were verified directly against api.rollbar.com, unauthenticated, on 2026-09-05.

## What this Recipe found

Checked live: no access token at all and a wrong one are different sentences, "access token required" versus "invalid access token". Routing runs ahead of the credential too: an unrouted path and a wrong method on a real path both answer "not found" needing nothing sent at all.

Every response is wrapped in `{"err": 0, "result": {...}}`, and `err` is a number, not a boolean -- zero means success. A client that checks for the presence of `err` rather than its value treats every successful call as a failure, and one that reads the body directly finds nothing at all.

An item's level and its status are independent, too: a critical error somebody has already resolved is still returned, still critical. Counting by level alone re-reports work that is already done.

## Sources

- Documentation: https://docs.rollbar.com/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve rollbar     # run it
cauldron verify rollbar -v # check every claim
```
