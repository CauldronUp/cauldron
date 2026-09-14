# kombo

Emulates the Kombo unified HRIS employee listing for local development and tests.

**9 conformance cases, 3 checked against the live API on 2026-09-13.**

Read from the OpenAPI document Kombo serves without a credential at [`api.kombo.dev/openapi.json`](https://api.kombo.dev/openapi.json), and struck live on 2026-09-13 with no credential, with a wrong token, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**An employee listing is required to carry identifying details.** The response schema's `required` array includes `ssn` ("The employee's social security number"), `tax_id`, `date_of_birth` and `nationality`, and the record also carries `home_address`, `bank_accounts`, `gender`, `ethnicity` and `marital_status`. The operation's query parameters are all filters; there is no field-selection parameter among them, so every caller who can list employees receives all of it.

**Two failure shapes under one envelope, and three keys vanish.**

```
no credential  401  {"status":"error","error":{"code":"PLATFORM.AUTHENTICATION_INVALID","title":"…","message":"No auth token specified! …","log_url":null}}
bad path       404  {"status":"error","error":{"message":"Can not GET /v1/cauldron-nope"}}
```

`code`, `title` and `log_url` are on the first and absent from the second, so the field a client switches on disappears exactly when the path is wrong. And the listing declares a 200 and a `default` and nothing else, so both of those statuses are outside the contract — which is why `cauldron drift` reports all three of this Recipe's failures as unbacked.

**`log_url` is always present and always null.** A field for a link to your own logs, on every authentication failure, filled in by nothing.

**The refusals use exclamation marks and ask a question.** "No auth token specified! Please include your API key in your headers like this: `Authorization: Bearer <token>`" and "The specified auth token does not seem valid! Is the region (us/eu) in the subdomain correct?" — the second puts a question to the caller, and blames geography for a token that was simply wrong. And the routing failure says "Can not", in two words.

**A filter can be silently ignored, by request.** `ignore_unsupported_filters` is a query parameter, so a caller can ask for a narrowed listing and receive a wider one with nothing in the response saying which filter was dropped. Beside it, `employment_status` and `employment_statuses` are two parameters for one field, and `include_deleted` puts deleted employees back in.

**Every record carries the time it was deleted.** `remote_deleted_at` is in `required`, on a listing of people who work there.

**And the raw upstream payload rides along.** `remote_data` is a required property, so each unified record carries the vendor-specific object it was derived from — which is the thing a unified API exists to hide.

**The success schema is named `…PositiveResponse`** rather than success or ok; `changed_at` is the only timestamp on the record that is not nullable; and the document's first server is the EU host, while the message for a wrong token asks whether you meant the other one.

## Sources

- [`api.kombo.dev/openapi.json`](https://api.kombo.dev/openapi.json) — served without a credential; recorded by `cauldron drift`, `spec_seen: 2026-09-13`.
- Live: `api.kombo.dev`, struck 2026-09-13 with no credential, a wrong token, an unrouted path, and a wrong method.

## Modelling limits

- **One route of a hundred and eleven.** The HRIS employee listing. ATS, assessment, AI-apply, time off, payroll, absences, attendance, webhooks and the integration-management surface are the rest.
- **One region.** The document declares an EU host and a US one; this Recipe serves one, and the finding is that the wrong-token message asks which you meant.
- **The sensitive fields are null in the fixture.** The schema requires `ssn`, `tax_id`, `home_address` and `bank_accounts` to be present; this Recipe serves them present and empty, which is what the contract permits, and does not invent values for them.
- **The success fixture is document-derived.** Listing employees needs a real key and a connected integration; the records here are the operation's own `results` schema with values of the declared types, and every case reading them is marked documentation-only.
- **Nothing is mapped in detection.** Kombo is reached through a plain HTTP call carrying a bearer token and an `X-Integration-Id` header, which does not resolve to this host through a dependency file. Checked 2026-09-13.
