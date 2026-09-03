# Sinch

Emulates the Sinch API (v1), for local development and tests.

**7 conformance cases, 2 checked against the live API on 2026-09-01.**

## What this Recipe found

**Whose documented API host has no public address**:
`api.sinch.com` resolves through public DNS to three unroutable `10.65.x.x`
addresses, and the SMS product actually answers elsewhere.

## Sources

- Documentation: https://developers.sinch.com/docs/sms/api-reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve sinch     # run it
cauldron verify sinch -v # check every claim
```
