# Google Address Validation

Emulates the Google Address Validation API (v1), for local development and tests.

**13 conformance cases, 4 checked against the live API on 2026-09-05.**

The resource cases cite protocol buffer definitions rather than an observation on a real project, because calling this API needs a billed key; the refusal cases were struck live, unauthenticated, against addressvalidation.googleapis.com.

## What this Recipe found

There is no field that says whether an address is valid. The verdict is spread across three separate granularities and five booleans, none of them called `valid`, `correct`, `deliverable` or `ok` -- what counts as valid is a policy the caller has to decide, not something the API states. Worse, this is proto3 JSON, so every `false` and every zero is omitted from the response entirely: a perfect address comes back with a verdict that has almost nothing in it, and a broken one has more fields, not fewer, so reading the presence of a key gives the opposite impression from what's actually true.

`OTHER` is the granularity that reads like a shrug and actually means undeliverable -- it's documented as every non-deliverable granularity bucketed together. And Google will silently correct part of an address and only mention it in a boolean nested several levels down inside one component (`replaced`), with the corrected value sitting in `formattedAddress` and no top-level flag announcing that anything changed.

Field names on the wire are also camelCase while the API is published as a snake_case proto -- `address_complete` becomes `addressComplete` -- so a client written straight from the proto definition reads nothing at all.

The live probe found the declared authentication error had never actually been reachable: it was named `invalid_key`, which nothing in this file wired to a credential failure, so every unauthenticated request fell through to a generic default instead. It is now named correctly, its content was already right, and a missing key entirely turns out to be a different failure altogether -- 403 `PERMISSION_DENIED`, naming an unregistered caller rather than a bad key.

## Sources

- Documentation: https://developers.google.com/maps/documentation/address-validation
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve googleaddress     # run it
cauldron verify googleaddress -v # check every claim
```
