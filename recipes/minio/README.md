# MinIO

Emulates the MinIO API (v4), for local development and tests.

**7 conformance cases, all of them checked against the live API on 2026-09-01.**

## What this Recipe found

It **never examines the name you asked about**: its own
default root username and an invented one answer byte-identically.

## Sources

- Documentation: https://min.io/docs/minio/linux/reference/minio-mc-admin.html
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve minio     # run it
cauldron verify minio -v # check every claim
```
