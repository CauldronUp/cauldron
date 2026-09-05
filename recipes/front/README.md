# Front

Emulates the Front API (1), for local development and tests.

**14 conformance cases, 4 checked against the live API on 2026-09-05.**

The resource cases cite documentation rather than an observation on a real inbox; the refusal cases were struck live, unauthenticated, against api2.frontapp.com.

## What this Recipe found

Front's envelope keys are all underscore-prefixed: a collection sits under `_results`, paging under `_pagination`, the error under `_error`. Code written by analogy with any other provider reaches for `data`, `results`, or `error` and finds nothing there. The envelope shape was already right; the live probe found the contents wrong. A missing credential gets "JSON Web Token error" and a present-but-wrong one gets "Invalid token" -- two different sentences under the identical "Unauthenticated" title, neither the "Unauthorized" / "The API token is invalid or missing" this file had declared. Front checks the credential before it looks at the path or the method at all, so an unrouted path and a wrong method get the same refusal a missing token does, not a 404 or a 405.

Timestamps are floating-point Unix seconds, not integers and not strings -- parsing them as integers drops the sub-second precision Front actually uses to order messages within a conversation. And like several shared-inbox providers in this collection, Front has no sandbox: a test that replies to a conversation sends a real email to a real customer, so the states most worth testing -- an archived conversation, an internal-only comment -- are the ones nobody rehearses.

## Sources

- Documentation: https://dev.frontapp.com/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve front     # run it
cauldron verify front -v # check every claim
```
