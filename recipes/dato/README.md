# dato

Emulates the DatoCMS Content Management API record listing for local development and tests.

**12 conformance cases, 4 checked against the live API on 2026-09-13.**

Read from the JSON Hyper-Schema DatoCMS serves without a credential at [`site-api.datocms.com/docs/site-api-hyperschema.json`](https://site-api.datocms.com/docs/site-api-hyperschema.json), and struck live on 2026-09-13 with no credential, with a wrong token, with no `Accept` header, with the wildcard one, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**Errors arrive in the JSON:API `data` member.** Every failure is shaped like this:

```json
{"data":[{"id":"02f293","type":"api_error","attributes":{
  "code":"INVALID_SITE","details":{},
  "doc_url":"https://www.datocms.com/docs/content-management-api/errors#INVALID_SITE"}}]}
```

JSON:API reserves `errors` for exactly this, and says a document's `data` and `errors` members "MUST NOT coexist". Here the errors *are* the data, typed `api_error` — so a client checking for an `errors` key finds none, and a client unwrapping `data` into its record type gets an error object instead.

**The per-occurrence id is not per occurrence.** JSON:API defines `id` on an error object as a unique identifier for that particular occurrence of the problem. Struck live:

| request | id | code |
|---|---|---|
| no `Authorization`, ×3 | `02f293` | `INVALID_SITE` |
| three *different* wrong tokens | `5413bc` | `INVALID_AUTHORIZATION_HEADER` |
| unrouted path | `cf3acb` | `INVALID_ENDPOINT` |
| wrong method on a real path | `cf3acb` | `INVALID_ENDPOINT` |

It is a constant per error kind, not per occurrence. The identifier that really is per request is `x-request-id`, in a header, beside `x-runtime` and `x-queue-time`.

**A missing credential is reported as an invalid site.** No `Authorization` header at all answers `INVALID_SITE` — the token is what names the site, so its absence is reported as the site being wrong rather than the credential being absent. A *wrong* token answers `INVALID_AUTHORIZATION_HEADER`, which blames the header for the value inside it.

**`Accept: */*` is not acceptable — and where that is checked depends on what else is wrong.** The header every HTTP client sends by default draws a 406 `INVALID_ACCEPT_HEADER`, and so does sending no `Accept` header at all. But only for a caller with no token. Struck live, all with no `Accept` header:

```
GET /items          (no token)      406  INVALID_ACCEPT_HEADER
GET /items          (wrong token)   401  INVALID_AUTHORIZATION_HEADER
GET /cauldron-nope  (no token)      404  INVALID_ENDPOINT
```

The `Accept` check sits between "your token is wrong" and "you have no token": a bad credential is reported before it, a missing one after it, and a bad path before both. `application/vnd.api+json` passes it; `application/json` passes it; `*/*` does not.

**A wrong method is reported as a wrong endpoint.** `PUT /items` and `GET /cauldron-nope` answer the same 404 `INVALID_ENDPOINT`, the same body, the same id.

**The path parameter is a percent-encoded JSON Pointer.** In the schema, one record lives at:

```
/items/{(%2Fschemata%2Fitem%23%2Fdefinitions%2Fitem%2Fdefinitions%2Fidentity)}
```

A URI Template variable whose name is a pointer into the document, url-encoded, sixty characters long, for what is in the end a record id. Every one of the resource's eleven single-record links carries it.

**The record's content is untyped.** `attributes` is "The JSON data associated to the record", declared `{"type": "object"}` with no `properties` and no `additionalProperties` — the entire payload of a CMS record, unconstrained by the schema that describes it.

**The id is a UUID that is not written as one, and the schema's own example disagrees.** `id` is "RFC 4122 UUID of record expressed in URL-safe base64 format", example `hWl-mnkWRYmMCSTq4z_piQ`. The example `attributes` printed directly beside it carries `"category": "24"` and `"upload_id": "20042921"` — short decimal strings, which is what DatoCMS ids used to be and what a reader copying the example will code against.

**A parameter's description names a parameter that does not exist.** `locale` reads "When `filter[query]` or `field[fields]` is defined, filter by this locale." The listing's parameters are `nested`, `filter`, `locale`, `page`, `order_by` and `version`. There is no `field`.

**And filters are documented to be ignored rather than refused.** From the listing's own description: some combinations "may be invalid, resulting in errors or ignored filters, which can lead to unexpected results". Three such combinations are listed, and which of the two outcomes each produces is not.

Also pinned: the resource is `item` on the wire and titled "Record" in the schema that defines it; three validity booleans ride on one record — `is_valid`, `is_current_version_valid` and `is_published_version_valid`, the last two nullable and the first not; `current_version` is "The ID of the current record version", a field named version holding an id, beside a `version` query parameter meaning published-or-current and a separate `item_version` resource; `relationships.creator.data` is an `anyOf` of five kinds including `access_token`, so a record's creator may be a credential; the listing's `meta` is `{total_count}` with `additionalProperties: false`, so nothing in the body says where the page ended; and the document declares its dialect as `http://interagent.github.io/interagent-hyper-schema`, over cleartext `http`.

## Sources

- [`site-api.datocms.com/docs/site-api-hyperschema.json`](https://site-api.datocms.com/docs/site-api-hyperschema.json) — 1.16 MB, `DatoCMS Site API`, 53 resources, served without a credential.
- [DatoCMS Content Management API documentation](https://www.datocms.com/docs/content-management-api).
- Live: `site-api.datocms.com`, struck 2026-09-13 with no credential, a wrong token, a missing `Accept` header, `Accept: */*`, `Accept: application/vnd.api+json`, an unrouted path, and a wrong method.

## Modelling limits

- **The description cannot be fingerprinted.** It is a JSON Hyper-Schema, not OpenAPI, and `cauldron drift` says so: "no openapi or swagger version field, so this is probably not an OpenAPI description". `spec_hash` is recorded as `-`, which is this format's way of saying the document exists and cannot be read.
- **The 406 is not served.** Live, the `Accept` check runs after a *wrong* credential and before a *missing* one. This Recipe checks the credential first in both cases, so there is no request that reaches the header check, and declaring an error nothing can reach would be worse than describing the ordering here. The three live probes above are the record of it.
- **Two routes of sixteen on this resource alone.** The listing and one record. Create, update, delete, duplicate, publish, unpublish, validate (twice), references, current-vs-published-state and four bulk operations are the rest — and `item` is one of 53 resources in the document.
- **`page[limit]` and `page[offset]` are served; `nested`, `filter`, `order_by` and `version` are not.** The response-mode switch (`nested=true`, which replaces block ids with block objects) changes the shape of `attributes`, and this Recipe serves one shape.
- **The fixture records are document-derived.** Listing records needs a real token and a real project; the records here are the schema's own `item` definition with its own example values, and every case reading them is marked documentation-only.
- **Nothing is mapped in detection.** DatoCMS is reached through `@datocms/cma-client` on npm or a plain bearer call, and neither resolves to this host through a dependency file. Checked 2026-09-13.
