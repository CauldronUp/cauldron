# Ecwid

Emulates the Ecwid API (v3), for local development and tests.

**15 conformance cases, 3 checked against the live API on 2026-09-05.**

The resource cases cite documentation rather than an observation on a real store; the refusal cases were struck live, unauthenticated, against app.ecwid.com.

## What this Recipe found

The default order listing quietly excludes abandoned carts. Ecwid's own docs say that with no filters set, the API returns all orders except unfinished ones -- to see them you have to ask for `paymentStatus=INCOMPLETE` by name. The response gives no sign anything was left out: `count` matches the array, paging matches `count`, and every number stays internally consistent. A shop reconciling takings never notices; a shop chasing cart recovery never learns to ask.

The live probe found this file's own claim about authentication wrong: the document's worked example shows a JSON `{"errorMessage":"Access token is invalid."}` on a bad credential, and the real API answers zero bytes -- reproducible with no token, a wrong one, or a right-shaped nonsense one alike. An unrouted path and a wrong method on a real path get the same emptiness at their own status codes, needing no credential either.

A handful of fields carry the same kind of trap. `acceptMarketing` is a boolean where `null` counts as a yes and only `false` means no, so the obvious `if (order.acceptMarketing)` check silently treats a genuine yes as a refusal. `discount` is documented as everything except the coupon discount, so the field with the obvious name is not the total order discount -- `couponDiscount` has to be added back in to get it. And `paymentMessage`, the reason a payment is pending, is cleared the moment the order becomes paid, so it's only ever visible while it doesn't matter yet.

## Sources

- Documentation: https://api-docs.ecwid.com/reference/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve ecwid     # run it
cauldron verify ecwid -v # check every claim
```
