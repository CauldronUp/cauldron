# Giphy

Emulates the Giphy API (v1), for local development and tests.

**6 conformance cases, 3 checked against the live API on 2026-09-01.**

## What this Recipe found

**A failure arrives with an empty results array.** A
401 carries `"data": []` -- the same key a successful listing fills -- so a
client that iterates `data`, which is the obvious thing to do with a listing,
gets zero results and no exception, and renders an empty grid where it should be
reporting that nobody is authenticated.

## Sources

- Documentation: https://developers.giphy.com/docs/api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve giphy     # run it
cauldron verify giphy -v # check every claim
```
