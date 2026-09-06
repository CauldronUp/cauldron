# Mollie

Emulates the Mollie API (v2), for local development and tests.

**18 conformance cases, 1 checked against the live API.**

Struck live 2026-09-05 against api.mollie.com, no account and no key -- a missing credential and a wrong one both answer the identical body. This file declared the auth failure as a 401 with an invented sentence; the real answer is a 400, "Invalid Authorization header". Fixed below.

## What this Recipe found

Mollie's webhook body is an id and nothing else -- deliberately, because the notification URL is public and anybody could post to it, so you fetch the payment yourself and read the status rather than trusting whatever the POST body claims. An integration built against an emulator that helpfully included the whole payment object learns a status the real webhook never carries, and finds out on its first live transaction.

`open` is not `pending`, and neither is `paid`: a payment starts `open`, meaning the customer has not finished at checkout, while `pending` is a separate state meaning the payment method is just slow and can still fail. Treating anything other than `failed` as success books orders nobody paid for. The checkout link (`_links.checkout`) is present only while the payment is open and disappears the moment it is not -- its absence is the only signal, and an application that caches the link and re-displays it sends a customer to a dead page. Money is a decimal string with the currency beside it (`"10.00"`, not `1000`), so treating Mollie like Stripe is out by a factor of a hundred in the direction that looks like a bargain.

Paying a payment is not modelled, only creating one -- the actual state transition happens at a hosted checkout page this format cannot reproduce, so which state a payment holds here is whatever a fixture puts there.

## Sources

- Documentation: https://docs.mollie.com/reference/v2/payments-api/get-payment
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve mollie     # run it
cauldron verify mollie -v # check every claim
```
