# Bunny Stream

Emulates the Bunny Stream API (v1), for local development and tests.

**14 conformance cases, 6 checked against the live API on 2026-09-03.**

## What this Recipe found

Stream's, where **the library id's shape is checked and its owner
is not**: numeric ids are indistinguishable, non-numeric ones are refused before
the credential.

## Sources

- Documentation: https://docs.bunny.net/reference/video_getvideo
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve bunnystream     # run it
cauldron verify bunnystream -v # check every claim
```
