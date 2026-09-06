# courier

Emulates the courier API (v1), for local development and tests.

**16 conformance cases, 2 checked against the live API.**

Two were struck live against api.courier.com on 2026-09-05, and found the failure this file describes disagreeing with itself: the wire status is 403, not the 401 this file had, and the body carries its own "status" field reading 401 anyway -- a client reading that field gets a different answer than the response it was actually given. There is also no "type" field on the wire at all.

## What this Recipe found

Courier decides which channel a message goes out on at delivery time, based on the recipient's own preferences and routing rules set in a dashboard -- so a send only ever answers 202 with a requestId, and that requestId is not a message id: one request can become several messages, and the message ids only exist after routing happens, so a client storing the requestId to look up a message later has stored the wrong handle. Whether a message reached anyone at all is a separate lookup from whether and how it routed, and routing that finds no eligible channel isn't an error -- the message just quietly reaches status UNDELIVERABLE, having answered 202 like every successful send.

SIMULATED is also a real, distinct status: Courier can be told to route without actually sending, which is the only honest way to test routing logic, and a client that treats anything other than FAILED as delivered counts simulations as real deliveries. No routing is actually performed here -- which channel a message went out on is whatever the fixture says, which is the point, since the real answer only ever comes from a lookup rather than the send response itself.

## Sources

- Documentation: https://www.courier.com/docs/reference/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve courier     # run it
cauldron verify courier -v # check every claim
```
