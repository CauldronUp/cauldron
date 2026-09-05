# Pusher

Emulates the Pusher API (v1), for local development and tests.

**13 conformance cases, 6 checked against the live API.**

Everything about publishing, channel state, and pagination still cites documentation rather than an observation, because reaching it needs a real app. The credential and routing checks were verified directly against api-mt1.pusher.com, unauthenticated, on 2026-09-05.

## What this Recipe found

An ordinary unauthenticated request does not get a crafted refusal at all, checked live: it crashes with sixteen bytes of plain text, "expected string" -- a parameter-parsing failure surfacing as the response, not a message about a key or a signature, and true whether the request carries no auth_key or a fictitious one. The signature-mismatch sentence this file already modelled is a different, deeper failure that needs a real app secret to produce. Routing also needs no credential: an unrouted path and a wrong method both answer plain-text "404 NOT FOUND".

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
