# Mangopay

Emulates the Mangopay API (v2.01), for local development and tests.

**8 conformance cases, 5 checked against the live API on 2026-09-02.**

## What this Recipe found

It **blocks a replayed idempotency key for a day**
whatever the first call did, where Checkout.com forgets it and Stripe replays the
failure.

## Sources

- Documentation: https://docs.mangopay.com/api-reference/payins/view-payin
- Machine-readable description: https://docs.mangopay.com/openapi.json, last checked 2026-09-02
  `cauldron drift mangopay` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve mangopay     # run it
cauldron verify mangopay -v # check every claim
```
