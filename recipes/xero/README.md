# Xero

Emulates the Xero API (2.0), for local development and tests.

**10 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What writing this Recipe changed

It answers a request for one invoice with a list of one, so code that expects a
single object finds an array where it did not look for one.

## Sources

- Documentation: https://developer.xero.com/documentation/api/accounting/overview
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve xero     # run it
cauldron verify xero -v # check every claim
```
