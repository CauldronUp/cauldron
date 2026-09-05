# Resend

Emulates the Resend API (v1), for local development and tests.

**11 conformance cases, 2 checked against the live API.**

Everything past the credential check still cites documentation rather than an observation, because reaching it needs a real account. The credential check itself was verified directly against api.resend.com on 2026-09-05.

## What this Recipe found

Checked live: a missing credential and a wrong one are not just different sentences, they are different HTTP statuses. No Authorization header answers 401 `{"name":"missing_api_key","message":"Missing API Key"}`; a syntactically fine but fictitious bearer answers 400, not 401, `{"name":"validation_error","message":"API key is invalid"}` -- a key that is definitely wrong is, on the wire, a bad request rather than an authentication failure.

Resend's `last_event` is exactly that, the last event, not a history, so a message that was delivered and later complained about reads `"complained"` with the delivery outcome simply gone, and there is no other endpoint that remembers it happened. A dashboard built by polling this field for delivered messages undercounts permanently rather than temporarily. A send itself returns only an id, meaning accepted rather than delivered, the same shape and the same available mistake as SES.

An API key is either `full_access` or `sending_access`, and a sending-only key trying to read domains is refused with a message about permission rather than about the key being invalid, so an integration built and tested with a full-access key in development can fail on a call nobody exercised once it is deployed with a restricted one.

## Sources

- Documentation: https://resend.com/docs/api-reference
- Machine-readable description: https://resend.com/openapi.json, last checked 2026-09-05
  `cauldron drift resend` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve resend     # run it
cauldron verify resend -v # check every claim
```
