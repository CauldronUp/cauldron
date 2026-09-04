# Google Pub/Sub

Emulates the Google Pub/Sub API (v1), for local development and tests.

**10 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Every Pub/Sub message body arrives base64-encoded, and reading `message.data` directly gets a string that looks plausible rather than obviously wrong -- it has a length, it will sit happily in a database, and nothing about it screams "decode me first." This is, per the header, the single most reproduced mistake against this API. Pulling nothing back returns an object with no `receivedMessages` key at all, so a loop written as `for (const m of response.receivedMessages)` throws on the quietest possible input, the same shape SQS has, and not a coincidence.

`deliveryAttempt` only appears on a message at all when the subscription has a dead-letter policy configured -- on every other subscription it is simply absent, missing exactly where a consumer trying to detect a retry would look for it. An `ackId` is scoped per-pull, not per-message, with a default ack deadline of ten seconds -- any handler doing real work will exceed that, and the message gets redelivered while the first attempt is still in flight.

## Sources

- Documentation: https://cloud.google.com/pubsub/docs/reference/rest
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve pubsub     # run it
cauldron verify pubsub -v # check every claim
```
