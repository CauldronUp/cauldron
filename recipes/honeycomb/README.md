# Honeycomb

Emulates the Honeycomb API (v1), for local development and tests.

**10 conformance cases, 3 checked against the live API on 2026-08-31.**

## What this Recipe found

**One endpoint answers `problem+json` and six do
not.** This one was written to find out whether RFC 7807 needed a new error
style. It did not -- a flat style with the message under `title` reproduces it
exactly -- and that negative result is why the file exists in the shape it does.

## Sources

- Documentation: https://docs.honeycomb.io/api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve honeycomb     # run it
cauldron verify honeycomb -v # check every claim
```
