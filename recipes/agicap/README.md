# agicap

Emulates the Agicap Open API for local development and tests.

**11 conformance cases, 5 checked against the live API on 2026-09-07.**

Written against Agicap's reference at `openapi.agicap.com` and struck live against `openapi.agicap.com` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**An unauthenticated request is a 500, and the reason is the rate limiter.**

```
GET /api/customers    (no credential)
500  Content-Type: application/problem+json
{"type":"https://tools.ietf.org/html/rfc7231#section-6.6.1",
 "title":"An error occurred while processing your request.",
 "status":500,
 "detail":"The current request and the key build strategy can't be able to build a valid key
           to control rate limit"}
```

The rate limiter partitions by client, derives its bucket key from the caller's identity, and an anonymous request has none — so building the key throws, the exception escapes, and the framework's unhandled-error page comes back as a problem document.

**So the first request anybody makes against this API, before they have wired up a credential, is a server error.** Not a refusal, not a 401 — a 500, which every retry policy in the world treats as "try again", and which will keep being a 500 however many times it is tried.

**The `type` URI points at the definition of 500.** `rfc7231#section-6.6.1` is the paragraph in the HTTP specification that defines the status code. RFC 9457 reserves `type` for identifying the *problem*; this identifies the status, which the document already carries in two other fields. It is the third `problem+json` document here and the third distinct way of getting that field wrong — see [revai](../revai) for the other two.

**`detail` is the internal failure, in broken English.** A sentence about a `KeyBuilder` strategy inside the middleware, printed to an anonymous caller. It names an internal component and describes a code path — and it is also, unusually, *true and specific*: a reader who understands it knows exactly why they got a 500.

**A wrong credential is a 404, not a 401.** With an Authorization header the same path answers 404 with no body. So sending nothing is a server error, sending something wrong is a missing endpoint, and **nothing this API does to a bad credential involves the status meant for one**.

**Currency is per customer rather than per company**, on a cash-flow product — so summing `outstandingAmount` across a listing adds euros to pounds unless the client groups first. And the amounts are bare JSON doubles.

## Modelling limits

- **One route.** Customers. Bank accounts, transactions, forecasts, budgets and the whole cash-flow surface each want their own evidence.
- **Nothing is mapped in detection.** Agicap's Open API is behind an account and publishes no client on any registry.
- **No `spec:`.** The reference is a rendered site behind a login.
