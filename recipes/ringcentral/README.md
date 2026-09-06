# RingCentral

Emulates the RingCentral API (v1.0), for local development and tests.

**21 conformance cases, 2 checked against the live API.**

Everything past the credential check still cites documentation rather than an observation, because reaching it means placing real calls. The credential check itself was verified directly against platform.ringcentral.com, unauthenticated, on 2026-09-05 -- exactly the production host the next paragraph explains this Recipe cannot go further into.

## What this Recipe found

Checked live: an absent credential and a fictitious one are different failures, and neither is the token-corruption case this file already modelled. No Authorization header answers `{"errorCode":"AGW-401","message":"Authorization header is not specified"}`. A syntactically fine but fictitious bearer answers a gateway code on the outside and a different one nested inside the errors array -- `{"errorCode":"TokenInvalid",...,"errors":[{"errorCode":"OAU-149","message":"OAuth token is invalid"}]}` -- the same sentence under two different machine codes depending on which one a client reads.

RingCentral retired its developer sandbox entirely on January 1, 2025, and now tells developers to build against production, which for a telephony API means a test run places real calls and sends real, billed messages to whoever is on the other end. This Recipe takes the position Mercury, Bill.com, Gusto and Deel take elsewhere in this collection, but further: nothing here dials or sends anything at all. Reading messages, call logs, and extensions is modelled; placing a call and sending a message are not.

A sent message being `Sent` does not mean delivered -- it moves Queued, then Sent, then Delivered, and `Sent` only means RingCentral handed it to a carrier, which can take minutes to disagree; the worse case is `DeliveryFailed` arriving after what already looked like a success. A deleted message stays in the message store with its content fully intact, marked `availability: Deleted`, so a listing that does not filter on that field shows messages the user believes are gone. `readStatus` and `messageStatus` are separate fields answering separate questions, so a message can be simultaneously `Read` and `DeliveryFailed`.

Rate limits are scoped per API group, Light, Medium, Heavy, rather than per account, so hammering one group leaves the others' budgets untouched and a single global counter in a client measures nothing real -- the 429 response names which group was exhausted.

## Sources

- Documentation: https://developers.ringcentral.com/api-reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve ringcentral     # run it
cauldron verify ringcentral -v # check every claim
```
