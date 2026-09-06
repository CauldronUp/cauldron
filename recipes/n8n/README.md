# n8n

Emulates the n8n API (v1), for local development and tests.

**18 conformance cases, 16 checked against the live API on 2026-09-01.**

## What this Recipe found

Its **credential message depends on which route you missed**,
splitting not by read against write nor by resource but by which helper each
handler happened to be written with.

## Sources

- Documentation: https://docs.n8n.io/api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve n8n     # run it
cauldron verify n8n -v # check every claim
```
