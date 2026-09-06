# opsgenie

Emulates the opsgenie API (v2), for local development and tests.

**18 conformance cases, 3 checked against the live API.**

Struck live 2026-09-05 against api.opsgenie.com, no account and no key. This file's declared "Could not authenticate" matched exactly, word for word, on no credential, the wrong scheme word, and a UUID-shaped key nobody issued. A key that is not UUID-shaped is a different answer entirely -- 422, "Key format is not valid!" -- because Opsgenie reads a credential's shape before its value. That one is now served too, and it caught a real fault: this Recipe's own fixture key was not a UUID, so the fake accepted a key the provider refuses outright. The key is a UUID now.

## What this Recipe found

An alert and an incident look like the same thing and are not: an alert is a page, somebody's phone rings, an incident is a coordination record, a channel, a timeline, stakeholder updates. They have separate ids, separate lifecycles, and separate close endpoints. Closing one does not close the other, so a monitor recovers, the alert closes automatically, and the incident sits open for days with nobody looking at it.

Almost every write here is asynchronous, answering 202 with a request id that is not the alert's own id -- finding out whether the thing you asked for actually happened means fetching the request status separately, which for a moment still says processing. An alert can also be addressed three different ways, its id, a short human-readable `tinyId`, or your own `alias`, and which one a path means is decided by a query parameter, so the same URL addresses a different alert depending on `identifierType`. Creating an alert against an alias that is already open does not create a new one -- it increments the existing alert's count and answers with the same 202 as a real creation, so a flapping monitor produces one alert with a count of forty rather than forty separate alerts, and a client counting responses counts wrong.

## Sources

- Documentation: https://docs.opsgenie.com/docs/api-overview
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve opsgenie     # run it
cauldron verify opsgenie -v # check every claim
```
