# Anthropic

Emulates the Anthropic API (2023-06-01), for local development and tests.

**21 conformance cases, 6 checked against the live API on 2026-09-02.**

## What this Recipe found

**One route resolves the id before the credential**
and its sibling on the same host does not.

## Sources

- Documentation: https://platform.claude.com/docs/en/api/messages
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve anthropic     # run it
cauldron verify anthropic -v # check every claim
```
