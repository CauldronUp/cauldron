# Formstack

Emulates the Formstack API (documents), for local development and tests.

**9 conformance cases, 8 checked against the live API on 2026-09-01.**

## What this Recipe found

Its merge **answers nine bytes and sends the file
elsewhere**: `{"success":1}`, no id and no URL, with the document going to a
Deliveries configuration set up beforehand.

## Sources

- Documentation: https://www.webmerge.me/developers
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve formstack     # run it
cauldron verify formstack -v # check every claim
```
