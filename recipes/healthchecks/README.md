# Healthchecks

Emulates the Healthchecks API (v3), for local development and tests.

**12 conformance cases, all of them checked against the live API on 2026-08-31.**

## What this Recipe found

**The failure the service exists to detect
cannot be modelled at all.** It notices a dead job by observing that *no ping
arrived* -- a verdict reached by a background process and announced to channels
the client never sees. Every mechanism here is request-in, response-out, so
there is no request to send for the thing the service is for. What is
observable: a ping to a UUID it has never seen answers `200 OK (not found)`,
fourteen bytes saying both things at once, and `/fail` answers identically -- so
signalling failure to a check that was never created is indistinguishable from
signalling it to one that was.

## Sources

- Documentation: https://healthchecks.io/docs/http_api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve healthchecks     # run it
cauldron verify healthchecks -v # check every claim
```
