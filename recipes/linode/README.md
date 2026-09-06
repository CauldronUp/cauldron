# Linode

Emulates the Linode API (v4), for local development and tests.

**17 conformance cases, 13 checked against the live API on 2026-08-31.** The unchecked one is the paging case: it sends the two parameter names this Recipe declares, which is a claim about the provider read from its own description rather than struck against it.

## What this Recipe found

They are here to be read beside Vultr's. Two
competitors, the same job, and nothing in common at the envelope: Linode wraps
everything under a generic `data` with `page`/`pages`/`results`, Vultr names
each collection after itself and pages by cursor. Their failures share no field
at all. Linode's timestamps are not RFC 3339 -- `2025-10-01T04:00:00`, no zone.

## Sources

- Documentation: https://techdocs.akamai.com/linode-api/reference/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve linode     # run it
cauldron verify linode -v # check every claim
```
