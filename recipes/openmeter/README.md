# openmeter

Emulates the OpenMeter Cloud meters API for local development and tests.

**11 conformance cases, 6 checked against the live API on 2026-09-13.**

Read from the OpenAPI document OpenMeter publishes in its own repository at [`api/openapi.cloud.yaml`](https://raw.githubusercontent.com/openmeterio/openmeter/main/api/openapi.cloud.yaml), and struck live against `openmeter.cloud` on 2026-09-13 with no credential, with a wrong bearer token, and with a Basic header.

## What this Recipe found

**The one field the contract requires is the one the server leaves out.** `UnexpectedProblemResponse` lists `required: [type, title, detail, instance]`.

```
404  {"type":"about:blank","title":"Not Found","status":404,"instance":"urn:request:45fc..."}
405  {"type":"about:blank","title":"Method Not Allowed","status":405,"instance":"urn:request:40a8..."}
```

No `detail`. And `status`, which every one of those documents does send, is not in the required list. The schema requires four fields and gets three, then volunteers a fourth it never promised.

**`invalid token: invalid token`.** A wrong bearer answers:

```
failed to validate security requirements: invalid token: invalid token
```

A Go error wrapped around itself — the same three words twice, with the name of the middleware that raised it in front. The sentence a client prints is the call stack.

**Every problem type is the one that means "no type".** `type` is declared `format: uri` with `default: about:blank`, and all five live documents send exactly that. RFC 9457 reserves `about:blank` for a problem with no semantics beyond its status code, so the machine-readable field is present on every failure to say there is nothing machine-readable here. (The schema's own description cites RFC 7807, which 9457 obsoleted in July 2023.)

**Three sentences behind one title.** All three 401s are `"title": "Unauthorized"` and differ only in `detail`:

```
(no header)      failed to validate security requirements: no matching security requirement
Bearer notreal   failed to validate security requirements: invalid token: invalid token
Basic Zm9v...    failed to validate security requirements: invalid authorization header
```

Three causes, one title — and the field that separates them is the field the 404 and the 405 omit.

**And the two failures with no `detail` are the two the operation never declares.** `listMeters` enumerates 400, 401, 403, 412, 500, 503 and a `default`. The 404 and the 405 this API actually answers are outside its own list, which is why `cauldron drift` reports them as unbacked — and the report is right.

**The identifier pattern accepts two spellings of the same identifier.** `^[0-7][0-9A-HJKMNP-TV-Za-hjkmnp-tv-z]{25}$` — a ULID in Crockford base32, which drops I, L, O and U so a human cannot mistranscribe one, and then permits both cases. `01G65Z755AFWAKHE12NY0CQ9FH` and `01g65z755afwakhe12ny0cq9fh` both validate, are the same identifier, and are not the same string.

**A record can tell you when it was permanently deleted.** `deletedAt` is "Timestamp of when the resource was permanently deleted", carried on a record in your hands, and `includeDeleted` is a query parameter on the listing that puts it there.

**The aggregation rule is a mini-language in a string.** `valueProperty` is `$.tokens` and every value in `groupBy` is another JSONPath expression — typed `string`, described in prose, and parsed by the server. The listing that returns them declares `page` and `pageSize` and answers a bare array, so there is no total and no next.

## Sources

- [`api/openapi.cloud.yaml`](https://raw.githubusercontent.com/openmeterio/openmeter/main/api/openapi.cloud.yaml) — recorded by `cauldron drift`, `spec_seen: 2026-09-13`.
- Live: `openmeter.cloud`, struck 2026-09-13 with no credential, a wrong bearer, and a Basic header.

## Modelling limits

- **Two routes.** Listing meters and creating one. Subscriptions, billing, entitlements, features, notifications, the portal and the event ingest are the rest of a document with over a hundred paths.
- **`instance` is a constant here.** The live field is "A URI reference that identifies the specific occurrence of the problem" and is a fresh 32-hex urn per request. A deterministic sandbox has no occurrences to identify, so it serves one fixed value; the live cases assert the shape with a regex rather than the value.
- **The cookie scheme is not modelled.** `CloudCookieAuth` reads `__session` from a cookie beside the two bearer schemes; this Recipe implements the bearer one, which is what an API client sends.
- **The declared pattern is quoted, not enforced.** Cauldron serves fixture values rather than validating them, so the two-case finding is what the document permits — demonstrated by serving both spellings — and not what the live API stores.
- **Nothing is mapped in detection.** OpenMeter is reached through `@openmeter/sdk`, `openmeter` on PyPI, or the Go client in its own repository, and none of those names resolve to this host through a dependency file. Checked 2026-09-13.
