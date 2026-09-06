# Brevo

Emulates the Brevo API (v3), for local development and tests.

**17 conformance cases, 4 checked against the live API on 2026-09-01.**

## What this Recipe found

Three are Mailjet's. Three are Lithic's. Three are Moov's.
Mailjet is the one worth naming: **one HTTP 200 can carry a success and a
failure together**, a per-recipient verdict the status line cannot express.

## Sources

- Documentation: https://developers.brevo.com/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve brevo     # run it
cauldron verify brevo -v # check every claim
```
