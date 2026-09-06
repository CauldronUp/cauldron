# Bringg

Emulates the Bringg API (oauth2), for local development and tests.

**10 conformance cases, 8 checked against the live API on 2026-09-01.**

## What this Recipe found

Its public API **leaks an untranslated i18n key**: Rails'
own placeholder for a string nobody wrote, in a live error field.

## Sources

- Documentation: https://developers.bringg.com/docs/bringg-api-access-management
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve bringg     # run it
cauldron verify bringg -v # check every claim
```
