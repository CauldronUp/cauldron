# cal

Emulates the Cal.com profile endpoint for local development and tests.

**7 conformance cases, 4 checked against the live API on 2026-09-13.**

Read from the OpenAPI document Cal.com publishes at [`cal.com/docs/api-reference/v2/openapi.json`](https://cal.com/docs/api-reference/v2/openapi.json), and struck live against `api.cal.com` on 2026-09-13 with no credential, with a wrong key, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**Every failure says the same sentence twice.**

```json
{"status":"error","timestamp":"…","path":"/v2/me","error":{
  "code":"UnauthorizedException",
  "message":"ApiAuthStrategy - api key - Your api key is not valid",
  "details":{
    "message":"ApiAuthStrategy - api key - Your api key is not valid",
    "error":"Unauthorized","statusCode":401}}}
```

`error.message` and `error.details.message` are byte-identical on every failure checked, so the prose is in the body twice, four levels deep.

**And the status is in it three ways.** `error.code` is a class name with `Exception` appended; `error.details.error` is HTTP's reason phrase; and `error.details.statusCode` is the number — beside a status line that already said it.

**The message begins with the name of the class that raised it.** `ApiAuthStrategy - api key - Your api key is not valid`: the guard's class name, then which branch of it ran, then the sentence, joined by space-hyphen-space.

**A wrong method is a `NotFoundException`.** `PUT /v2/me` answers `"code":"NotFoundException"` with `"Cannot PUT /v2/me"` — the path that does exist, reported missing, under a code named after a 404.

**The anonymous failure teaches the whole auth contract.** "Either pass an API key as 'Bearer' header or OAuth client credentials as 'x-cal-secret-key' and 'x-cal-client-id' headers" — three header names, in a message, on an API whose published document declares no `securitySchemes` at all.

**The document declares no servers either.** `servers` is `[]`, so a generated client has an operation for each of its 228 paths and no address to send them to.

**And `/v2/me` declares one response, the 200.** Not the 401 every anonymous caller meets, not the 404 a typo meets — so `cauldron drift` reports all four of this Recipe's failures as unbacked, and each report is right.

**`status` is an enum of `["success", "error"]` on the success envelope.** The same key carries the same two values on failures, so `status` is not what tells them apart — and a 200 may legally say `"error"`.

**`timeFormat` is a number.** Twelve or twenty-four, as an integer, on a field whose name says format.

**Twelve fields are `required` on the profile and none is declared nullable** — including `bio`, `avatarUrl` and `organizationId`, which a personal account has none of. And the operation is `MeController_getMe`: the controller class and the method name joined by an underscore.

## Sources

- [`cal.com/docs/api-reference/v2/openapi.json`](https://cal.com/docs/api-reference/v2/openapi.json) — recorded by `cauldron drift`, `spec_seen: 2026-09-13`.
- Live: `api.cal.com`, struck 2026-09-13 with no credential, a wrong key, an unrouted path, and a wrong method.

## Modelling limits

- **One route of two hundred and twenty-eight.** The profile. Bookings, event types, schedules, availability, teams, webhooks, OAuth clients, the platform surface and the atoms are the rest.
- **`timestamp` is a constant in the failures.** Live it is the moment of the failure, to the millisecond.
- **The OAuth credential pair is not modelled.** The anonymous message names `x-cal-secret-key` and `x-cal-client-id` as an alternative to the bearer; this Recipe implements the bearer, which is what an API key sends.
- **The success fixture is document-derived.** Reading a profile needs a real key; the record here is `MeOutput`'s own required field set with values of the declared types, and every case reading it is marked documentation-only.
- **Nothing is mapped in detection.** Cal.com is reached through `@calcom/atoms` or a plain HTTP call carrying a bearer token, and neither resolves to this host through a dependency file. Checked 2026-09-13.
