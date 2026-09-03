# Kaleyra

Emulates the Kaleyra API (v1), for local development and tests.

**5 conformance cases, 4 checked against the live API on 2026-09-01.**

## What this Recipe found

**One identifier has three names**: a parent id, that
id with `:1` appended, and the same value called campaign_id on the status route.

## Sources

- Documentation: https://developers.kaleyra.io/docs/sms-api-overview
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve kaleyra     # run it
cauldron verify kaleyra -v # check every claim
```
