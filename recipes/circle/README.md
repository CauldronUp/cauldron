# circle

Emulates the Circle payments API for local development and tests.

**9 conformance cases, 6 checked against the live API on 2026-09-06.**

Written against Circle's developer documentation at `developers.circle.com` and struck live against `api.circle.com` on 2026-09-06. Every live request used a deliberately invalid credential; the balance read needs an account and is documentation-derived.

## What this Recipe found

**A not-found is `code: -1`.**

```
GET /v1/cauldron-nonexistent
404 {"code":-1,"message":"Resource not found"}
```

Every other failure carries a `code` that *is* the HTTP status — `401` beside a 401. The router's own failure carries **-1**: a sentinel where a code belongs, negative so it cannot collide with a status, and explained nowhere. A client switching on `code` handles 401 and then meets a number no table covers.

**Three different 401s, and the shape of the key decides which.** All three struck live:

| Sent | `message` |
|---|---|
| nothing | `malformed authorization. Missing API key in authorization header. Make sure to use Bearer authorization type` |
| `Bearer not-a-real-key` | `malformed API key. API key should contain three substrings, separated by a colon. First substring denotes the environment, e.g. TEST_API_KEY. Second substring denotes the ID of the key, e.g. ebb3ad72232624921abc4b162148bb84. Third substring denotes the secret of the key, e.g. 019ef3358ef9cd6d08fc32csfe89a68d. This format only applies to API keys generated after May 2023` |
| `Bearer TEST_API_KEY:abc:def` | `Invalid credentials.` |

The middle one is **a four-sentence tutorial served as an error message**, and it prints example key material — a sample key id and a sample secret — into the body of a 401 that anyone can trigger without an account. They are fictional, and they are also exactly what a real one looks like, which is the reason for including them and the reason they sit in a response rather than only in the documentation, where the same paragraph already lives.

It also carries a date caveat: *"This format only applies to API keys generated after May 2023."* So the error describing the credential format is conditional on when the credential was issued — and does not say what the older format looks like, which is the case a reader hitting this message is most likely to be in.

The punctuation is worth noticing too: two of the three begin lower case and end without a full stop, and the third does neither. A product showing `message` to a user gets three registers from one status.

**The health check is not in the API.** `GET /ping` answers `{"message":"pong"}` — 200, no `code`, the only response here without one, and the only one needing no credential. `GET /v1/ping`, the same name one segment in, answers `{"code":-1,"message":"Resource not found"}`.

So the endpoint a monitor calls sits outside the versioned surface and outside its envelope, and putting it where the rest of the API lives produces a 404 that reads like an outage.

**Money is a string with its places kept.** `"1250.00"`, and a zero balance is `"0.00"` — not `"0"`, not `0`. That is the right decision, and there is no `scale` or `decimals` field to say how many places to expect, so a client formatting the value has to infer it from the currency.

## Detection

`@circle-fin/developer-controlled-wallets` names `api.circle.com` 13 times in its published archive and is mapped. It is Circle's wallets SDK rather than a client of the Circle Mint surface this Recipe serves — the same host, the same credential, a different path family — so the mapping says "this project talks to Circle", which is what detection is for, and not "this project calls these two routes".

## Modelling limits

- **Two routes.** Balances and the health check. Transfers, payouts, wallets, payment intents and the whole settlement lifecycle each want their own evidence, and none of it is reachable without an account.
- **`data.unsettled` is not served.** The balances response carries `available` and `unsettled` as two arrays under one object; this Recipe serves the available half and the shape of the envelope, which is where the nesting shows.
- **No `spec:`.** Circle's reference is a documentation site; no machine-readable description is published at a stable URL that `cauldron drift` could fingerprint.
