# Pipedrive

Emulates the Pipedrive API (v1), for local development and tests.

**18 conformance cases, 1 checked against the live API.**

Everything past the credential check still cites documentation rather than an observation, because reaching it needs a company account. The credential check itself was verified directly against api.pipedrive.com on 2026-09-05.

## What this Recipe found

Checked live: no Authorization header, a fictitious bearer token, an unrouted path, and a wrong method all answer `{"success":false,"error":"unauthorized access","errorCode":401,"error_info":"Please check developers.pipedrive.com"}` -- the same body regardless of which of the four produced it. Two things this file had assumed rather than observed turned out wrong: the sentence was written as "You need to be authorized to make this request." rather than "unauthorized access", and `error_info` -- read here as a machine code -- is a constant human sentence pointing at the docs site, not something that varies with the failure.

Every Pipedrive response carries a `success` boolean, and the HTTP status is not the whole story -- a client that branches on status code alone can miss a failure, and one that reads `data` without first checking `success` reads null. A deal's `status` and its `stage` are also different things: a deal can sit in the final pipeline stage and still be open, so code that treats reaching the last stage as a win reports revenue that has not actually closed.

## Sources

- Documentation: https://developers.pipedrive.com/docs/api/v1
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve pipedrive     # run it
cauldron verify pipedrive -v # check every claim
```
