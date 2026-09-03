# Vultr

Emulates the Vultr API (v2), for local development and tests.

**12 conformance cases, all of them checked against the live API on 2026-08-31.**

## What this Recipe found

**The total is the page size plus one.**
`meta.total` on `/v2/regions` reports 3 for `per_page=2` and 33 for
`per_page=32` -- a one-page lookahead, confirmed across six page sizes -- while
the other three catalogues on the same host report the real figure. Its `os` and
`applications` ids are JSON numbers where `regions` and `plans` are strings.

## Sources

- Documentation: https://www.vultr.com/api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve vultr     # run it
cauldron verify vultr -v # check every claim
```
