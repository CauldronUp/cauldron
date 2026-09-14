# unifiedto

Emulates the Unified.to HRIS employee listing for local development and tests.

**10 conformance cases, 4 checked against the live API on 2026-09-13.**

Read from the OpenAPI document Unified.to serves without a credential at [`api.unified.to/openapi.json`](https://api.unified.to/openapi.json), and struck live on 2026-09-13 with no credential, with a wrong token, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**A credential that was sent is described as missing.**

```
(no header)      401  {"statusCode":401,"error":"Unauthorized","message":"Missing authentication"}
Bearer notreal   401  the identical body
```

"Missing authentication", to a request carrying a token — so a caller with a revoked key and one with no key at all read the same sentence.

**`raw` means two different things on one operation.** It is a record property, `{"type": "object", "additionalProperties": true}`, holding the untransformed upstream payload. It is also a query parameter, and there it is a string:

> Raw parameters to include in the 3rd-party request. Encoded as a URL component. eg. raw parameters: `foo=bar&zoo=bar` -> `raw=foo%3Dbar%26zoo%3Dbar`

A query string, percent-encoded by hand, inside a query parameter's value.

**Two thirds of the document's named types exist because a property was an object.** 826 of its 1,327 schemas are called `property_<Parent>_<field>` — `property_HrisEmployee_emails`, `property_HrisEmployee_compensation` — one per nested property, generated rather than named.

**Nothing on an employee is required.** `HrisEmployee` has no `required` array at all, so `id` is optional in the contract that describes it — and the listing declares exactly one response, the 200, so every failure this Recipe serves is outside the contract and `cauldron drift` says so.

**A field name concatenates two countries.** `ssn_sin` — the United States' Social Security Number and Canada's Social Insurance Number, joined by an underscore, on a record served from three data regions, one of which is Australia.

**An employee record carries a storage quota, three times.** `storage_quota_allocated`, `storage_quota_used` and `storage_quota_available`, all `number`, on a person — and the third is the first minus the second. And it carries a security posture: `has_mfa`, a boolean, beside `date_of_birth`, `pronouns`, `salutation`, `bio`, `marital_status` and `ssn_sin`.

**The field selector is a closed list of every field.** `fields` is an array whose items are an `enum` naming each property, so adding a field to the record means adding it to the parameter that selects fields.

**The document's title is `"Unified.to  API"`, with two spaces**; the security scheme is called `jwt` and declared `{"type": "apiKey", "name": "authorization", "in": "header"}`, so the value is a bare secret in a lowercase header; the collection path is singular — `/hris/{connection_id}/employee` lists employees; and three servers are declared, North American, European and Australian, with no field on any record saying which answered.

## Sources

- [`api.unified.to/openapi.json`](https://api.unified.to/openapi.json) — served without a credential; recorded by `cauldron drift`, `spec_seen: 2026-09-13`.
- Live: `api.unified.to`, struck 2026-09-13 with no credential, a wrong token, an unrouted path, and a wrong method.

## Modelling limits

- **One route of three hundred and forty-eight.** The HRIS employee listing. Accounting, ATS, CRM, commerce, martech, messaging, storage, payment, ticketing, task and KMS are the rest, each with the same `/{category}/{connection_id}/{object}` shape.
- **One region.** Three servers are declared; this Recipe serves one, and the finding is that no record says which answered.
- **`ssn_sin` is null in the fixture.** The schema declares it a string; this Recipe serves it present and empty rather than inventing a value for it.
- **The success fixture is document-derived.** Listing employees needs a real key and a live connection; the record here is `HrisEmployee`'s own properties with values of the declared types, and every case reading them is marked documentation-only.
- **Nothing is mapped in detection.** Unified.to is reached through `@unified-api/typescript-sdk` or a plain HTTP call carrying an `authorization` header, and neither resolves to this host through a dependency file. Checked 2026-09-13.
