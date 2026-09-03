# Browserless

Emulates the Browserless API (rest), for local development and tests.

**8 conformance cases, 6 checked against the live API on 2026-09-01.**

## What this Recipe found

It **has no field for the target's status at all** --
no header, no body field, so a page that 404d and a page that rendered are the
same successful response. Sending the credential as an `Authorization` header,
the way nearly every other provider here takes it, crashes the edge with a 500.

## Sources

- Documentation: https://docs.browserless.io/rest-apis/content
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve browserless     # run it
cauldron verify browserless -v # check every claim
```
