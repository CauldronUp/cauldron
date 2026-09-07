# eightbyeight

Emulates the 8x8 analytics API for local development and tests.

**9 conformance cases, 4 checked against the live API on 2026-09-07.**

Written against 8x8's reference at `developer.8x8.com` and struck live against `api.8x8.com` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**The request id is the same on two different requests.**

Struck live, seconds apart:

```
GET /vms/v1/subaccounts    (no credential)
404 {"message":"no Route matched with those values","request_id":"658c4f55152fc611583ee84dfaa8800e"}

GET /vms/v1/subaccounts    Authorization: Bearer notreal
404 {"message":"no Route matched with those values","request_id":"658c4f55152fc611583ee84dfaa8800e"}

GET /
404 {"message":"no Route matched with those values","request_id":"8011fa5566f3b055456e3d4e8ab27d45"}
```

The first two carry an **identical** `request_id`. The third, a different path, carries a different one. So the value is derived from the request rather than minted per call — **two calls to the same URL are one id**, and the field that exists to identify *this* request identifies the route.

That defeats the thing a request id is for. A caller reporting "this failed at 14:02" hands over a value that also belongs to every other caller who made the same request, so support cannot find the call.

It is also the only correlation id in this collection that is *reproducible*: two clients on different continents get the same 32 hex characters for the same URL.

**Everything is Kong.** `no Route matched with those values` is Kong's own sentence, and it answers a credentialed request, an uncredentialed one and the host root alike — so 8x8's own application is never reached from outside and the gateway is the entire observable API.

This is the **third** Kong sighting here. [productboard](../productboard) answers it on unknown paths while serving its own 401 on real ones; [browserbase](../browserbase) leaks Fastify the same way; here the gateway answers everything. Three products, one piece of infrastructure, and in each case the sentence a developer sees was written by somebody who has never heard of the company they are integrating with.

**A credential failure is reported as a missing route.** The route is only mounted for authenticated consumers, so an anonymous request is not refused — it is unrouted. The status says the endpoint does not exist. It does; what does not exist is the caller's permission to see it.

**The listing is Spring Data's `Page`, on the wire.** `content`, `totalElements`, `number`, `size` — a Java framework's serialisation rather than a designed envelope, and `number` is not what any client would guess for "which page".

**A suspended subaccount has no date saying when.** A status and a creation time, and nothing recording the change — so one suspended this morning and one suspended last year are the same record.

## Modelling limits

- **One route.** Subaccounts. Calls, SMS, chat, quality data and the whole Contact Centre surface each want their own evidence.
- **Nothing is mapped in detection.** No client for this API turned up on npm, Packagist or the Go module proxy under an obvious name on 2026-09-07.
- **No `spec:`.** 8x8 publishes a rendered reference site.
