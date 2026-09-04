# truelayer

Emulates the truelayer API (v1), for local development and tests.

**14 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A consent lasts ninety days from when it was granted, not from last use, and on day ninety-one every call against that connection fails -- the only fix is the customer going back to their bank to re-authorise. The expiry date is stated exactly once, in the original token response; miss it and you find out by failing. A revoked consent and an expired one also both arrive as the same 403, distinguished only by an error code buried inside it, and only one of those two is something prompting the customer sooner would have prevented.

Pending and settled transactions are two different lists for the same money, and a pending transaction's `transaction_id` changes once it settles -- matching the two lists on id loses every transaction exactly once.

## Sources

- Documentation: https://docs.truelayer.com/docs/data-api-basics
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve truelayer     # run it
cauldron verify truelayer -v # check every claim
```
