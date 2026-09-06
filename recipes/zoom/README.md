# Zoom

Emulates the Zoom API (v2), for local development and tests.

**15 conformance cases, 3 checked against the live API on 2026-09-05.**

The meeting and recurrence cases still cite documentation, since a real meeting needs a real account. The credential and routing shapes needed no account at all, and checking them live confirmed this Recipe's existing claim and found one more.

## What this Recipe found

A meeting's numeric id and its UUID are not interchangeable, and using one where the other belongs gets a 404 that looks exactly like a missing meeting -- reports and recordings are keyed by UUID, meetings themselves by id, and the id is large enough that parsing it as a JavaScript number is unsafe. Meeting type is also an integer enum, and a recurring meeting with no fixed time (type 3) has no `start_time` at all, so code reading that field unconditionally finds nothing on exactly the meetings least likely to be tested.

## What checking it live found

This Recipe's `authentication_error` was already right -- for both an absent credential and a present, garbage one, which had never actually been checked before. A path nothing declares answers a different code and sentence again, resolved before any credential is judged.

## Sources

- Documentation: https://developers.zoom.us/docs/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve zoom     # run it
cauldron verify zoom -v # check every claim
```
