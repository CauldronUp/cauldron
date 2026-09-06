# Textline

Emulates the Textline API (v1), for local development and tests.

**6 conformance cases, 4 checked against the live API on 2026-09-01.**

## What this Recipe found

It has **no delivery-status vocabulary at all** -- a
Post carries only human triage fields, never anything about carriage.

## Sources

- Documentation: https://help.textline.com/en/articles/6660329-api-documentation
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve textline     # run it
cauldron verify textline -v # check every claim
```
