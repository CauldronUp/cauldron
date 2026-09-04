# Gusto

Emulates the Gusto API (v1), for local development and tests.

**9 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

`processed` means submitted, not paid. A payroll with `processed: true` has been sent, but the money doesn't move until `check_date`, which can be days later, and the payroll can still be cancelled up until the deadline. Treating `processed` as "these people have been paid" is wrong for the entire window that actually matters -- and a payroll carries three different dates (pay period start/end, submission deadline, check date), so reporting off the wrong one puts earnings in the wrong month, which is a tax problem rather than a display bug.

A terminated employee is still returned by the API and still carries payroll history -- correct to exclude from a new payroll run, wrong to exclude from a year-end report, and it's the same list serving both uses. And an employee can exist, be listed, and still be impossible to actually pay: `onboarded: false` means the record is incomplete, with nothing else on the object explaining why.

## Sources

- Documentation: https://docs.gusto.com/embedded-payroll/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve gusto     # run it
cauldron verify gusto -v # check every claim
```
