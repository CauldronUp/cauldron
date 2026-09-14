# unstructured

Emulates the Unstructured Platform API workflow listing for local development and tests.

**10 conformance cases, 5 checked against the live API on 2026-09-13.**

Read from the OpenAPI document served without a credential at [`platform.unstructuredapp.io/openapi.json`](https://platform.unstructuredapp.io/openapi.json), and struck live on 2026-09-13 with no header, with a malformed key, with a bearer token, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**The path without the trailing slash redirects to `http`.**

```
GET https://platform.unstructuredapp.io/api/v1/workflows

HTTP/1.1 307 Temporary Redirect
location: http://platform.unstructuredapp.io/api/v1/workflows/
```

The request arrived over TLS and the address it is sent to has none. A 307 preserves the method and the body, so a client that follows redirects automatically repeats the request — credential header included — at a cleartext URL, unless its own stack refuses the downgrade.

**And the slash is the API's convention.** `/api/v1/workflows/`, `/api/v1/sources/`, `/api/v1/destinations/`, `/api/v1/jobs/` and `/api/v1/templates/` all end in one, while every item path under them does not — and `/api/v1/notifications` does not either, so the convention is not consistent inside its own document. The address a person would write by hand is the one that redirects.

**A bearer token is a 500.** The 401 says "provide either Bearer token or API key in header". Sending the first of the two:

```
Authorization: Bearer notarealkey    500  {"detail":"Internal Server Error verifying token"}
```

One of the two credential kinds the failure recommends crashes the verifier, and the crash is reported in the caller's own words.

**A malformed key is answered with typing advice.** "Failed to authenticate API Key: API key is malformed, please type the API key correctly in the header." — "please type … correctly", to a program.

**The document declares no security at all.** `components.securitySchemes` is `null`. The credential appears instead as an ordinary header parameter on every operation, `required: false`, typed `anyOf: [string, null]` — so the contract says the key is optional and may be null, on an API that answers 401 without it, and a generated client's signature makes it optional too.

**A workflow type enum contains a price.** `WorkflowType` is `["basic", "advanced", "platinum", "custom"]`: three descriptions of a workflow's structure and one plan name.

**A listing can be asked for the deleted ones.** `show_only_soft_deleted` is a query parameter, `anyOf: [boolean, null]`, default `false` — so soft-deleted workflows are a filter away, and the flag is nullable as well as boolean. `page` and `page_size` are nullable integers with defaults too, so the paging position has three states and one of them is null.

**And the operation declares a 200 and a 422 and nothing else**, on an API that answered 401, 404, 405 and 500 in the probes above — so `cauldron drift` reports every failure this Recipe serves as unbacked, and each report is right.

**`key` on the record is an unlabelled nullable string beside `id`**, and `sources` and `destinations` are arrays of bare uuids while `workflow_nodes` beside them carries whole objects.

## Sources

- [`platform.unstructuredapp.io/openapi.json`](https://platform.unstructuredapp.io/openapi.json) — served without a credential; recorded by `cauldron drift`, `spec_seen: 2026-09-13`.
- Live: `platform.unstructuredapp.io`, struck 2026-09-13 with no header, a malformed key, a bearer token, an unrouted path, and a wrong method.

## Modelling limits

- **Two routes.** The workflow listing, and the slashless address that redirects to it.
- **The 500 is not served.** Sending `Authorization: Bearer …` crashes the live verifier; this Recipe records the transcript rather than emulating a crash, and the sandbox treats an unknown scheme as an unreadable credential.
- **The redirect target is the live one.** `Location` is the exact value the live response carries, cleartext scheme and all, so a client that follows it in a test goes where it would go in production.
- **The success fixture is document-derived.** Listing workflows needs a real key; the records here are `WorkflowInformation`'s own properties with values from its declared enums, and every case reading them is marked documentation-only.
- **Nothing is mapped in detection.** Unstructured is reached through `unstructured-client` on PyPI or a plain HTTP call carrying an `unstructured-api-key` header, and neither resolves to this host through a dependency file. Checked 2026-09-13.
