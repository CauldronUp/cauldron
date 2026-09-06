# Vonage

Emulates the Vonage API (1), for local development and tests.

**16 conformance cases, 4 checked against the live API on 2026-09-05.**

Vonage has no test mode for sending, so the send-path cases still cite documentation. The credential shape needed no account at all, and checking it live found two different provider APIs living in one Recipe, each with a failure shape that had never been run.

## What writing this Recipe changed

It reports a successful send with the string `"0"`, which is what made the
conformance checker compare a scalar's kind as well as its value. Before that,
`"0"` and `0` were indistinguishable and this Recipe's case passed whichever the
emulator sent.

## What checking it live found

The SMS API (`/sms/json`) answers a credential failure as this Recipe's own 200-status message-count envelope, with codes `2` (missing) and `4` (wrong) -- not the 401 and single sentence this Recipe had claimed. The Account API (`/account/numbers`) answers something else again: RFC 7807 problem details, `{type, title, detail}`, on a 422 for an absent credential and a 401 for a present, wrong one -- the shape this Recipe's `authentication_error` had actually been modelling, just at the wrong status with the wrong sentence. The two are now told apart with a per-route `auth` block. A third route, `/search/message`, answered a bare 404 regardless of credential and could not be checked further -- left as documentation.

## Sources

- Documentation: https://developer.vonage.com/en/api/sms
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve vonage     # run it
cauldron verify vonage -v # check every claim
```
