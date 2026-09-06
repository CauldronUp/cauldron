# Bitwarden

Emulates the Bitwarden API (1), for local development and tests.

**12 conformance cases, 10 checked against the live API on 2026-09-03.** The unchecked one is the paging case: it sends the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

Its **two refusals differ only in a header clients
ignore**: both 401 with zero bytes, separated by the challenge alone.

## Sources

- Documentation: https://bitwarden.com/help/public-api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve bitwarden     # run it
cauldron verify bitwarden -v # check every claim
```
