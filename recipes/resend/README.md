# Resend

Emulates the Resend API (v1), for local development and tests.

**9 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Resend's `last_event` is exactly that, the last event, not a history, so a message that was delivered and later complained about reads `"complained"` with the delivery outcome simply gone, and there is no other endpoint that remembers it happened. A dashboard built by polling this field for delivered messages undercounts permanently rather than temporarily. A send itself returns only an id, meaning accepted rather than delivered, the same shape and the same available mistake as SES.

An API key is either `full_access` or `sending_access`, and a sending-only key trying to read domains is refused with a message about permission rather than about the key being invalid, so an integration built and tested with a full-access key in development can fail on a call nobody exercised once it is deployed with a restricted one.

## Sources

- Documentation: https://resend.com/docs/api-reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve resend     # run it
cauldron verify resend -v # check every claim
```
