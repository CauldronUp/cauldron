# Zoom

Emulates the Zoom API (v2), for local development and tests.

**11 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A meeting's numeric id and its UUID are not interchangeable, and using one where the other belongs gets a 404 that looks exactly like a missing meeting -- reports and recordings are keyed by UUID, meetings themselves by id, and the id is large enough that parsing it as a JavaScript number is unsafe. Meeting type is also an integer enum, and a recurring meeting with no fixed time (type 3) has no `start_time` at all, so code reading that field unconditionally finds nothing on exactly the meetings least likely to be tested.

## Sources

- Documentation: https://developers.zoom.us/docs/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve zoom     # run it
cauldron verify zoom -v # check every claim
```
