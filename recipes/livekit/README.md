# LiveKit

Emulates the LiveKit API (v1), for local development and tests.

**15 conformance cases, 11 checked against the live API on 2026-09-02.**

## What this Recipe found

**A bad credential outranks routing and an absent one
does not**, so sending nothing tells you more than sending something wrong.

## Sources

- Documentation: https://docs.livekit.io/home/server/managing-rooms/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve livekit     # run it
cauldron verify livekit -v # check every claim
```
