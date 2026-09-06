# Docusign

Emulates the Docusign API (v2.1), for local development and tests.

**16 conformance cases, 4 checked against the live API on 2026-09-05.**

The resource shapes cite documentation rather than an observation on a real account; the refusal cases were struck live, unauthenticated, against demo.docusign.net.

## What writing this Recipe changed

It counts with strings rather than numbers. Together with Vonage's `"0"`, that
is what made the conformance checker compare a scalar's kind as well as its
value -- until then a case could pass whichever of the two the emulator
sent.

The live probe found that an absent credential and a wrong one are not the same failure: no Authorization header at all gets 401 PARTNER_AUTHENTICATION_FAILED, naming a missing integrator key, while a well-formed but invented bearer token gets the AUTHORIZATION_INVALID_TOKEN this file already declared. An unrouted path and a wrong method on a real path both answer an empty 404 before Docusign's application ever looks at authentication -- with a non-standard HTTP reason phrase this format has no way to reproduce.

## Sources

- Documentation: https://developers.docusign.com/docs/esign-rest-api/reference/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve docusign     # run it
cauldron verify docusign -v # check every claim
```
