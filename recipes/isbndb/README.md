# ISBNdb

Emulates the ISBNdb API (v2), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-09-03.**

## What this Recipe found

Its **method rejection leaks its own internal route** --
a path segment that appears in no documentation and in nothing that was sent.

## Sources

- Documentation: https://isbndb.com/apidocs/v2
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve isbndb     # run it
cauldron verify isbndb -v # check every claim
```
