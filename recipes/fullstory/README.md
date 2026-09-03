# FullStory

Emulates the FullStory API (v2), for local development and tests.

**6 conformance cases, all of them checked against the live API on 2026-09-01.**

## What this Recipe found

**There is no route to fetch one session**.
`/v2/sessions/{id}` answers the same generic 404 pure nonsense gets, while the
events beneath it answer 401 -- so the sub-resource exists and its parent does
not.

## Sources

- Documentation: https://developer.fullstory.com/server/getting-started/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve fullstory     # run it
cauldron verify fullstory -v # check every claim
```
