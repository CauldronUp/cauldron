# marqeta

Emulates the marqeta API (v3), for local development and tests.

**22 conformance cases, 3 checked against the live API.**

Struck live 2026-09-05 against sandbox-api.marqeta.com, no account and no key. This file declared the 401 as code "401" / "Unauthorized"; the real one, identical whether no credential is sent, a wrong one, or the two halves swapped, is code "401001" / "Invalid credential has been detected". Fixed below.

## What this Recipe found

Marqeta has a real sandbox, but the reason to emulate it anyway is that a card payment is not one event, it is two: an authorization that holds money and a separate clearing, days later, for a possibly different amount, linked back to the authorization rather than replacing it. A ledger built by summing transactions counts the same payment twice, and a restaurant tab that authorizes the bill and clears the bill plus a tip is normal, not a bug.

There are three balances and only one is spendable: `available_balance` is what can actually be spent, `ledger_balance` includes holds that have not cleared, and `pending_credits` is money that has arrived but is not usable yet. Reading the wrong one authorizes a payment that gets declined. A declined transaction still shows up in the listing with an amount and moved no money, so a total that does not filter on status overcounts by everything that failed. The card number itself is never returned, only a masked form and the last four digits.

Just-in-time funding, arguably the most distinctive thing Marqeta does, is not modelled -- it puts a webhook in the authorization path with a response deadline, and Cauldron delivers webhooks without waiting for an answer, so the one behaviour that would most reward reproducing is structurally out of reach.

## Sources

- Documentation: https://www.marqeta.com/docs/core-api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve marqeta     # run it
cauldron verify marqeta -v # check every claim
```
