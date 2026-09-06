# Milvus

Emulates the Milvus API (v2), for local development and tests.

**10 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

Its insert **looks the most finished of any here and is
not**: a real count, returned when the rows reach a message queue, under a
default consistency level its own documentation says allows inconsistency.

## Sources

- Documentation: https://docs.zilliz.com/reference/restful
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve milvus     # run it
cauldron verify milvus -v # check every claim
```
