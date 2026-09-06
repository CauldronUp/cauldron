# DNSimple

Emulates the DNSimple API (v2), for local development and tests.

**17 conformance cases, 9 checked against the live API on 2026-09-03.**

## What this Recipe found

**An unmentioned scheme crashes the server** -- Basic,
which its documentation never names, answers 500 where everything else is 401.

## Sources

- Documentation: https://developer.dnsimple.com/v2/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve dnsimple     # run it
cauldron verify dnsimple -v # check every claim
```
