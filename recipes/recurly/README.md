# Recurly

Emulates the Recurly API (v2021-02-25), for local development and tests.

**15 conformance cases, 2 checked against the live API.**

Everything past the version gate and the credential check still cites documentation rather than an observation, because reaching it needs a real site. The credential check itself was verified directly against v3.recurly.com on 2026-09-05.

## What this Recipe found

Recurly versions its API through the Accept header rather than a path segment, and an ordinary request that does not name one of a dozen exact `application/vnd.recurly.vYYYY-MM-DD` values is refused with a 406 before the Authorization header is even looked at -- checked live, this runs ahead of authentication entirely and is described rather than encoded, since this format has no mechanism for a header whose value must match a set. Once a valid version is sent, an absent credential and a wrong one are different sentences under the identical `unauthorized` type: "You must provide a valid Authorization header" versus "Invalid API key". This file had one message for both.

A canceled Recurly subscription is not an expired one -- cancelling only stops the renewal, and the customer keeps access until the current period actually ends, so code that revokes access the moment cancellation happens takes away something the customer already paid for. Money is an integer in minor units with the currency alongside it, and the amount on the subscription itself is not necessarily the amount on the invoice once a discount or proration gets involved -- the two numbers can legitimately disagree for the same billing event.

## Sources

- Documentation: https://recurly.com/developers/api/v2021-02-25/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve recurly     # run it
cauldron verify recurly -v # check every claim
```
