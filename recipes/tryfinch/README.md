# tryfinch

Emulates the Finch employee directory for local development and tests.

**9 conformance cases, 4 checked against the live API on 2026-09-13.**

Read from the OpenAPI document Finch serves without a credential at [`api.tryfinch.com/openapi.json`](https://api.tryfinch.com/openapi.json), and struck live on 2026-09-13 with no header at all, with a wrong bearer, on a path that does not exist, with a method the path does not take, and with a version header Finch does not know.

## What this Recipe found

**Nothing gets past the version header.** With no `Finch-API-Version` on the request:

```
GET /employer/directory           400  Missing 'Finch-API-Version' header
GET /employer/directory + token   400  Missing 'Finch-API-Version' header
GET /cauldron-nope                400  Missing 'Finch-API-Version' header
PUT /employer/directory           400  Missing 'Finch-API-Version' header
```

A wrong path, a wrong method and a wrong credential are all indistinguishable until the header is right. And the header is declared in the document as a *parameter*, `required: true`, with `default: "2020-09-17"` — so a generated client fills it in and a hand-written one does not, and the two have entirely different first experiences of the API. Its pattern, `([12]\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01]))`, is unanchored, so any value containing a date anywhere passes the declared check.

**The 401 tells you the name of the middleware.** A wrong token answers "AUTHENTICATION MIDDLEWARE: could not match any resolver for the token format" — capitals, a colon, and an internal word ("resolver") that appears nowhere in the documentation.

**And the vendor's own error code is missing exactly from the credential failures.**

```
400  name: invalid_request_error         finch_code: invalid_request
404  name: not_found_error               finch_code: item_not_found
401  name: unauthorized_request_error    (no finch_code)
401  name: invalid_token_error           (no finch_code)
```

The field a client switches on to stay version-independent is absent from the two failures every integration meets first. Where it is present, `name` and `finch_code` say the same word with and without `_error` — except the 404, where one says "route" and the other says "item".

**And the operation declares 200, 202 and 422, so neither the 401 nor the 404 is in the contract at all.** `cauldron drift` reports all four of this Recipe's credential and routing failures as unbacked, and each report is right.

**A 202 on a GET is shaped exactly like a failure.** The document's own example:

```json
{"code": 202, "name": "sync_in_progress", "finch_code": "data_sync_in_progress", "message": "The data being requested is being fetched. Please check back later."}
```

The same four keys as every error, on a successful status. Code that decides "this body has `code` and `name`, so it is an error" is right every time except this one.

**The identifier pattern accepts a uuid with no hyphens, and any string containing one.** `[0-9a-fA-F]{8}[-]?(?:[0-9a-fA-F]{4}[-]?){3}[0-9a-fA-F]{12}` is unanchored and every hyphen in it is optional, under a description that reads "A stable Finch `id` (UUID v4)".

**Every field but the identifier is nullable.** `first_name`, `middle_name`, `last_name`, `manager`, `department` and `is_active` are each `nullable: true`, so an individual in the directory may have no name — and `is_active` is a boolean with three states, on "is this person employed".

**The count is optional and the offset is not.** `Paging` requires `offset` and leaves `count`, "the total number of elements for the entire query", out of `required`.

**Nine `/sandbox/` paths sit in the production document**, at the production host, beside the real ones; `/employer/individual` and `/employer/employment` are POSTs that read; and the document's version is the same date the header wants, so the API's version and its request header are one string in two places.

## Sources

- [`api.tryfinch.com/openapi.json`](https://api.tryfinch.com/openapi.json) — served without a credential; recorded by `cauldron drift`, `spec_seen: 2026-09-13`.
- Live: `api.tryfinch.com`, struck 2026-09-13 with no header, a wrong bearer, an unrouted path, a wrong method, and an unknown version.

## Modelling limits

- **One route of forty-four.** The employee directory. Company, employment, individual, pay statements, benefits, plans, payments, jobs, providers, the connect flow and the sandbox surface are the rest.
- **The version header is checked after the credential here, and live it is checked first.** That ordering is the finding above; a declarative Recipe authenticates first, so the case that provokes the missing-header failure presents a good credential, and the live transcript for the other three combinations is recorded rather than emulated.
- **The success fixture is document-derived.** Reading a directory needs a real access token; the records here are the operation's own `individuals` schema with values of the declared shapes, and every case reading them is marked documentation-only.
- **The 202 is not served.** A sync still running answers a success status with an error-shaped body; a deterministic sandbox has no sync to be waiting on, so the shape is recorded here rather than emulated.
- **Nothing is mapped in detection.** Finch is reached through `@tryfinch/finch-api` on npm or `finch-api` on PyPI, and neither resolves to this host through a dependency file. Checked 2026-09-13.
