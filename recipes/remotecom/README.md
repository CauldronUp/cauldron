# Remote.com

Emulates the Remote.com API (v1), for local development and tests.

**16 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

It **says twice who employs somebody**: `type` and
`employment_model` are both on every record and neither is derivable from the
other.

## Sources

- Documentation: https://developer.remote.com/reference/welcome-to-remote-api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve remotecom     # run it
cauldron verify remotecom -v # check every claim
```
