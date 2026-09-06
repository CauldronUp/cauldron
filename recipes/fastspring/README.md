# fastspring

Emulates the fastspring API (v1), for local development and tests.

**21 conformance cases, 3 checked against the live API on 2026-09-05.**

The resource cases cite documentation rather than an observation on a real store; the refusal cases were struck live, unauthenticated, against api.fastspring.com.

## What this Recipe found

One purchase produces three different objects with three different identifiers, and which one a webhook hands you depends on which event fired. An `order.completed` webhook carries the order id; a `subscription.charge.completed` carries the subscription id plus a different order id for that period's charge; a `subscription.deactivated` carries only the subscription. Code that keys everything off "the order id" ends up with three different values over one customer's lifetime, with no single lookup that reconciles them.

The live probe found FastSpring sitting behind an AWS API Gateway that refuses before anything downstream runs: a missing credential, an invented one, an unrouted path, and a wrong method on a real path all answer the identical 401 with zero bytes, not the JSON `{"error":"Unauthorized"}` this file declared before the probe. There is no distinguishing a bad path or a bad method from a bad credential here -- the gateway answers the same way to all three.

A cancelled subscription also keeps billing until its period actually ends: `state` reads `"active"` with a `deactivationDate` set in the future, so checking state alone bills someone who already left, and checking the date alone cuts off someone who's paid through next month. Every amount besides arrives twice -- a number and a localized, preformatted string with a currency symbol baked in -- and test and live orders share the same API, distinguished only by a boolean a client has to remember to filter on.

## Sources

- Documentation: https://developer.fastspring.com/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve fastspring     # run it
cauldron verify fastspring -v # check every claim
```
