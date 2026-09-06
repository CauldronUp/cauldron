# Twitch

Emulates the Twitch API (helix), for local development and tests.

**8 conformance cases, 3 checked against the live API on 2026-08-31.**

## What this Recipe found

It **needs two credentials and checks one first.** Sending
a correct `Client-Id` with no token earns a byte-identical response to sending
nothing at all.

## Sources

- Documentation: https://dev.twitch.tv/docs/api/reference/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve twitch     # run it
cauldron verify twitch -v # check every claim
```
