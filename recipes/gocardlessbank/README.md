# gocardlessbank

Emulates the gocardlessbank API (v2), for local development and tests.

**13 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

This is the API formerly called Nordigen, not the GoCardless direct-debit API -- same company, different product, different base URL and credentials entirely, which is why the Recipe is named after the product rather than the company.

Booked and pending transactions arrive as two arrays in one response, not two endpoints, and the same purchase appears in `pending` first and `booked` later under a different `transactionId` -- code that merges the two arrays counts it twice. That id is also optional: some banks never send one, and for those the only way to identify a transaction is date, amount and remittance text together, so deduplicating on `transactionId` alone silently keeps every duplicate from those banks. And `transactionAmount.amount` is a signed string -- a debit is a negative string -- so parsing it as unsigned or comparing it as plain text both fail, in opposite directions.

A requisition's status is a two-letter code (`CR`, `GC`, `UA`, `RJ`, `SA`, `GA`, `LN`), not a boolean, and only `LN` means actually linked, with nothing in the response explaining what any of the codes mean.

## Sources

- Documentation: https://developer.gocardless.com/bank-account-data/overview
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve gocardlessbank     # run it
cauldron verify gocardlessbank -v # check every claim
```
