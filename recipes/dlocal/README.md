# dlocal

Emulates the dLocal payments API for local development and tests.

**11 conformance cases, 4 checked against the live API on 2026-09-07.**

Written against dLocal's reference at `docs.dlocal.com` and struck live against `api.dlocal.com` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**The refusal is a ladder, and the first rung is a clock.**

Struck live, three requests:

```
GET /payments    (no headers at all)
400 {"code":5001,"message":"Invalid parameter","param":"X-Date"}

GET /payments    X-Date: 2026-09-07T21:30:00.000Z
400 {"code":5001,"message":"Missing parameter(s) [or non valid values]. login, key"}

GET /payments    X-Date, X-Login, X-Trans-Key, Authorization
403 {"code":3001,"message":"Invalid credentials"}
```

**The first thing dLocal checks is not the credential — it is the timestamp.** `X-Date` is part of the HMAC signature scheme, so a request without one cannot be verified at all, and the API says so before it says anything about who you are.

That is defensible, and it means the first failure a new integrator meets names a header they have never heard of, on an endpoint they were trying to authenticate against.

**The same code covers two different missing things.** `5001` for a missing `X-Date` and `5001` again for missing `login` and `key` — and the first carries a `param` field naming the header while the second does not, listing its two names inside the sentence instead. So the machine-readable half is one number for two problems, and the structured field that would say which is present on only one of them.

**The final refusal is 403.** For a request carrying every required header with wrong values in them. Not 401 — so across the whole ladder this API answers **400, 400 and 403** to what is, from the caller's side, one problem: not being authenticated yet.

**The codes are four digits whose leading digit is a category.** 5001 is a parameter problem and 3001 a credential one, and nothing on the wire says so.

**`PENDING` is the normal state rather than a transient one.** dLocal's local payment methods — boleto, PIX, cash at a shop — settle over days. A client polling for `PAID` with a short timeout gives up on a payment that will arrive.

**`status_code` is a string that reads like an HTTP status.** `"200"` on a payment record, which is not one — so a log line carrying it is ambiguous about what succeeded.

**Timestamps carry an offset with no colon.** `+0000` rather than `+00:00` or `Z` — legal ISO 8601 basic format and not RFC 3339, so a strict parser refuses it and a lenient one does not. The request's own `X-Date` uses `Z`, so one exchange has two spellings of UTC.

## Modelling limits

- **One route.** Payments. Refunds, payouts, chargebacks, the whole payment-method catalogue and the country-specific flows each want their own evidence.
- **The HMAC signature is not modelled.** This Recipe checks `X-Login`, which is the header the ladder's third rung refuses; reproducing the signature would mean modelling a shared secret the emulator has no business holding.
- **Nothing is mapped in detection.** No client for this API turned up on npm, Packagist or the Go module proxy under an obvious name on 2026-09-07.
- **No `spec:`.** dLocal publishes a rendered reference site.
