# SendGrid

Emulates the SendGrid API (v3), for local development and tests.

**10 conformance cases, 3 checked against the live API.**

Everything past the credential check still cites documentation rather than an observation, because reaching it needs a real account. The credential check itself was verified directly against api.sendgrid.com on 2026-09-05.

## What this Recipe found

Checked live: a missing credential and a wrong one are different sentences. No Authorization header answers "Permission denied, wrong credentials"; a syntactically fine but fictitious bearer answers "The provided authorization grant is invalid, expired, or revoked" -- the message this file already had exactly right, just not paired with what an absent header actually says. A wrong method on a real path needs no credential at all and answers a 405 naming the one method that is allowed; an unrouted path answers a third sentence again. SendGrid checks the credential and the route in an order that depends on which mistake was made, which is recorded in the file header rather than encoded as a case.

A successful send answers 202 with an empty body -- the message id lives in a header, not the JSON. A client that calls `.json()` on that response throws, and an emulator that helpfully invents a body would hide the bug entirely. Failures come back as an array, because a single request can be wrong in more than one way at once, so code written for a single error object breaks on the first request that has two problems.

## Sources

- Documentation: https://www.twilio.com/docs/sendgrid/api-reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve sendgrid     # run it
cauldron verify sendgrid -v # check every claim
```
