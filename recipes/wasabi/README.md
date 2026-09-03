# Wasabi

Emulates the Wasabi API (v1), for local development and tests.

**3 conformance cases, all of them checked against the live API on 2026-09-01.**

## What this Recipe found

Its S3 endpoint **is a marketing redirect until you try to
authenticate**. No headers at all gets a 303 to the company's home page in HTML;
send any Authorization value, however wrong, and the XML gate appears.

## Sources

- Documentation: https://docs.wasabi.com/docs/apis
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve wasabi     # run it
cauldron verify wasabi -v # check every claim
```
