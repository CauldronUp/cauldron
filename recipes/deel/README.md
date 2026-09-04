# Deel

Emulates the Deel API (v2), for local development and tests.

**11 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A Deel contract that reads in_progress is signed and actively running, with invoices accruing -- the word reads like a draft and means the opposite. waiting_for_client_sign is the state that actually means nothing has started, and it doesn't sort anywhere near in_progress alphabetically or logically. An employee and a contractor also come back from the same endpoint as genuinely different objects: a contractor has a contract and no employment record, an employee has both, so code that reads employment fields finds nothing on half a workforce.

Amounts are decimal strings with a separate currency field, and a contract's currency isn't the organisation's -- a team paying people in four currencies gets four different scales in one list, and summing without converting produces a number with no meaning. A terminated contract also keeps accruing invoices after the termination date, since final invoices arrive later, so a report cut off at the termination date silently drops the last one.

Contract creation and termination are deliberately not modelled -- an emulator that made hiring or firing somebody easy to exercise would be inviting the exact mistake this header warns about, the same restraint Mercury's payment creation takes.

## Sources

- Documentation: https://developer.deel.com/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve deel     # run it
cauldron verify deel -v # check every claim
```
