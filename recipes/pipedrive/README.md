# Pipedrive

Emulates the Pipedrive API (v1), for local development and tests.

**10 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Every Pipedrive response carries a `success` boolean, and the HTTP status is not the whole story -- a client that branches on status code alone can miss a failure, and one that reads `data` without first checking `success` reads null. A deal's `status` and its `stage` are also different things: a deal can sit in the final pipeline stage and still be open, so code that treats reaching the last stage as a win reports revenue that has not actually closed.

## Sources

- Documentation: https://developers.pipedrive.com/docs/api/v1
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve pipedrive     # run it
cauldron verify pipedrive -v # check every claim
```
