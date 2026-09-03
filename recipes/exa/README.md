# Exa

Emulates the Exa API (unversioned), for local development and tests.

**10 conformance cases, 8 checked against the live API on 2026-09-03.**

## What this Recipe found

**The one search API that labels a result cached**, per
document, where its two neighbours label nothing.

## Sources

- Documentation: https://docs.exa.ai/reference/search
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve exa     # run it
cauldron verify exa -v # check every claim
```
