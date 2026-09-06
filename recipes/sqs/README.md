# AWS SQS

Emulates the AWS SQS API (2012-11-05), for local development and tests.

**18 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Receiving a message doesn't remove it. It becomes invisible for the visibility timeout, and if the consumer doesn't delete it in time it comes back -- code that treats receive as consume processes the same message repeatedly, and only under load, which is exactly when nobody is watching. `ApproximateReceiveCount` is the only signal that a message has been seen before, and the word "Approximate" is a warning rather than modesty.

An empty queue is also a trap: `ReceiveMessage` returning no messages is a 200 with no `Messages` key at all, not an empty array, so `for (const m of response.Messages)` throws on an idle queue instead of just doing nothing.

## Sources

- Documentation: https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve sqs     # run it
cauldron verify sqs -v # check every claim
```
