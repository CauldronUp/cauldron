# Ably

Emulates the Ably API (v1), for local development and tests.

**18 conformance cases, 3 checked against the live API.**

Most cases here cite documentation rather than an observation, and the Recipe's own header says why for the ones that still do. Three were struck live against rest.ably.io on 2026-09-05 with no app at all, and they corrected a mistake: every auth failure in this file had been modelled as 40100 "Invalid credentials", and the live host answers an entirely absent credential with 40101 "No authentication information provided" instead. 40100 turns out to be what a doc page says a wrong secret gets for a real app, and this project has no real, paid Ably app to check that half against.

## What this Recipe found

Publishing to Ably tells a caller nothing: the response is 204 with no body at all, not even an empty object, so calling .json() on it throws and there's no message id to correlate anything with afterwards. History comes back newest first by default, so code that iterates the array expecting chronological order replays a conversation backwards and only notices when somebody reads it.

Two smaller traps: a token's capability is a JSON string, not an object, so it looks like structured data and reading token.capability["*"] finds undefined until it's parsed a second time. And errors carry a five-digit code and an HTTP status separately, and they're different numbers -- 40400 is not 404, and support threads quote the long one.

Struck live: a method Ably does not accept on a real path, and a path it has never heard of, both come back 404 code 40400 with identical wording. Ably does not appear to answer 405 for anything this Recipe found.

This covers only the REST surface; Cauldron serves no WebSocket or realtime transport, and token requests are not signed or verified -- only the key is checked.

## Sources

- Documentation: https://ably.com/docs/api/rest-api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve ably     # run it
cauldron verify ably -v # check every claim
```
