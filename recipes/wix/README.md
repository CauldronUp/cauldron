# Wix

Emulates the Wix API (v1), for local development and tests.

**15 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

The shipping address is not where an order ships to, on a pickup order. Wix's own field description says `shippingDestination` holds the pickup location's address for pickup orders, not the customer's, and the actual recipient is in a differently-named field, `recipientInfo`. Every integration that prints a label reaches for the field with "shipping" in the name, and on a click-and-collect order that's the shop's own depot address rather than the customer's -- the parcel gets posted to the pickup point it was meant to be collected from.

Importing historical orders is also invisible to normal event handlers: `POST /orders/import` fires a dedicated `order_imported` event and never `order_created`, so a shop migrating years of history from another platform produces zero of the event most integrations actually listen for.

## Sources

- Documentation: https://dev.wix.com/docs/rest/business-solutions/e-commerce/orders/orders
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve wix     # run it
cauldron verify wix -v # check every claim
```
