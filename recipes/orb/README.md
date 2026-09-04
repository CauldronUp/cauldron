# Orb

Emulates the Orb API (v1), for local development and tests.

**12 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

An Orb invoice in `draft` status is not a bill, it is a running total -- `amount_due` moves every time a new usage event arrives, and the only invoice guaranteed not to change is one that has been issued. Nothing about the object announces that difference except the status word itself, and the recurring failure is someone reading `amount_due`, showing it to a customer or writing it into their own ledger, and finding a different number an hour later with nothing having gone wrong.

Even an issued invoice is not fully final: a backfill can amend usage that has already been invoiced, so a closed period can be superseded by a new invoice carrying a different total for the same dates. Money is sent as a decimal string (`"1284.50"`, not `128450` or `1284.5`) specifically because floating point cannot be trusted with a bill, and reporting the same usage event twice is sometimes free and sometimes a double charge -- deduplication is by idempotency key within a window, so a retry inside the window costs nothing and the same retry an hour later bills twice.

Usage aggregation and the write path that reports events are not modelled -- the amounts here are whatever a fixture puts there rather than a sum Cauldron performs, reproducing the shapes and the state transitions but not the underlying arithmetic.

## Sources

- Documentation: https://docs.withorb.com/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve orb     # run it
cauldron verify orb -v # check every claim
```
