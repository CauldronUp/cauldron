# Sendcloud

Emulates the Sendcloud API (v2), for local development and tests.

**9 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

The call that buys a shipping label from a real carrier doesn't look like a purchase. `POST /parcels` creates a parcel, and a boolean on that same request, `request_label`, decides whether it also buys a label from the carrier -- defaulted to `true` in every quickstart example. There's no separate "buy" step to forget; there's a field to forget to set to `false`, on the call that reads as "create a record." The field never even appears on the response it produced.

A cancel is not a refund, and the consequence is stated in units of insurance rather than money: cancelling a parcel that has already shipped voids its insurance coverage without stopping the carrier from moving it. And this host turns out to be two different pieces of software behind one hostname -- an absent Authorization header reaches Sendcloud's own app and gets its own 401, while any credential at all, right or wrong or malformed, is intercepted first by a gateway in front of it, answering a differently shaped 401 even for a path that was never routed.

## Sources

- Documentation: https://sendcloud.dev/api/v2/parcels/index.md
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve sendcloud     # run it
cauldron verify sendcloud -v # check every claim
```
