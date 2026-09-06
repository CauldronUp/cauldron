# Axiom

Emulates the Axiom API (v1), for local development and tests.

**8 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

**Whose own guide documents an endpoint that does not exist**.
The path its ingest guide names 404s; the one that answers is labelled
"(Legacy)" in its generated reference.

## Sources

- Documentation: https://axiom.co/docs/restapi/introduction
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve axiom     # run it
cauldron verify axiom -v # check every claim
```
