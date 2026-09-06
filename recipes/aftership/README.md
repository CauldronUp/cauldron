# AfterShip

Emulates the AfterShip API (v4), for local development and tests.

**11 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A tracking number alone doesn't identify a parcel -- AfterShip keys every tracking by carrier slug and number together, both in the path, so code that stores the number alone and looks it up later may retrieve a different carrier's tracking for someone else's parcel. The response envelope changes shape with cardinality too: a listing is under data.trackings, a single tracking under data.tracking, so a client with one response helper reads undefined from whichever shape it wasn't written for.

Status is two fields and only one is safe to branch on: tag is a coarse, fixed-set state, while subtag is finer and grows over time (AfterShip adds new suffixes like InTransit_001), so code that switches on subtag breaks when they add one. Delivered is also not terminal -- a delivered parcel can move to Exception afterwards (refused, returned, damaged) -- so code that stops polling on Delivered stops watching exactly when something goes wrong. Creating the same tracking twice is refused outright rather than treated as idempotent, so a retry after a timeout can look like the tracking was never created.

## Sources

- Documentation: https://www.aftership.com/docs/tracking/quickstart/api-quick-start
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve aftership     # run it
cauldron verify aftership -v # check every claim
```
