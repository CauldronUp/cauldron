# Wix

Emulates the Wix API (v1), for local development and tests.

**18 conformance cases, 4 checked against the live API on 2026-09-05.**

The order and fulfilment cases still cite documentation, since a real order needs a real site. The credential shape needed no site at all, and checking it live found this Recipe's single failure model wrong for every operation it checked.

## What this Recipe found

The shipping address is not where an order ships to, on a pickup order. Wix's own field description says `shippingDestination` holds the pickup location's address for pickup orders, not the customer's, and the actual recipient is in a differently-named field, `recipientInfo`. Every integration that prints a label reaches for the field with "shipping" in the name, and on a click-and-collect order that's the shop's own depot address rather than the customer's -- the parcel gets posted to the pickup point it was meant to be collected from.

Importing historical orders is also invisible to normal event handlers: `POST /orders/import` fires a dedicated `order_imported` event and never `order_created`, so a shop migrating years of history from another platform produces zero of the event most integrations actually listen for.

## What checking it live found

Wix does not have one credential-failure shape -- it has one per operation. SearchOrders answers 403, shaped like a permission refusal and naming the operation in both the sentence and the code (`"read order: permission denied"`, `READ_ORDER_FORBIDDEN`); GetOrder answers 428 Precondition Required with a different sentence and code again (`MISSING_IDENTITY_CONTEXT`). Neither matched this Recipe's single `401 UNAUTHENTICATED`, including an existing case that had asserted it for the exact GetOrder request checked live. A write (CreateOrder), sent with an empty body and no credential, never surfaced a credential complaint at all -- it answered a plain body-validation error instead. A path nothing declares is a 404 with nothing in it, resolved before any of the above.

## Sources

- Documentation: https://dev.wix.com/docs/rest/business-solutions/e-commerce/orders/orders
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve wix     # run it
cauldron verify wix -v # check every claim
```
