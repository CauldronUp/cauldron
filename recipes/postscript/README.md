# Postscript

Emulates the Postscript API (2.0), for local development and tests.

**15 conformance cases, 7 checked against the live API on 2026-09-01.**

## What this Recipe found

Its **created subscriber cannot be fetched** -- pending
until they confirm, with an id its own documentation says will not work and
confirmation arriving only as a webhook.

## Sources

- Documentation: https://developers.postscript.io/reference/get-subscriber
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve postscript     # run it
cauldron verify postscript -v # check every claim
```
