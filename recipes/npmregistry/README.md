# npm registry

Emulates the npm registry API (registry), for local development and tests.

**14 conformance cases, 8 checked against the live API on 2026-08-22.**

## What this Recipe found

What was checked is the shape each case turns on rather than the values: that a
per-version document carries no time, no versions and no dist-tags, that
`deprecated` is a string and is absent rather than false when it does not
apply, that `integrity` sits inside `dist`, and that a missing package answers
404 with `{"error":"Not found"}`. The Recipe records which real packages those
were observed on, so the dates can be audited rather than taken on trust.

## Sources

- Documentation: https://github.com/npm/registry/blob/master/docs/REGISTRY-API.md
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve npmregistry     # run it
cauldron verify npmregistry -v # check every claim
```
