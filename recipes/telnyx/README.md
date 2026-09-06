# Telnyx

Emulates the Telnyx API (v2), for local development and tests.

**15 conformance cases, 3 checked against the live API on 2026-09-05.**

Telnyx has no sandbox, so the message and number cases still cite documentation. The credential shape needed no key at all, and checking it live found this Recipe's own message wrong.

## What this Recipe found

A message has no single status -- it has one per recipient, nested under a `to` array, so a message sent to three numbers can be delivered to two and undelivered to the third, and reading a top-level `status` finds nothing. The initial status is `queued`, and `delivery_failed`, if it comes, arrives minutes later on a webhook, so a 200 from the send is a receipt for the request rather than proof of anything. Cost is also absent from the response until the message finishes, so a spend tracker reading it right after the send sums nothing instead of zero.

## What checking it live found

No credential and a garbage one get different sentences under the same code (10009), and neither was `"The API key is invalid."` -- absence says so in as many words, and a credential that was sent is called "malformed" rather than "invalid" regardless of how it was wrong. A path nothing declares is a third failure again, `10005 Resource not found`, resolved before either credential question is asked.

## Sources

- Documentation: https://developers.telnyx.com/api
- Machine-readable description: https://telnyx.com/openapi.json, last checked 2026-08-31
  `cauldron drift telnyx` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve telnyx     # run it
cauldron verify telnyx -v # check every claim
```
