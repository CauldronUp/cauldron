# Easyship

Emulates the Easyship API (v2), for local development and tests.

**15 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A cross-border shipping rate that sorts first because it's cheapest is often cheaper only because the customer, not the seller, is on the hook for duty and tax -- DDU rates omit import_duty_charge and import_tax_charge entirely and so are lower by exactly that amount, while DDP rates fold them in. Sorting by total_charge and taking the first one routes customers into paying a bill on the doorstep, and some refuse delivery, sending the parcel back at the seller's expense with nothing having errored anywhere.

Delivery is also always a range (min_delivery_time, max_delivery_time), never a date, so a client that promises the minimum is promising the best case. And the currency on a rate is the store's, not the customer's, so the figure quoted isn't necessarily the figure that ends up on a card statement.

No rate is actually calculated here -- which couriers can carry a parcel and what duty a country charges is the entire product, and the rates in this Recipe are fixtures rather than a real calculation for the given HS code and destination.

## Sources

- Documentation: https://developers.easyship.com/reference/rates
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve easyship     # run it
cauldron verify easyship -v # check every claim
```
