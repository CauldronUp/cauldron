# Ably

Emulates the Ably API (v1), for local development and tests.

**10 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Publishing to Ably tells a caller nothing: the response is 204 with no body at all, not even an empty object, so calling .json() on it throws and there's no message id to correlate anything with afterwards. History comes back newest first by default, so code that iterates the array expecting chronological order replays a conversation backwards and only notices when somebody reads it.

Two smaller traps: a token's capability is a JSON string, not an object, so it looks like structured data and reading token.capability["*"] finds undefined until it's parsed a second time. And errors carry a five-digit code and an HTTP status separately, and they're different numbers -- 40400 is not 404, and support threads quote the long one.

This covers only the REST surface; Cauldron serves no WebSocket or realtime transport, and token requests are not signed or verified -- only the key is checked.

## Sources

- Documentation: https://ably.com/docs/api/rest-api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve ably     # run it
cauldron verify ably -v # check every claim
```
