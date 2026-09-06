# Buildkite

Emulates the Buildkite API (v2), for local development and tests.

**12 conformance cases, 1 checked against the live API.**

Struck live against api.buildkite.com on 2026-09-05, both with no Authorization header and with a made-up Bearer token: byte-identical 401. This file had truncated the message to "Authentication required"; Buildkite sends a longer sentence with a link to its own auth docs.

## What this Recipe found

"blocked" is a real Buildkite build state meaning the build is waiting at a manual gate for a human to unblock it -- neither a pass nor a failure, and neither is "scheduled" or "not_run" -- so a deployment gate that treats anything but "passed" as failed rejects a build that's simply waiting on someone. A build that hasn't started has no started_at, and one that hasn't finished has no finished_at, so computing a duration without checking for either produces a date in 1970.

One fidelity gap: Buildkite pages with a Link header, which Cauldron doesn't model.

## Sources

- Documentation: https://buildkite.com/docs/apis/rest-api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve buildkite     # run it
cauldron verify buildkite -v # check every claim
```
