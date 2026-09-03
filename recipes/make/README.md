# Make

Emulates the Make API (v2), for local development and tests.

**6 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

**A token from the wrong region looks like a token that
never existed** -- Make is partitioned by hostname, and using a real credential
on the wrong one produces the same "Not authorized." that a fabricated one does.

## Sources

- Documentation: https://developers.make.com/api-documentation/api-reference/scenarios
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve make     # run it
cauldron verify make -v # check every claim
```
