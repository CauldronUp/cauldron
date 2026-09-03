# Serper

Emulates the Serper API (unversioned), for local development and tests.

**9 conformance cases, 7 checked against the live API on 2026-09-01.**

## What this Recipe found

**Whose `statusCode` field only repeats the status line** --
so its missing-key and wrong-key failures share a status and a `statusCode`, and
differ only in whether the message ends with an invitation to sign up.

## Sources

- Documentation: https://serper.dev
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve serper     # run it
cauldron verify serper -v # check every claim
```
