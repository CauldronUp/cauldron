# hunter

Emulates the Hunter email verifier for local development and tests.

**9 conformance cases, 3 checked against the live API on 2026-09-13.**

Read from the OpenAPI document Hunter serves without a credential at [`api.hunter.io/openapi.json`](https://api.hunter.io/openapi.json), and struck live on 2026-09-13 with no credential, with a wrong key in the query string, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**A hundred and twenty of the document's hundred and thirty-three operations point at a placeholder that says so.** `generic_resource_response` is `{"type": "object", "additionalProperties": true}`, and its description reads:

> Generic envelope for V2 resource operations whose detailed schema is not yet in the contract.

Leads, campaigns, companies, api-keys, bulk operations — every one is typed as "an object", by a named schema that admits it. Thirteen operations have a real one.

**Any failure may arrive in either of two shapes, and one is called legacy.** `error_response` is a `oneOf` of `standard_error_response` — `{errors: [{id, code, details}]}` — and `legacy_error_response`, which is `{error: {}}`: a required key whose schema is the empty object, so its value may be anything at all.

**"No user found for the API key supplied" is said for supplying none.**

```json
{"errors":[{"id":"authentication_failed","code":401,"details":"No user found for the API key supplied"}]}
```

A wrong key in the query string gets the identical body, so a key that is absent and a key that is wrong are the same sentence.

**`id` is the code and `code` is the status.** The field named `id` holds `authentication_failed`; the field named `code` holds the integer 401. And the prose lives in `details`, plural, holding one sentence.

**The secret is a declared query parameter.** `apiKeyQuery` is a first-class security scheme — `{"type": "apiKey", "in": "query", "name": "api_key"}` — beside `apiKeyHeader` and `bearerAuth`. Three schemes for one credential, and one of them puts it in the URL.

**An unrouted path answers a branded HTML page.** `404 Not Found • Hunter`, with a stylesheet, on an API whose every other response is JSON.

**A deprecation notice is a required response field.** The verifier's `data` object requires `_deprecation_notice` — a leading underscore, and prose the contract insists on, on every verification.

**`status` and `result` are both required strings.** Two verdict fields on one record, and nothing says how they differ.

**And `webmail` appears twice in one response.** Once inside `data` and once inside `meta`, at two levels, under one name.

**Nine booleans on the record** — `regexp`, `gibberish`, `disposable`, `webmail`, `mx_records`, `smtp_server`, `smtp_check`, `accept_all`, `block` — with no enum anywhere; `meta.params.email` echoes the request back inside the response; and every one of those objects is `additionalProperties: false`, a closed contract, in a document whose other hundred and twenty operations are open.

## Sources

- [`api.hunter.io/openapi.json`](https://api.hunter.io/openapi.json) — served without a credential; recorded by `cauldron drift`, `spec_seen: 2026-09-13`.
- Live: `api.hunter.io`, struck 2026-09-13 with no credential, a wrong key in the query string, an unrouted path, and a wrong method.

## Modelling limits

- **One route of ninety-seven.** The email verifier, which is one of the thirteen operations the document actually describes. Domain search, email finder, discover, leads, campaigns, companies and the bulk surface are the rest — and a hundred and twenty of them share the placeholder schema.
- **`meta` is served as route constants.** Live, `meta.domain`, `meta.company_name` and `meta.webmail` are derived from the address being checked; the fixture carries the shape with one address's values.
- **The HTML 404 is not served.** An unrouted path answers a branded page from the web application; this Recipe records that rather than emulating it.
- **The success fixture is document-derived.** Verifying an address needs a real key; the record here is `email_verifier_response`'s own required field set with values of the declared types, and every case reading it is marked documentation-only.
- **Nothing is mapped in detection.** Hunter is reached through a plain HTTP call carrying `api_key` or `X-API-KEY`, which does not resolve to this host through a dependency file. Checked 2026-09-13.
