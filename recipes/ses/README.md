# AWS SES

Emulates the AWS SES API (v2), for local development and tests.

**11 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A successful send means accepted, not delivered. `SendEmail` answers 200 with a `MessageId`, and that's a promise to try -- the bounce, if there is one, arrives minutes later on a completely separate channel, so code that treats the 200 as delivery reports a hundred percent success rate regardless of what actually happened.

The suppression list makes this worse in the other direction: an address that has bounced is added automatically, and a later send to it is accepted with a `MessageId` and then silently dropped. No error, no bounce, no event -- the only way to know is to check the list, and almost nothing does.

## Sources

- Documentation: https://docs.aws.amazon.com/ses/latest/APIReference-V2/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve ses     # run it
cauldron verify ses -v # check every claim
```
