# Toast

Emulates the Toast API (v2), for local development and tests.

**10 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A restaurant's day doesn't end at midnight. `businessDate` is the day the money belongs to and `openedDate` is when the order actually happened, and they disagree for a few hours every night -- an order rung in at 1am Saturday belongs to Friday's business date. A report grouped by the calendar date in `openedDate` instead splits every late night across two days, and the total for the month is still right, which is exactly the shape of bug that survives longest. `businessDate` also isn't a date at all -- it's the integer `20260201`, which sorts correctly and parses as a date in nothing.

An order also has no total of its own: money lives on its checks, one check per bill, so a table split three ways is one order with three checks, and anything summing orders has to walk two levels down first.

## Sources

- Documentation: https://doc.toasttab.com/openapi/orders
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve toast     # run it
cauldron verify toast -v # check every claim
```
