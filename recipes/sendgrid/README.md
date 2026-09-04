# SendGrid

Emulates the SendGrid API (v3), for local development and tests.

**8 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A successful send answers 202 with an empty body -- the message id lives in a header, not the JSON. A client that calls `.json()` on that response throws, and an emulator that helpfully invents a body would hide the bug entirely. Failures come back as an array, because a single request can be wrong in more than one way at once, so code written for a single error object breaks on the first request that has two problems.

## Sources

- Documentation: https://www.twilio.com/docs/sendgrid/api-reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve sendgrid     # run it
cauldron verify sendgrid -v # check every claim
```
