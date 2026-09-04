# Recurly

Emulates the Recurly API (v2021-02-25), for local development and tests.

**8 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A canceled Recurly subscription is not an expired one -- cancelling only stops the renewal, and the customer keeps access until the current period actually ends, so code that revokes access the moment cancellation happens takes away something the customer already paid for. Money is an integer in minor units with the currency alongside it, and the amount on the subscription itself is not necessarily the amount on the invoice once a discount or proration gets involved -- the two numbers can legitimately disagree for the same billing event.

## Sources

- Documentation: https://recurly.com/developers/api/v2021-02-25/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve recurly     # run it
cauldron verify recurly -v # check every claim
```
