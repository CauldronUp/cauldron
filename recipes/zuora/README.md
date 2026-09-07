# zuora

Emulates the Zuora billing API for local development and tests.

**12 conformance cases, 5 checked against the live API on 2026-09-07.**

Written against Zuora's reference at `developer.zuora.com` and struck live against `rest.zuora.com` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**One code, two sentences, and the key order changes between them.**

```
GET /v1/accounts    (no credential)
401 {"reasons":[{"code":"90000011","message":"Unauthorized"}],"success":false}

GET /v1/accounts    Authorization: Bearer notreal
401 {"reasons":[{"message":"Authentication error","code":"90000011"}],"success":false}
```

The code is `90000011` for both. The sentence is not. **The field that exists to be branched on cannot separate "you sent nothing" from "you sent something wrong", and the prose can** — which is the wrong way round.

And the keys arrive in a different order: `code` first on one, `message` first on the other. Two code paths building the same object. The only reader that notices is one comparing bytes — which is every recorded fixture, every snapshot test and every cache key built from a response body.

**The code is eight digits in a string.** `"90000011"` — wide enough to encode a category and an ordinal, and quoted, so `code === 90000011` is false and `parseInt` comes first. There is no separate category field, so whatever structure lives inside those digits is undocumented on the wire.

**The success flag rides the failure.** `"success": false` beside the reasons, so a client checks one boolean rather than a status. The good half of this envelope, and the third field name for the same idea in this collection after [loops](../loops)'s `success` and [sift](../sift)'s `status: 0`.

**The field is called `reasons`.** Not `errors`, not `messages` — plural, because Zuora's validation genuinely can fail several ways at once and each gets an entry. A better name than `errors` for what it holds.

**An unknown path answers the credential failure**, so routing is judged after authentication.

**An account owing money is still `Active`.** `status` is the subscription lifecycle rather than the dunning state, so a client reading it to decide whether to chase somebody reads the wrong field; the balance is the one that knows.

**Money is a bare JSON number.** `1249.50` as a double, on a billing API, with the currency in a separate field — so every sum a client computes is binary floating-point.

## Modelling limits

- **One route.** Accounts. Subscriptions, invoices, payments, amendments, the whole revenue-recognition surface and the Object Query Language endpoint each want their own evidence.
- **Nothing is mapped in detection.** Zuora's SDKs are generated per tenant and per version and none resolves under an obvious registry name; nothing found on 2026-09-07 calls `rest.zuora.com`.
- **No `spec:`.** Zuora publishes a rendered reference site.
