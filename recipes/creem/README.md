# creem

Emulates the Creem product listing for local development and tests.

**12 conformance cases, 4 checked against the live API on 2026-09-13.**

Read from the OpenAPI document Creem publishes at [`docs.creem.io/api-reference/openapi.json`](https://docs.creem.io/api-reference/openapi.json), and struck live against `api.creem.io` on 2026-09-13 with no key, with a wrong key, on a path that does not exist, and with a method a real path does not take.

## What this Recipe found

**A 404 whose body says 500.**

```
GET /v1/cauldron-nope   404  {"statusCode":500,"timestamp":"2026-09-13T21:58:34.410Z","path":"/v1/cauldron-nope"}
PUT /v1/products        404  {"statusCode":500,"timestamp":"2026-09-13T21:58:34.607Z","path":"/v1/products"}
```

The status line says the thing was not found and the body says the server broke. A client that trusts the body reports an outage; a client that trusts the header reports a typo; both are reading the same response.

**And that body is a different shape from the other failures.** A missing key answers:

```json
{"trace_id":"d3f88409-...","status":401,"error":"Unauthorized","message":["API Key is missing"],"timestamp":1789336714203}
```

So on one API the status lives under `status` or under `statusCode`; the timestamp is epoch milliseconds as a number or an ISO string; and a failure carries a trace id and a sentence, or a path and neither. Two error contracts, and which one arrives depends on how far the request got.

**The sentence is an array and the category beside it is not.** `message` is `["API Key is missing"]` — one string, in a list — while `error` is the bare string `"Unauthorized"`. Two prose fields on one object, one of them plural for no reason a caller can see.

**A refused request carries an entity tag.** `etag: W/"92-JJAKp4OszPkasCCGzKi3bcBVyVI"` on the 401, so the refusal is cacheable and revalidatable, with `access-control-allow-credentials: true` and `vary: Origin` beside it.

**The environment enum has three values and two of them mean the same thing.** `EnvironmentMode` is `["test", "prod", "sandbox"]`. And the reference's own webhook samples ship `"mode": "local"` on live-shaped records — a fourth value the schema does not allow, printed in the published documentation.

**A product says how often it bills, twice, in two languages.** `billing_period` is an English phrase from `["every-month", "every-three-months", "every-six-months", "every-year", "every-day", "once", "custom"]`, and `recurring_interval: "month"` with `recurring_interval_count: 3` sits on the same record. The pair can express every fourth month; the phrase cannot. `once` is not a period, and `custom` says nothing at all.

**The image is singular and plural at once.** `image_url` is a string and `image_urls` is an array, both on the record, and the array's example is the string's example wrapped in brackets.

**Two absences on one record.** `default_success_url` is the empty string on one product and `null` on another, for the same "not set".

**The collection answers at two addresses and only one is documented.** `/v1/products/search` is the operation the reference describes; `/v1/products` answers the same listing live and appears nowhere in the document, which is why `cauldron drift` reports it undeclared.

**And the price is a bare number.** `price: 1100` with `currency: "EUR"` — minor units, with nothing on the record saying so, next to a `tax_mode` that decides whether that number already contains the tax.

## Sources

- [`docs.creem.io/api-reference/openapi.json`](https://docs.creem.io/api-reference/openapi.json) — recorded by `cauldron drift`, `spec_seen: 2026-09-13`.
- `docs.creem.io/llms-full.txt` — the webhook samples that print `"mode": "local"`.
- Live: `api.creem.io`, struck 2026-09-13 with no key, a wrong key, an unrouted path, and a wrong method.

## What this needed from Cauldron

`cursor_number` on the list envelope, new for this Recipe. Creem's paging fields are numbers — `next_page: 2`, `prev_page: 1`, null at each end — where every previous Recipe's pointer was an opaque token or an address. A cursor rendered as the string `"2"` reads identically in a diff and is a different value on the wire, and a client adding one to it gets `"21"`. Declared, the next and previous fields are JSON numbers; omitted, nothing about any existing Recipe moves.

## Modelling limits

- **Two routes, one collection.** Products, at both of its addresses. Customers, subscriptions, checkouts, discounts, licences, transactions and the webhook surface are the rest.
- **`trace_id` and `timestamp` are constants in the failures.** Live they are a fresh uuid and the moment of the failure. The live cases assert their shape with a regex rather than their value.
- **The 404's `path` is right where the runtime supplies it.** An unrouted path receives the path it was asked for; the wrong-method failure is pinned to the path its case uses.
- **The success fixture is document-derived.** Listing products needs a real key; the records here are `ProductEntity`'s own fields with the reference's example values, and every case reading them is marked documentation-only.
- **Nothing is mapped in detection.** Creem is reached through `@creem_io/nodejs` or a plain HTTP call carrying an `x-api-key` header, and neither resolves to this host through a dependency file. Checked 2026-09-13.
