# CircleCI

Emulates the CircleCI API (v2), for local development and tests.

**9 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What writing this Recipe changed

Its fixture carries a workflow waiting on approval -- neither a pass nor a
failure, and the state most likely to be missing from a client that branches on
success.

## Sources

- Documentation: https://circleci.com/docs/api/v2/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve circleci     # run it
cauldron verify circleci -v # check every claim
```
