# CircleCI

Emulates the CircleCI API (v2), for local development and tests.

**11 conformance cases, 3 checked against the live API.**

Three were struck live against circleci.com on 2026-09-05, and the finding is a genuine trap: an absent Circle-Token answers real JSON, {"message":"You must log in first."}, mislabelled Content-Type text/plain, while a present, wrong token answers plain text with no JSON in it at all, "Invalid token provided." A client that calls .json() unconditionally survives the first case by accident and throws on the second. This file's own message, "Authentication failed.", was never a sentence CircleCI sends.



## What writing this Recipe changed

Its fixture carries a workflow waiting on approval -- neither a pass nor a
failure, and the state most likely to be missing from a client that branches on
success.

## Sources

- Documentation: https://circleci.com/docs/api/v2/
- Machine-readable description: https://circleci.com/api/v2/openapi.json, last checked 2026-09-05
  `cauldron drift circleci` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve circleci     # run it
cauldron verify circleci -v # check every claim
```
