# Google Pub/Sub

Emulates the Google Pub/Sub API (v1), for local development and tests.

**18 conformance cases, 3 checked against the live API.**

Everything past the credential and routing checks still cites documentation rather than an observation, because reaching it needs a real project. Those checks were verified directly against pubsub.googleapis.com, unauthenticated, on 2026-09-05.

## What this Recipe found

No credential at all and a wrong one are not just different sentences, checked live -- they are different HTTP statuses. No Authorization header answers 403 PERMISSION_DENIED, refused for naming nobody; a syntactically fine but fictitious bearer answers 401 UNAUTHENTICATED instead, and the sentence this file already had was a truncation of the real one. An unrouted path never reaches Pub/Sub's own error handling at all -- it answers Google's generic front-end 404 HTML page, needing no credential either.

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
