# schematic

Emulates the Schematic company listing for local development and tests.

**10 conformance cases, 4 checked against the live API on 2026-09-14.**

Read from the API reference at [`docs.schematichq.com`](https://docs.schematichq.com/api-reference/companies/list-companies) and from the types in Schematic's own Go SDK, and struck live on 2026-09-14 with no key, with a wrong key, with an empty header, on a path that does not exist, with a method the path does not take, and on the root.

## What this Recipe found

**This API says almost everything twice, once singular and once plural.** On the company record:

| singular | plural |
|---|---|
| `plan` (optional) | `plans` (always sent) |
| `billing_subscription` (optional) | `billing_subscriptions` (always sent) |
| `billing_profile` (optional) | `billing_profiles` (optional) |

Three pairs, each an optional singular beside an array of the same thing, and nothing saying whether the singular is the first of the array, the current one, or something else. `traits` and `entity_traits` are a fourth pair in two different shapes — a `map[string]any` of names to values, and an array of trait objects.

**And the listing's parameters do it too.** Twenty-one query parameters, among them `plan_id` and `plan_ids`, `plan_version_id` and `plan_version_ids`, and three pairs of opposites: `with_subscription` and `without_subscription`, `without_plan`, `with_entitlement_for` beside `without_feature_override_for`. Nothing documents what happens when a caller sends both halves of a pair.

**The listing has no count and no cursor.** The envelope is `{data, params}` — the records, and the request's own parameters echoed back. Paging is `limit` (default 100) and `offset` (default 0), and nothing in a full page says whether another exists.

**A missing key and a wrong key say different things; an empty header counts as missing.**

```
(no header)                       401  {"error":"Unauthorized"}
X-Schematic-Api-Key: (empty)      401  {"error":"Unauthorized"}
X-Schematic-Api-Key: notrealkey   401  {"error":"API key invalid"}
```

Two distinct sentences, which is rarer than it should be — most of the APIs in this collection have one — and the header present but empty is judged absent rather than malformed.

**A wrong path and a wrong method are Go's default page, in plain text**, answered before the credential: `404 page not found`, `Content-Type: text/plain`, from an API that is JSON everywhere else.

**And the API host serves a web application at its root.** `GET /` answers `200` with `text/html` and `<title>Schematic</title>`.

**Money is a floating-point number.** `billing_credit_balances` is `map[string]float64` — credit balances, per credit type, in binary floating point, on a billing record.

Also pinned: the generated Go SDK gives every struct a private `explicitFields *big.Int`, commented "Private bitmask of fields set to an explicit value and therefore not to be omitted", so a caller's intent to send a zero travels as a big integer beside the data; and a company record carries `pending_migration` and `scheduled_downgrade`, two pieces of billing-process state, on the same object as its name.

## Sources

- [List companies](https://docs.schematichq.com/api-reference/companies/list-companies) — the `{data, params}` envelope, the `X-Schematic-Api-Key` header, and the twenty-one query parameters.
- [`schematichq/schematic-go`](https://github.com/schematichq/schematic-go) — `types.go`, for `CompanyDetailResponseData` and the singular/plural pairs.
- Live: `api.schematichq.com`, struck 2026-09-14 with no key, an empty header, a wrong key, an unrouted path, a wrong method, and the root.

## Modelling limits

- **No description is published.** Schematic serves no OpenAPI document at any of the usual paths, so there is no `spec` to fingerprint.
- **The success side is SDK-derived.** Listing companies needs a real key; the records here are the Go SDK's own `CompanyDetailResponseData` fields with values of the declared types, and every case reading them is marked documentation-only.
- **Sixteen fields of a much wider record.** `CompanyDetailResponseData` also carries `add_ons`, `custom_plan_billings`, `default_payment_method`, `entitlements`, `keys`, `metrics`, `payment_methods`, `pending_migration`, `rules` and `scheduled_downgrade`.
- **Two of twenty-one parameters are served.** `limit` and `offset` page; `q` filters on the name. The singular/plural and with/without pairs are the finding, and serving them would mean inventing behaviour for combinations the reference does not describe.
- **`params` is served with the defaults, not the request.** Live it echoes what was asked for; this format adds constants to an envelope rather than reflecting the query, so the case asserting it says so.
- **One route of many.** The company listing. Companies by key, users, features, entitlements, plans, components and the event surface are the rest.
- **Nothing is mapped in detection.** Schematic is reached through `schematic-go`, `@schematichq/schematic-js` or a plain header call, and none of them resolves to this host through a dependency file. Checked 2026-09-14.
