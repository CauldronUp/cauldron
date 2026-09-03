# Docusign

Emulates the Docusign API (v2.1), for local development and tests.

**8 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What writing this Recipe changed

It counts with strings rather than numbers. Together with Vonage's `"0"`, that
is what made the conformance checker compare a scalar's kind as well as its
value -- until then a case could pass whichever of the two the emulator
sent.

## Sources

- Documentation: https://developers.docusign.com/docs/esign-rest-api/reference/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve docusign     # run it
cauldron verify docusign -v # check every claim
```
