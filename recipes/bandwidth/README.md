# bandwidth

Emulates the bandwidth API (v2), for local development and tests.

**16 conformance cases, 2 checked against the live API.**

Two were struck live against messaging.bandwidth.com on 2026-09-05, and found an absent credential and a wrong one are not remotely the same answer. This file had a present, wrong Basic credential getting the same JSON error body an absent header gets; the live host actually sends no body at all for a wrong credential -- Content-Length: 0 -- with a WWW-Authenticate challenge header instead, and reserves the JSON {"type", "description"} shape for when nothing was sent at all. The description text was also wrong: "Authentication failed" was never the sentence; it is "Your request could not be authenticated".

## What this Recipe found

The same Bandwidth message comes back under entirely different property names depending on which endpoint returns it: a send response calls it id, from, to and direction: "out"; searching for that same message a minute later calls it messageId, sourceTn, destinationTn and messageDirection: "OUTBOUND" -- not one name survives, and even the direction string changes. A client reconciling a sent message against a later search has to translate between two vocabularies for the same object, with no hint from either response that the other vocabulary exists.

A send also answers 202, not 200 or 201 -- accepted means queued, nothing has reached a handset yet, and the real delivery outcome arrives later at a callback URL, so code treating any 2xx as "sent" reports success for messages that were never delivered. "to" stays an array even for a single recipient, so one send to three numbers is one message id but three separate delivery callbacks. And segmentCount isn't derivable from the text yourself -- one non-GSM character (a pasted curly quote) drops the per-segment limit from 160 characters to 70, so a message that was one segment yesterday is three today at three times the price.

## Sources

- Documentation: https://dev.bandwidth.com/apis/messaging
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve bandwidth     # run it
cauldron verify bandwidth -v # check every claim
```
