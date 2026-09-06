# Whereby

Emulates the Whereby API (v1), for local development and tests.

**9 conformance cases, 6 checked against the live API on 2026-09-01.**

## What this Recipe found

**The URL and the room expire together**, which is why
the room URL is unsigned where PDFMonkey's and Api2Pdf's are not.

## Sources

- Documentation: https://docs.whereby.com/reference/whereby-rest-api-reference/meetings
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve whereby     # run it
cauldron verify whereby -v # check every claim
```
