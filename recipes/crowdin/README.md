# Crowdin

Emulates the Crowdin API (v2), for local development and tests.

**16 conformance cases, 11 checked against the live API on 2026-09-02.** The 2 unchecked ones are the paging cases: they send the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

Its **missing project needs no credential to be
missing** -- a route naming an identifier resolves it first, while the
collection one segment up refuses the credential.

## Sources

- Documentation: https://support.crowdin.com/developer/api/v2/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve crowdin     # run it
cauldron verify crowdin -v # check every claim
```
