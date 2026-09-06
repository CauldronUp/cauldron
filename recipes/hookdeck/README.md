# Hookdeck

Emulates the Hookdeck API (2024-03-01), for local development and tests.

**20 conformance cases, 3 checked against the live API on 2026-09-05.**

The resource cases cite Hookdeck's own OpenAPI description rather than an observation on a real project; the refusal cases were struck live, unauthenticated, against api.hookdeck.com.

## What this Recipe found

"Did my webhook work" has three different answers here, and they're allowed to disagree. Hookdeck models one arriving webhook as three separate objects: a Request (what arrived), an Event (what gets delivered, with six possible states), and an Attempt (one try, with only two states). A request can be received, verified, and answered with a 200, and still produce nothing downstream at all -- if nothing was configured to receive it, `events_count` on the request is simply zero, and the original sender never sees anything wrong.

The states also don't mean what their names suggest across layers: an event whose last two delivery attempts both failed is still `SCHEDULED`, because the next retry hasn't run yet -- `FAILED` on an event specifically means retries are exhausted, a different sentence entirely from `FAILED` on an attempt. Code that alerts on event status stays silent through every retry until the very end; code that alerts on individual attempts fires on failures that go on to succeed.

Thirty-two of thirty-three documented attempt error codes mean the request never reached your server at all (DNS failures, connection resets, TLS errors) -- only one, `BAD_RESPONSE`, means your server actually answered and said no. An attempt that never reached you carries `response_status: null` rather than a 5xx, so code that counts failures by checking `response_status >= 500` misses every transport failure entirely.

The live probe found authentication failures here are 12 bytes of plain text, "Unauthorized", not the JSON envelope this file declared and every other failure in this API actually uses. An unrouted path and a wrong method on a real path both answer Express's own default HTML 404 before authentication is ever consulted, naming the exact request that hit it.

## Sources

- Documentation: https://hookdeck.com/docs/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve hookdeck     # run it
cauldron verify hookdeck -v # check every claim
```
