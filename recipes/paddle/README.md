# Paddle

Emulates the Paddle API (v1), for local development and tests.

**10 conformance cases, 2 checked against the live API.**

Struck live 2026-09-05 against api.paddle.com, no account and no key. The status and code this file declared for a missing credential were already right, 403 / `forbidden`; the sentence was invented, and the `documentation_url` this file shares across every failure is really specific to each one. Fixed below. A third real shape, a credential in a form Paddle does not recognise at all (`authentication_malformed`), was also found and is stated rather than modelled, since no fixture key this Recipe could invent is verifiably Paddle-shaped.

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
