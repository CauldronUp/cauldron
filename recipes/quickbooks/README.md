# QuickBooks

Emulates the QuickBooks API (v3), for local development and tests.

**8 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A QuickBooks list response is nested two levels deep under a capitalized, resource-specific key -- a customer list arrives as `QueryResponse.Customer`, so code reading `data` or `customers` at either the top level or the first level finds nothing. Failures arrive under `Fault.Error` as an array, and the HTTP status can still be 200 for a batch request, so code reading a top-level error message finds nothing there either. `SyncToken`, the optimistic-concurrency check on every record, is a string of digits -- sending a stale one is refused, and sending none at all silently overwrites whatever changed since the record was last read.

## Sources

- Documentation: https://developer.intuit.com/app/developer/qbo/docs/api/accounting/most-commonly-used/account
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve quickbooks     # run it
cauldron verify quickbooks -v # check every claim
```
