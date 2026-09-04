# Pusher

Emulates the Pusher API (v1), for local development and tests.

**9 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Publishing to a Pusher channel answers with a bare empty object, `{}`, not an id, not a count, not a status, so there is nothing to correlate a later event against, and code that reads `response.id` off a successful call gets undefined. The channel list is a map keyed by channel name rather than an array, so looping over it as a list finds nothing, and a channel with no occupants is simply absent from the map rather than present with a zero.

`user_count` and `subscription_count` are different numbers -- one person with three browser tabs open is one user and three subscriptions -- so a dashboard reading the wrong one is off by however many tabs people happen to have open, and `user_count` only exists on presence channels at all. Failures are plain text, not JSON, the same trap Trello has elsewhere in this collection, so calling `.json()` on one throws instead of reporting the actual problem.

The `auth_signature` used to authorize private and presence channel subscriptions is an HMAC that is not verified here -- only the `auth_key` is checked, so a request with a valid key and a completely wrong signature is accepted, which is stated plainly rather than left for someone to discover.

## Sources

- Documentation: https://pusher.com/docs/channels/library_auth_reference/rest-api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve pusher     # run it
cauldron verify pusher -v # check every claim
```
