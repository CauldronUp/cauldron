# heroku

Emulates the heroku API (3), for local development and tests.

**12 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A successful list answers 206, not 200. Heroku pages with the `Range` header rather than query parameters, and a 206 with a `Next-Range` header means there's more to fetch -- only the final page answers 200. That makes two common client patterns wrong in opposite directions: checking `response.status === 200` treats every page but the last as a failure, while checking `response.ok` accepts the 206 and never reads `Next-Range`, silently processing only the first couple hundred records and reporting that as everything.

The `Accept` header isn't optional content negotiation, it's where the API version actually lives -- omit it and the response is 406, quoting the exact header that was needed. Errors are keyed by an `id` string (`"not_found"`) rather than a code, so code written to switch on `error.code` finds `undefined` on every single Heroku failure. And a formation scaled to quantity zero still exists as a process type -- scaling to nought doesn't remove it, so a client counting formations counts things that consume nothing, and one checking whether a process type exists gets the wrong answer about whether it's actually running.

## Sources

- Documentation: https://devcenter.heroku.com/articles/platform-api-reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve heroku     # run it
cauldron verify heroku -v # check every claim
```
