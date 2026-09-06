# Fireworks

Emulates the Fireworks API (v1), for local development and tests.

**10 conformance cases, 6 checked against the live API on 2026-09-01.**

## What this Recipe found

It **quotes back a path the caller never sent**: the
gateway strips its own prefix before composing "Path not found", so grepping your
logs for the path you were told about finds nothing.

## Sources

- Documentation: https://docs.fireworks.ai/api-reference/list-models
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve fireworks     # run it
cauldron verify fireworks -v # check every claim
```
