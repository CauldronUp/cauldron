# Mastodon

Emulates the Mastodon API (v1), for local development and tests.

**12 conformance cases, 9 checked against the live API on 2026-09-02.** The 3 unchecked ones are the paging cases: they send the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

Its **id is local and whose uri is not**, and whose
reference and wire disagree about a missing status on the same server, the same
day.

## Sources

- Documentation: https://docs.joinmastodon.org/methods/statuses/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve mastodon     # run it
cauldron verify mastodon -v # check every claim
```
