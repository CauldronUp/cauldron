# Sumo Logic

Emulates the Sumo Logic API (v1), for local development and tests.

**13 conformance cases, 7 checked against the live API on 2026-09-01.**

## What this Recipe found

Logic's, whose search job **hands back a real handle and a warning
about its own failure mode**: an id, a link to poll, and a field explaining that
the poll depends on cookie affinity.

## Sources

- Documentation: https://help.sumologic.com/docs/api/search-job/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve sumologic     # run it
cauldron verify sumologic -v # check every claim
```
