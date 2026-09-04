# revenuecat

Emulates the revenuecat API (v1), for local development and tests.

**12 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

RevenueCat's own SDK guidance tells you to read `customerInfo.entitlements[id].isActive`, a property the client SDKs compute for you that does not exist in the API response itself. The moment the question moves server-side, a webhook handler, your own entitlement check, somebody has to write that comparison by hand, and the advice they were following no longer applies. The real comparison is `expires_date` against now, with four special cases the SDK property was quietly absorbing: `null` means lifetime access, not expired; an `unsubscribe_detected_at` timestamp means auto-renewal is off but the customer keeps access until `expires_date`, possibly months away; a trial with a future date is a paying-nothing customer who is nonetheless entitled to everything; and `ownership_type: FAMILY_SHARED` is an entitlement from a purchase this customer did not make and cannot manage. The fixture models four active entitlements, no two active for the same reason.

The endpoint that fetches a customer is a GET that can create one -- it is literally titled Get or Create Customer, returning 200 if the customer existed and 201 if this very request invented them, with an identical, empty-entitlements body either way. A status code on a GET is not something most client code branches on, so this is the only observable difference between an existing customer with no subscription and a customer who did not exist a moment ago.

Two maps on the customer object are keyed differently -- `entitlements` by the identifiers named in the dashboard, `subscriptions` by the product identifiers the app stores assigned -- and the only bridge between them is a `product_identifier` field nested inside each entitlement, so reaching into `subscriptions` with an entitlement's own key finds nothing.

## Sources

- Documentation: https://www.revenuecat.com/docs/api-v1
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve revenuecat     # run it
cauldron verify revenuecat -v # check every claim
```
