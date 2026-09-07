# mercadopago

Emulates the Mercado Pago payments API for local development and tests.

**11 conformance cases, 5 checked against the live API on 2026-09-07.**

Written against Mercado Pago's reference at `developers.mercadopago.com` and struck live against `api.mercadopago.com` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**A timestamp and a UUID, joined by a semicolon, in one string.**

```json
{"message":"Unauthorized use of live credentials","error":"unauthorized","status":401,
 "cause":[{"code":7,
           "description":"Unauthorized use of live credentials",
           "data":"07-09-2026T21:50:21UTC;b05e2e13-0c9b-478d-8bab-fb2b47a2d3c3"}]}
```

`data` is two values in one field.

The first is a timestamp in `DD-MM-YYYYThh:mm:ssUTC` — day first, `T` in the middle, and the literal letters `UTC` where an offset belongs. It is neither ISO 8601 nor RFC 3339, and **no date library parses it**. The second is a request id.

So the two things a caller most wants from a failure — when it happened and what to quote to support — are both present, concatenated, in the one field whose name promises neither.

**"Unauthorized use of live credentials" is the answer to sending none.** Struck live with no `Authorization` header at all. The sentence names a specific, real and common mistake in this API — a production key where a test one belongs — and it is the answer to a request that carried no key of any kind. The sixth provider in this collection to describe the wrong failure, and the most specific about the wrong thing.

**The message appears twice and the status three times** — `message` and `cause[0].description` are the same string; `status`, the HTTP status and `error` all carry the same fact.

**The 404 is in Spanish.**

```json
{"error":"resource not found",
 "message":"Si quieres conocer los recursos de la API que se encuentran disponibles
            visita el Sitio de Desarrolladores de MercadoLibre (https://developers.mercadopago.com)"}
```

The `error` field is English and the `message` beside it is Spanish, in one body — and it names **MercadoLibre**, the parent marketplace, while linking `mercadopago.com`. The sentence names one company and the URL another.

**The reason a payment is held lives in a second field.** `status` is one word for a payment held for fraud review, one held for funds and one held for a document check; `status_detail` is the difference, so a client switching on `status` cannot tell a customer what to do.

## Modelling limits

- **One route.** Payment search. Preferences, refunds, chargebacks, merchant orders, subscriptions and the whole Checkout Pro surface each want their own evidence.
- **Nothing is mapped in detection.** Mercado Pago's SDKs are named for the company and span several countries' products.
- **No `spec:`.** Mercado Pago publishes a rendered reference site.
