# convertapi

Emulates the ConvertAPI account endpoint for local development and tests.

**8 conformance cases, 5 checked against the live API on 2026-09-14.**

Read from the documentation at [`convertapi.com/docs/user-information`](https://www.convertapi.com/docs/user-information), and struck live on 2026-09-14 with no credential, with a wrong bearer token, with the secret as a query parameter, on a path that does not exist, with a method the path does not take, on a conversion path, and on the root.

## What this Recipe found

**The root tells an anonymous caller their own IP address and the build number.**

```
GET /
200  text/plain

Welcome <the caller's public IP> to ConvertAPI 1.2.67.108
```

No credential, no path, no parameters — a greeting carrying the address the request came from and a four-part internal version number. It is the only 200 this API gives away.

**The error code is the HTTP status with a digit on the end.** `4011`. A client storing the vendor code and a client storing the status are storing the same number, one of them with an extra digit, and `4011` sorts nowhere near the 401 it means.

**The message admits it cannot tell the two cases apart.** "Unauthorized. Invalid **or missing** API credentials." Struck live, the identical body answers no credential, a wrong bearer token, and a wrong secret in the query string.

**The keys are PascalCase.** `{"Code":4011,"Message":"…"}` — capitalised where almost every JSON API in this collection is not, and the same casing runs through the success: `Active`, `FullName`, `ConversionsTotal`.

**And the routing failures are not JSON at all.** An unrouted path answers `404` with an empty body and no `Content-Type`. A method the path does not take answers `405` the same way. So the credential failure is a typed document, and the two failures that would say what was wrong with the *request* say nothing.

**One endpoint needs a stronger credential than the rest of the API.** The documentation says `/user` requires "the Master Token" and that "regular API Tokens are not accepted" — and a regular token draws the same `4011` as no token at all, so a caller holding a working key for every other endpoint cannot tell this one apart from a broken credential.

**The documentation declines to say what the response contains.** "The response may include additional legacy fields kept for backward compatibility" — the five documented fields are a lower bound, and what else arrives is whatever an older plan once needed.

Also pinned: `ConversionsTotal` and `ConversionsConsumed` are two counters with no *remaining* between them, so the number a caller actually wants is a subtraction they do themselves; and the account holder's full name and email ride on what is otherwise a usage endpoint.

## Sources

- [ConvertAPI user information](https://www.convertapi.com/docs/user-information) — the response fields and the Master Token requirement.
- Live: `v2.convertapi.com`, struck 2026-09-14 with no credential, a wrong bearer token, a wrong `Secret` query parameter, an unrouted path, a wrong method, `/convert/pdf/to/jpg`, and the root.

## Modelling limits

- **The root's IP address is a documentation one.** Live, the greeting carries the caller's own public address. The fixture uses `192.0.2.1` — RFC 5737's TEST-NET-1 — because a Recipe must not carry anybody's real address, and the finding is that the field is there at all.
- **The build number is the one observed on 2026-09-14.** It will have moved by the time you read this; the shape is the point.
- **No description is published.** ConvertAPI serves no OpenAPI document, so there is no `spec` to fingerprint.
- **The success side is document-derived.** Reading the account needs a real Master Token; the record here is the documentation's own example, and the cases reading it are marked documentation-only. The "additional legacy fields" the documentation mentions are not served, because it does not name them.
- **Two routes of many.** The root and the account. The whole conversion surface — `/convert/{from}/to/{to}` in every combination — is the rest, and it answers the same `4011`.
- **Nothing is mapped in detection.** ConvertAPI is reached through `convertapi` on npm, PyPI, NuGet or RubyGems, and none of them resolves to this host through a dependency file. Checked 2026-09-14.
