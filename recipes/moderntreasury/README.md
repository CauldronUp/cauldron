# Modern Treasury

Emulates the Modern Treasury API (v1), for local development and tests.

**25 conformance cases, 3 checked against the live API.**

Struck live 2026-09-05 against app.moderntreasury.com, no account and no key -- and found three distinct real failures where this file declared one, "Invalid credentials", never exercised by any case. A missing credential gets a JSON sentence naming the auth scheme; a secret carrying neither the `live-` nor the `test-` prefix is refused before its value is read at all, quoting that rule back; and a correctly-prefixed key nobody issued gets no body at all, zero bytes. All three are served. The middle one is why `auth.shape` exists: this Recipe was one of four that could describe a shape check and not model it, because the only mechanism available accepted every correctly-shaped string outright.

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
