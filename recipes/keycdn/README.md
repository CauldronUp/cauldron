# KeyCDN

Emulates the KeyCDN API (v1), for local development and tests.

**11 conformance cases, 5 checked against the live API on 2026-09-03.**

## What this Recipe found

**The clean no to its group's question** -- purge answers in
the past tense with nothing to poll.

## Sources

- Documentation: https://www.keycdn.com/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve keycdn     # run it
cauldron verify keycdn -v # check every claim
```
