# Allegro

Emulates the Allegro API (v1), for local development and tests.

**11 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

An Allegro order exists before it's paid, and the endpoint's own name says so: it isn't /orders, it's /order/checkout-forms, and Allegro's own state descriptions spell out why -- FILLED_IN means "checkout form filled in but payment is not completed yet so data could still change", and only READY_FOR_PROCESSING means the money actually arrived. A row can appear in the listing before the buyer has entered an address, and the address that appears next isn't final; code that treats "a new order appeared" as "ship it" ships against a delivery address the buyer can still edit, for an order that may never be paid.

There are two independent status fields that both use CANCELLED to mean different things: status tracks the buyer's checkout, fulfillment.status tracks the seller's own fulfillment -- and per Allegro's own docs, fulfillment.status can move on its own (RETURNED can't be set by the seller at all). Every failure also carries two separate messages: one in English for developers, one translated and safe to show a shopper -- except on a 401, which arrives as a bare OAuth2 {error, error_description} body instead, so the unwrapping every other failure needs finds nothing there.

Money throughout is a decimal string ("the amount provided in a string format to avoid rounding errors", in Allegro's own words) and the API version is chosen via the Accept header rather than the URL, with parts of the API only reachable through a beta media type -- asking for plain application/json gets a documented 406.

## Sources

- Documentation: https://developer.allegro.pl/documentation
- Machine-readable description: https://developer.allegro.pl/swagger.yaml, last checked 2026-08-31
  `cauldron drift allegro` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve allegro     # run it
cauldron verify allegro -v # check every claim
```
