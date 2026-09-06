# TalkJS

Emulates the TalkJS API (v1), for local development and tests.

**12 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

Its conversation **exists and cannot be seen**: created by a
caller-chosen PUT, and invisible until somebody sends a message.

## Sources

- Documentation: http://talkjs.com/docs/REST_API/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve talkjs     # run it
cauldron verify talkjs -v # check every claim
```
