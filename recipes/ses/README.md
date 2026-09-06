# AWS SES

Emulates the AWS SES API (v2), for local development and tests.

**13 conformance cases, 2 checked against the live API.**

Everything past the credential check still cites documentation rather than an observation, because reaching it needs real AWS credentials. The credential check itself was verified directly against email.us-east-1.amazonaws.com, unsigned, on 2026-09-05.

## What this Recipe found

`GetAccount` was modelled as a listing and is not one. This Recipe served `{"Accounts": [{...}]}` for `GET /v2/email/account`, and AWS's own service model describes a `GetAccountResponse` whose members sit at the top level -- there is no collection and no identifier. Two of those members were flat here as well, where AWS nests them: `Max24HourSend` and `SentLast24Hours` belong under `SendQuota`. So code written against this fake read `Accounts[0].SendingEnabled` and would have found `undefined` against SES, and code written against SES read `SendQuota.Max24HourSend` and found `undefined` here. Both are fixed, and the case that had asserted the invented shape now asserts the real one.

Checked live: a request with no Authorization header at all answers `{"message":"Missing Authentication Token"}` with the machine type carried in a response header, `x-amzn-ErrorType`, not in the body at all. This REST-routed AWS service puts its error type somewhere different from a JSON-1.1 one like Secrets Manager, and this file had assumed the body convention here too, pairing the right status with the wrong code and the wrong sentence.

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
