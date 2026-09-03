# Wayback Machine

Emulates the Wayback Machine API (waybackmachine), for local development and tests.

**7 conformance cases, all of them checked against the live API on 2026-08-28.**

## What this Recipe found

Machine's, where asking whether a page existed in 1990
answers yes, from 2002. `closest` means closest available, `available` is `true`,
`status` is `"200"`, and nothing in the response says the snapshot is twelve
years from the date that was asked for. The distance is computable, but only by
subtracting two timestamps in two different formats, and no field states it.

**And `timestamp` is in the document twice, meaning two different things.** The
one inside `closest` is when the snapshot was taken, fourteen digits. The one at
the top level is the date the caller asked for, eight digits, echoed back -- and
absent entirely when the caller sent none. So the same key is a request
parameter in one place and an answer in the other, at two lengths, one of them
conditional. `url` is doubled the same way. `status` is the string `"200"`, so
`=== 200` is false and `== 200` is true. A URL nothing has archived answers
`archived_snapshots: {}` with a 200, where the hits have `closest` -- and a
request with no `url` at all answers the plain text `Error: no url parameter`,
also with a 200, in `text/html`, from an endpoint whose successes are JSON.

## Sources

- Documentation: https://archive.org/help/wayback_api.php
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve waybackmachine     # run it
cauldron verify waybackmachine -v # check every claim
```
