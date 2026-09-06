# Mercury

Emulates the Mercury API (v1), for local development and tests.

**13 conformance cases, 2 checked against the live API.**

Struck live 2026-09-05 against api.mercury.com, no account and no key. This file declared one authentication_error, "Invalid or missing API token", for every failure; the real API sends two distinct codes and sentences -- `noAuthTokenHeader` / "No Authorization header present for this request" for a missing credential, `noTokenInDB` / "No matching token found" for one nobody issued. Split below.

## What this Recipe found

Mercury has no sandbox -- an API token reaches a real bank account, and a send in the wrong direction is a wire somebody has to phone about. Amounts are floating-point dollars, not integer cents, the opposite convention from Ramp elsewhere in this collection, so a shared ledger that assumes either format is wrong for half its inputs. An outgoing transaction is represented by a negative amount rather than a separate direction field, so summing a statement without checking signs gives a net figure where somebody wanted a total, and taking the absolute value erases the distinction entirely.

A pending transaction has no posted date and can still fail; its amount counts against the available balance rather than the current one, which is why the two balances disagree and why reading the wrong one authorizes spending against money that is already committed elsewhere. A failed transfer still shows up in the list with a status of failed and a reason -- the money never left, but summing without filtering on status counts it anyway.

Recipient and payment-creation endpoints are deliberately left out: this Recipe only reads accounts and transactions, on the reasoning that a Recipe which made sending money look easy would encourage exactly the risk the header warns about.

## Sources

- Documentation: https://docs.mercury.com/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve mercury     # run it
cauldron verify mercury -v # check every claim
```
