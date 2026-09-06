# Paddle

Emulates the Paddle API (v1), for local development and tests.

**20 conformance cases, 3 checked against the live API.**

Struck live 2026-09-05 against api.paddle.com, no account and no key. The status and code this file declared for a missing credential were already right, 403 / `forbidden`; the sentence was invented, and the `documentation_url` this file shares across every failure is really specific to each one. Fixed below. A credential Paddle does not recognise at all is a different failure again -- same 403, code `authentication_malformed`, its own sentence and its own `documentation_url` -- so a client branching on the status cannot tell it from sending nothing. That one is served too. What could not be established is the well-formed-but-wrong case: Paddle's key layout is public, and a string built to exactly that layout is still malformed to it, so something beyond the shape is checked and no key this Recipe can invent gets past it.

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
