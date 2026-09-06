# SignNow

Emulates the SignNow API (v2), for local development and tests.

**13 conformance cases, 9 checked against the live API on 2026-09-01.**

## What this Recipe found

Its **document status is not a field at all**: its SDK
marks it `[JsonIgnore]` and computes it on the client from invite records, so
there is nothing on the wire to be honest about.

## Sources

- Documentation: https://docs.signnow.com/docs/signnow/welcome
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve signnow     # run it
cauldron verify signnow -v # check every claim
```
