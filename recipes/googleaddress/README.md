# Google Address Validation

Emulates the Google Address Validation API (v1), for local development and tests.

**8 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

There is no field that says whether an address is valid. The verdict is spread across three separate granularities and five booleans, none of them called `valid`, `correct`, `deliverable` or `ok` -- what counts as valid is a policy the caller has to decide, not something the API states. Worse, this is proto3 JSON, so every `false` and every zero is omitted from the response entirely: a perfect address comes back with a verdict that has almost nothing in it, and a broken one has more fields, not fewer, so reading the presence of a key gives the opposite impression from what's actually true.

`OTHER` is the granularity that reads like a shrug and actually means undeliverable -- it's documented as every non-deliverable granularity bucketed together. And Google will silently correct part of an address and only mention it in a boolean nested several levels down inside one component (`replaced`), with the corrected value sitting in `formattedAddress` and no top-level flag announcing that anything changed.

Field names on the wire are also camelCase while the API is published as a snake_case proto -- `address_complete` becomes `addressComplete` -- so a client written straight from the proto definition reads nothing at all.

## Sources

- Documentation: https://developers.google.com/maps/documentation/address-validation
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve googleaddress     # run it
cauldron verify googleaddress -v # check every claim
```
