# Plivo

Emulates the Plivo API (v1), for local development and tests.

**9 conformance cases, 3 checked against the live API on 2026-09-01.**

## What this Recipe found

**Whose failure its own SDK cannot parse** -- a 401 in
`text/html` on a JSON API, where the vendor's published client is written to
read `{"error": ...}` and substitutes a generic sentence when parsing throws.

## Sources

- Documentation: https://www.plivo.com/docs/messaging/api/overview
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve plivo     # run it
cauldron verify plivo -v # check every claim
```
