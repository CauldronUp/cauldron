# OMDb

Emulates the OMDb API (v1), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-09-03.**

## What this Recipe found

**The success-shaped failure that no longer is** -- the
reputation was real and is out of date, though the envelope still holds the
string False rather than a boolean.

## Sources

- Documentation: https://www.omdbapi.com/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve omdb     # run it
cauldron verify omdb -v # check every claim
```
