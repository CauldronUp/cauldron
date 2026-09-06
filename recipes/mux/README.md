# Mux

Emulates the Mux API (v1), for local development and tests.

**20 conformance cases, 2 checked against the live API.**

Struck live 2026-09-05 against api.mux.com, no account and no key. This file declared one authentication_error, "Invalid credentials", for every failure; the real API sends two different sentences, neither of which is that -- a missing credential gets "This action must be completed through the dashboard interface.", and a secret nobody issued gets the plainer "Unauthorized request". Split below.

## What this Recipe found

An asset's id is not its playback id -- the playback URL is built from a separate object nested inside the asset, and using the asset id instead gives a 404 from a hostname that otherwise looks completely right, because both are opaque strings of similar length with nothing to tell them apart. A newly created, still-preparing asset has no playback IDs, no duration, and no aspect ratio at all -- not empty, not null, simply absent -- so code that reads `asset.playback_ids[0].id` immediately after creating throws on undefined rather than getting a helpful empty value.

A signed playback policy needs a JWT, and a signed playback ID looks identical to a public one -- the URL built from it just returns 403 instead of video, so an integration developed against a public asset breaks the day somebody switches the policy on the account. Failures also carry their message inside a plural array: `error.message` is undefined and the actual sentence lives at `error.messages[0]`, which is backwards from every other API in this collection.

Assets never become ready on their own here; which state an asset holds is whatever a fixture puts there, removing the real waiting rather than simulating it. Test mode has the same timing problem as production, so it does not actually help with the part of a Mux integration that is hard to get right.

## Sources

- Documentation: https://docs.mux.com/api-reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve mux     # run it
cauldron verify mux -v # check every claim
```
