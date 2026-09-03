# Sardine

Emulates the Sardine API (v1), for local development and tests.

**10 conformance cases, all of them checked against the live API on 2026-09-01.**

## What this Recipe found

Its **decision shape is sales-gated**, so what a risk API
returns could not be established at all -- and whose refusal echoes a mistyped
client id back.

## Sources

- Documentation: https://docs.sardine.ai/guides/api-reference/overview
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve sardine     # run it
cauldron verify sardine -v # check every claim
```
