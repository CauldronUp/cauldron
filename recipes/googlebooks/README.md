# Google Books

Emulates the Google Books API (v1), for local development and tests.

**10 conformance cases, 5 checked against the live API on 2026-09-03.**

## What this Recipe found

Books', where **sending a key fails faster than sending none**:
any key skips the exhausted anonymous quota and fails on a per-project check.

## Sources

- Documentation: https://developers.google.com/books/docs/v1/using
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve googlebooks     # run it
cauldron verify googlebooks -v # check every claim
```
