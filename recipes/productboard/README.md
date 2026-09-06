# productboard

Emulates the Productboard API for local development and tests.

**9 conformance cases, 6 checked against the live API on 2026-09-06.**

Written against Productboard's developer documentation and struck live against `api.productboard.com` on 2026-09-06. Every live request used a deliberately invalid credential.

## What this Recipe found

**The API is two months past its own sunset date and still answering.**

Every response carries, struck live:

```
Deprecation: @1775692799
Sunset: Wed, 08 Jul 2026 23:59:59 GMT
```

`@1775692799` is `2026-04-08T23:59:59Z`. The sunset is exactly three months later, 2026-07-08. This was written on 2026-09-06. So the surface announced its deprecation in April, said it would stop in July, and served every request behind this Recipe in September with both headers still attached.

**And the two dates are in two different formats on the same response.** `Sunset` is an HTTP-date, which RFC 8594 requires. `Deprecation` is `@1775692799` — a Unix timestamp with an at-sign, the syntax an early draft used and which RFC 9745 replaced with an HTTP-date. One lifecycle, two headers, two parsers, and the one a client is most likely to mis-parse is the one that fires first.

A client that checks `Sunset` and refuses to call a sunset API stopped working in July, against an API that works. One that ignores both learns nothing. Neither behaviour is wrong, exactly — the headers are correct HTTP and simply are not true.

**The 401 reads your token and tells you what is wrong with its payload.**

Three refusals, all 401, all obtained with no account behind them:

| Sent | `message` |
|---|---|
| nothing | `Unauthorized` |
| `Bearer not-a-real-key` | `Bad token; invalid JSON` |
| `Bearer eyJ…` (well-formed) | `No mandatory 'iss' in claims` |

The second says it base64-decoded the credential and tried to parse the result. The third says it parsed it, looked inside, and found no issuer claim.

So the failure path is a **JWT-shaping oracle**: anyone can iterate on the structure of a token and be told, claim by claim, what the server wants next — and the signature check, the part that cannot be guessed at, happens after all of it. That is genuinely helpful to somebody debugging their own integration, and it is exactly the same help to somebody who is not.

Each of those carries `WWW-Authenticate: Bearer error="invalid_token"`, which is RFC 6750 done properly, on the same responses.

**The 404 is the gateway, not the product.**

```
GET /cauldron-nonexistent
404 {"message":"no Route matched with those values"}
```

That sentence is Kong's. So `message` — the only field present in every failure here — has **two producers**: Productboard writes it for a credential and Kong writes it for a path, in the same key, with nothing in the body saying which answered. A client logging `message` files two systems under one name.

**A status is an object, not a string.** `data[0].status` is `{id, name}`, so a client comparing `feature.status` to `"Candidate"` compares an object to a word.

## Modelling limits

- **One route.** Features. Notes, components, products, releases, custom fields and the webhook surface each want their own evidence.
- **`X-Version` is declared and not exercised as a failure.** Struck live, the credential is checked first, so what happens when the version header is missing is not observable from outside.
- **Nothing is mapped in detection.** No client package for this API turned up on npm, Packagist or the Go module proxy under any obvious name.
- **No `spec:`.** Productboard's reference is a rendered documentation site.
