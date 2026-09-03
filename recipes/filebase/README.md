# Filebase

Emulates the Filebase API (v1), for local development and tests.

**3 conformance cases, all of them checked against the live API on 2026-09-01.**

## What this Recipe found

Its **opaque trace header names the pod that served
you**. `X-Amz-Id-2` took two values across a dozen requests, `czNndy0x` and
`czNndy0y`, which are base64 for `s3gw-1` and `s3gw-2`.

## Sources

- Documentation: https://docs.filebase.com
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve filebase     # run it
cauldron verify filebase -v # check every claim
```
