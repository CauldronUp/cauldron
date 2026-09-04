# fastspring

Emulates the fastspring API (v1), for local development and tests.

**13 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

One purchase produces three different objects with three different identifiers, and which one a webhook hands you depends on which event fired. An `order.completed` webhook carries the order id; a `subscription.charge.completed` carries the subscription id plus a different order id for that period's charge; a `subscription.deactivated` carries only the subscription. Code that keys everything off "the order id" ends up with three different values over one customer's lifetime, with no single lookup that reconciles them.

A cancelled subscription also keeps billing until its period actually ends: `state` reads `"active"` with a `deactivationDate` set in the future, so checking state alone bills someone who already left, and checking the date alone cuts off someone who's paid through next month. Every amount besides arrives twice -- a number and a localized, preformatted string with a currency symbol baked in -- and test and live orders share the same API, distinguished only by a boolean a client has to remember to filter on.

## Sources

- Documentation: https://developer.fastspring.com/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve fastspring     # run it
cauldron verify fastspring -v # check every claim
```
