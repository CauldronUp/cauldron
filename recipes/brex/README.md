# Brex

Emulates the Brex API (v2), for local development and tests.

**11 conformance cases, 1 checked against the live API.**

Struck live against platform.brexapis.com on 2026-09-05: an unauthorized request carries no body at all, byte for byte the same whether the Authorization header is absent or a made-up Bearer value. This file had claimed the message "The access token is invalid.", which Brex does not send.

## What this Recipe found

Brex reports money as an object with amount and currency, where the amount is still cents -- the same value Ramp reports as a bare integer with no unit -- so code handling one finds a number where it expected an object or the reverse, and doing arithmetic on Brex's object directly produces NaN rather than a silently wrong figure, which is the better of the two failures and still a failure. A card also exists before it can be used: status moves through PENDING to ACTIVE, and a card created a moment ago declines everything while looking fully issued, with no field anywhere saying why.

Expenses and transactions are separate objects describing the same money -- a transaction is what the card did, an expense is the record somebody has to attach a receipt to -- so an expense can be missing a receipt while its transaction is perfectly settled. And pagination's cursor is only ever present when there's a next page; an empty string is not the same as absent, and Brex omits the field entirely rather than sending one.

Card creation is deliberately not modelled here, for the same reason Mercury's payment creation isn't: an emulator that makes issuing corporate cards easy to exercise would be inviting the exact mistake it warns about.

## Sources

- Documentation: https://developer.brex.com/openapi/expenses_api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve brex     # run it
cauldron verify brex -v # check every claim
```
