# Klarna

Emulates the Klarna API (v1), for local development and tests.

**6 conformance cases, 4 checked against the live API on 2026-09-02.**

## What this Recipe found

It **tells you not to trust its own status field**, in a
sentence from its own reference.

## Sources

- Documentation: https://docs.klarna.com/payments/after-payments/order-management/manage-orders-with-the-api/view-and-change-orders/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve klarna     # run it
cauldron verify klarna -v # check every claim
```
