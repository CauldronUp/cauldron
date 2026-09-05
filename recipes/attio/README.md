# Attio

Emulates the Attio API (2.0.0), for local development and tests.

**8 conformance cases, 2 checked against the live API.**

Struck live against api.attio.com on 2026-09-05, and it caught something worse than a wrong message: the auth block declared no keys at all, which this project's own engine reads as "accepts anything". Every request to this emulator succeeded whether or not a credential was sent, for as long as this Recipe has existed. The two live cases fixed it and, along the way, recorded Attio's own wording for the two failures: an absent Authorization header is refused with a sentence about the header itself, and a present, wrong key is refused with a sentence naming the length a real one should be.

## What writing this Recipe changed

This Recipe was excluded as impossible, and the assessment was simply wrong.
Attio is REST and publishes an OpenAPI description; it had been recorded as
GraphQL and nobody reread the reason.

Both entries on that list sat unchallenged for the same reason: nothing rereads
a list of things that cannot be done, and an expired reason reads exactly like a
live one.

## Sources

- Documentation: https://api.attio.com/openapi/api
- Machine-readable description: https://attio.com/openapi.json, last checked 2026-09-02
  `cauldron drift attio` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve attio     # run it
cauldron verify attio -v # check every claim
```
