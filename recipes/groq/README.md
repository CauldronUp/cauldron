# Groq

Emulates the Groq API (v1), for local development and tests.

**11 conformance cases, 5 checked against the live API on 2026-08-31.**

## What this Recipe found

**The wrong verb and the wrong path are the same
failure.** `POST` to a real GET-only path answers the identical 404 a path
nobody defined answers, with the method folded into a sentence rather than into
a status -- there is no 405 anywhere. Its routes also resolve *before* the
credential, which is the opposite order from Scaleway's, written the same day.

## Sources

- Documentation: https://console.groq.com/docs/api-reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve groq     # run it
cauldron verify groq -v # check every claim
```
