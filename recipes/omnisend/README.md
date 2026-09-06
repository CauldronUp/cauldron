# Omnisend

Emulates the Omnisend API (2026-03-15), for local development and tests.

**13 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

Its ordinary create is an upsert -- "Existing contact
updated (upsert by email)" -- with **no guard-rail language anywhere**.

## Sources

- Documentation: https://api-docs.omnisend.com/reference/contacts
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve omnisend     # run it
cauldron verify omnisend -v # check every claim
```
