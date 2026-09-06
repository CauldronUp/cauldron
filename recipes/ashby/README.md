# Ashby

Emulates the Ashby API (v1), for local development and tests.

**21 conformance cases, 6 checked against the live API on 2026-09-02.**

## What this Recipe found

Its **own auth documentation is wrong about its API**: it
promises 401 for missing and 403 for wrong, and five credential shapes all
answered one 401 in plain text.

## Sources

- Documentation: https://developers.ashbyhq.com/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve ashby     # run it
cauldron verify ashby -v # check every claim
```
