# truelayer

Emulates the truelayer API (v1), for local development and tests.

**26 conformance cases, 2 checked against the live API on 2026-09-05.**

The consent and transaction cases still cite documentation, since a real connection needs a real bank. The credential shape needed no consent at all, and checking it live found this Recipe's own message wrong.

## What this Recipe found

A consent lasts ninety days from when it was granted, not from last use, and on day ninety-one every call against that connection fails -- the only fix is the customer going back to their bank to re-authorise. The expiry date is stated exactly once, in the original token response; miss it and you find out by failing. A revoked consent and an expired one also both arrive as the same 403, distinguished only by an error code buried inside it, and only one of those two is something prompting the customer sooner would have prevented.

Pending and settled transactions are two different lists for the same money, and a pending transaction's `transaction_id` changes once it settles -- matching the two lists on id loses every transaction exactly once.

## What checking it live found

An absent credential and a garbage bearer value carry the same code, `invalid_token`, and disagree only on `error_description` -- empty for nothing sent, `"The signature is invalid"` for something sent that does not verify. Neither was the generic sentence this Recipe had claimed for both. A path nothing declares gets the absent-shaped answer too, checked before routing.

## Sources

- Documentation: https://docs.truelayer.com/docs/data-api-basics
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve truelayer     # run it
cauldron verify truelayer -v # check every claim
```
