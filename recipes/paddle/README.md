# Paddle

Emulates the Paddle API (v1), for local development and tests.

**9 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What writing this Recipe changed

Its fixture carries a subscription with a cancellation scheduled but not yet
applied: still active, already ending, and absent from any client that treats
cancellation as a single moment.

## Sources

- Documentation: https://developer.paddle.com/api-reference/overview
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve paddle     # run it
cauldron verify paddle -v # check every claim
```
