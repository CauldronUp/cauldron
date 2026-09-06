# Airbyte

Emulates the Airbyte API (v1), for local development and tests.

**13 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

It **does not call a partial failure a success**:
`incomplete` is a status distinct from both succeeded and failed. It is also
undocumented -- present in the enum, absent from Airbyte's own prose.

## Sources

- Documentation: https://reference.airbyte.com
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve airbyte     # run it
cauldron verify airbyte -v # check every claim
```
