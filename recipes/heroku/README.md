# heroku

Emulates the heroku API (3), for local development and tests.

**32 conformance cases, 4 checked against the live API on 2026-09-05.**

The resource cases cite documentation rather than an observation on a real account; the refusal cases were struck live, unauthenticated, against api.heroku.com.

## What this Recipe found

A successful list answers 206, not 200. Heroku pages with the `Range` header rather than query parameters, and a 206 with a `Next-Range` header means there's more to fetch -- only the final page answers 200. That makes two common client patterns wrong in opposite directions: checking `response.status === 200` treats every page but the last as a failure, while checking `response.ok` accepts the 206 and never reads `Next-Range`, silently processing only the first couple hundred records and reporting that as everything.

The `Accept` header isn't optional content negotiation, it's where the API version actually lives -- omit it and the response is 406, quoting the exact header that was needed. Errors are keyed by an `id` string (`"not_found"`) rather than a code, so code written to switch on `error.code` finds `undefined` on every single Heroku failure. And a formation scaled to quantity zero still exists as a process type -- scaling to nought doesn't remove it, so a client counting formations counts things that consume nothing, and one checking whether a process type exists gets the wrong answer about whether it's actually running.

The live probe found the declared authentication error had never actually been reachable: it was named `unauthorized`, which nothing in this file wired to a credential failure, so it sat unused while every real refusal fell through to a generic default. It is now named correctly, and a missing credential turns out to be a different sentence from a wrong one, both under the same `id`. An unrouted path and a wrong method on a real path both answer before authentication is ever consulted, with a different `not_found` sentence from the one this file already declares for a missing app.

## Sources

- Documentation: https://devcenter.heroku.com/articles/platform-api-reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve heroku     # run it
cauldron verify heroku -v # check every claim
```
