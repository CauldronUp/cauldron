# Last.fm

Emulates the Last.fm API (2.0), for local development and tests.

**7 conformance cases, 5 checked against the live API on 2026-09-02.**

## What this Recipe found

**The one place two schemes name one thing**: an artist
resolves identically by name and by MusicBrainz identifier.

## Sources

- Documentation: https://www.last.fm/api/show/artist.getInfo
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve lastfm     # run it
cauldron verify lastfm -v # check every claim
```
