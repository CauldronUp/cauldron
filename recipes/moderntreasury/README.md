# Modern Treasury

Emulates the Modern Treasury API (v1), for local development and tests.

**14 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Modern Treasury's ledger is the actual product, and the trap is that summing amounts across it always gives zero -- every movement is two entries with equal amounts and opposite directions that cancel by design, so a naive sum that looks like a bug is actually the system working correctly. `direction` is credit or debit, but what that does to a balance depends on the account: an account whose `normal_balance` is debit goes up on a debit and down on a credit, and a credit-normal account is the reverse, so the same word means opposite things on two accounts in the same transaction.

A ledger account carries three balances that disagree on purpose: `pending` includes entries not yet posted, `posted` includes only the ones that are, and `available` is the only one safe to spend against. Reading the wrong one authorizes money that is not actually there. Nothing is ever edited -- a posted transaction is immutable, and a correction is a new transaction in the opposite direction, so the history is append-only and a balance is a sum over everything rather than a field anyone updates directly.

Payment orders, the half of Modern Treasury that actually instructs a bank, are deliberately not modelled, the same call Mercury and Bill.com make elsewhere in this collection. Idempotency keys are accepted and ignored, which the header calls out as a gap worth naming loudly for a ledger API specifically.

## Sources

- Documentation: https://docs.moderntreasury.com/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve moderntreasury     # run it
cauldron verify moderntreasury -v # check every claim
```
