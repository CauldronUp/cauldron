# calcom

Emulates the calcom API (v2), for local development and tests.

**14 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A Cal.com slot listing is a view, not a reservation, and nothing in the API says so -- you show someone a free slot, they pick it, and between those two moments somebody else booked it, so the booking request fails with a message that reads like a bug in your code and is actually the system working as designed. A booking also has two identifiers and the public one isn't the numeric one: id is Cal.com's internal integer, uid is what shows up in URLs and the attendee's cancellation link, and the paths take the uid.

Cancelling doesn't delete a booking -- it stays, its status becomes cancelled, and the slot frees, so code counting bookings counts cancelled ones too, and code expecting a 404 after cancelling gets a 200 instead. An event type requiring confirmation also produces pending bookings that hold the slot without being accepted or rejected either way -- the slot is gone regardless of which way it resolves.

No slot is ever actually taken out from under a reader here; the double-booking race is described by an error that can be armed on purpose, since reproducing it for real needs two clients and precise timing.

## Sources

- Documentation: https://cal.com/docs/api-reference/v2/introduction
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve calcom     # run it
cauldron verify calcom -v # check every claim
```
