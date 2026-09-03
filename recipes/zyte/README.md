# Zyte

Emulates the Zyte API (v1), for local development and tests.

**7 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

It **draws the line the other two blur**: a target that
answered badly is a 200 wrapping the real status, and a failure to reach the
target at all is a genuine non-2xx -- 521 permanent, 520 ban, 421 unresolvable.

## Sources

- Documentation: https://docs.zyte.com/zyte-api/usage/errors.html
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve zyte     # run it
cauldron verify zyte -v # check every claim
```
