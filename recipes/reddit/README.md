# Reddit

Emulates the Reddit API (v1), for local development and tests.

**5 conformance cases, 4 checked against the live API on 2026-09-02.**

## What this Recipe found

It is **mostly a wall and says so** -- every post endpoint
blocked on all three hosts, so the identifier scheme is cited and never
fixtured.

## Sources

- Documentation: https://github.com/reddit-archive/reddit/wiki/OAuth2
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve reddit     # run it
cauldron verify reddit -v # check every claim
```
