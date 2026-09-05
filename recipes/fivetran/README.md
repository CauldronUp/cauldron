# fivetran

Emulates the fivetran API (v1), for local development and tests.

**15 conformance cases, 4 checked against the live API on 2026-09-05.**

The resource cases cite documentation rather than an observation on a real account; the refusal cases were struck live, unauthenticated, against api.fivetran.com.

## What this Recipe found

Triggering a sync that's already running answers 200 and doesn't start a new one, and nothing in the response says so. It's the identical success shape as actually starting a sync, so a client that triggers and then waits for `succeeded_at` to change is waiting on a sync it never started, and a pipeline that fires on every webhook calls this endpoint a hundred times and syncs once.

A connector also carries three separate state fields that don't imply each other: `setup_state`, `sync_state` and `update_state` can independently read broken, syncing, and on-schedule at the same time, and no single one of them means the connector is actually working. `paused` and `sync_frequency` are likewise unrelated -- a connector that isn't paused but has no schedule never runs, and a paused one keeps its schedule and simply ignores it. Every response, including failures, comes back `{"code": "Success", "data": {...}}`-shaped, with the real outcome carried in the `code` string rather than in the HTTP status alone.

The live probe found this file's own authentication error wrong on both fields: real Fivetran answers `AuthFailed`, not `Unauthorized`, and an absent credential and a present-but-wrong one share that code while carrying different sentences ("Missing authorization header" against "Invalid authorization credentials"), neither of which this file had declared. It also found an unrouted path answers a distinct 404 before authentication is ever consulted, while a wrong method on a real path gets the same auth refusal a missing credential does, not a 405.

## Sources

- Documentation: https://fivetran.com/docs/rest-api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve fivetran     # run it
cauldron verify fivetran -v # check every claim
```
