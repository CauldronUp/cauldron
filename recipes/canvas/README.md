# Canvas

Emulates the Canvas API (v1), for local development and tests.

**14 conformance cases, 11 checked against the live API on 2026-09-03.**

## What this Recipe found

It is an **unauthenticated existence oracle**: an
allocated course id answers 401 and an unallocated one 404, with no credential
either way, and a deleted course is indistinguishable from one that never was.

## Sources

- Documentation: https://canvas.instructure.com/doc/api/courses.html
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve canvas     # run it
cauldron verify canvas -v # check every claim
```
