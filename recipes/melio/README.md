# Melio

Emulates the Melio API (v1), for local development and tests.

**11 conformance cases, all of them checked against the live API on 2026-09-01.**

## What this Recipe found

A cancel route is **proven to exist without ever being
called**: three other verbs on that path answer "Could not find method X ... in
openapi spec" while POST reaches the credential gate.

## Sources

- Documentation: https://docs.melio.com/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve melio     # run it
cauldron verify melio -v # check every claim
```
