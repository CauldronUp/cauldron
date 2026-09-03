# Webflow

Emulates the Webflow API (2.0.0), for local development and tests.

**8 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What writing this Recipe changed

This Recipe found a bug that had already shipped. Every timestamp field was
being filled in automatically, so a site that had never been published still
carried a `lastPublished`. The emulator was claiming events that never happened,
and no test written against it could have caught that -- the value was always
there, so nothing ever looked wrong.

## Sources

- Documentation: https://developers.webflow.com/data/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve webflow     # run it
cauldron verify webflow -v # check every claim
```
