# adobesign

Emulates the Adobe Acrobat Sign REST API v6 for local development and tests.

**13 conformance cases, 5 checked against the live API on 2026-09-07.**

Written against Adobe's published OpenAPI document for REST API v6 and struck live against `api.na1.adobesign.com` on 2026-09-07 with no credential, with a deliberately invalid one, with a credential carrying no scheme prefix, and on a path that does not exist.

## What this Recipe found

**Three failures, three codes, and every sentence is true.**

```
(no header)        401 {"code":"NO_AUTHORIZATION_HEADER",
                        "message":"Authorization header not provided"}
Bearer notreal     401 {"code":"INVALID_ACCESS_TOKEN",
                        "message":"Access token provided is invalid or has expired"}
/cauldron-nope     404 {"code":"NOT_FOUND","message":"Resource not found"}
```

One shape, three codes, and each says what actually happened. The first names the **header** rather than the key, so a caller who put the token in a query parameter is pointed at the right place. The second says "invalid **or has expired**", which is the honest version of the confident claim [wave](../wave) makes — Wave tells a request carrying no credential that its authentication expired; Adobe tells a request carrying a wrong one that it might have.

After seven providers in this collection that answer "invalid key" to a request holding none, this is the counterexample.

**And the published document says the body is something else entirely.** The specification's own examples for this endpoint's 401 are:

```json
{"reason": "NO_AUTHORIZATION_HEADER: Authorization header not provided.", "status": 401}
```

Two fields, `reason` and `status`. The code and the sentence **concatenated into one string** by a colon. A full stop at the end.

The wire sends `code` and `message`, with the code and the sentence apart, and no full stop. **Not one of those three things matches.** A client generated from this document reads `body.reason`, gets `undefined`, and has no way to discover why.

The schema beside those examples is `{"type": "object"}` with no properties at all — so the examples are the only description the document offers of this body, and they are wrong about every part of it.

**The document declares the credential optional.** `Authorization` is listed as a header parameter with `required: false`, on an operation whose every response to a request without one is the 401 above. A generator reading this emits a client whose auth argument can be left out.

**And it does not say where the API is.** `servers` is a single entry whose `url` is `/api/rest/v6` — a path with no host. Acrobat Sign shards by region, so `na1`, `eu1`, `jp1` and the rest are all real hosts and the document names none of them. The shard a customer belongs to is not discoverable from the specification at all, and a request to the wrong one does not work.

**`id` and `status` are not required, and `hidden` is.** The document guarantees exactly seven fields on an agreement: `displayDate`, `displayParticipantSetInfos`, `groupId`, `esign`, `hidden`, `latestVersionId` and `name`.

So a conforming response must say whether an agreement is hidden from you, and need not say what it is called by or what is happening to it. **The two fields any client reads first are the two the schema does not promise.**

**The status is written from the reader's point of view.** Thirty values, of which eight begin `WAITING_FOR_MY_`: `WAITING_FOR_MY_SIGNATURE`, `WAITING_FOR_MY_APPROVAL`, `WAITING_FOR_MY_DELEGATION`, `WAITING_FOR_MY_ACKNOWLEDGEMENT`, `WAITING_FOR_MY_ACCEPTANCE`, `WAITING_FOR_MY_FORM_FILLING`, `WAITING_FOR_MY_HOSTING`, `WAITING_FOR_MY_VERIFICATION`.

So the same agreement has a different `status` depending on who asked for it. A cached copy is only true for the account that fetched it, and two participants comparing notes are both right.

**The date format is a Java pattern.** `displayDate` is described as "Format would be `yyyy-MM-dd'T'HH:mm:ssZ`" — a `SimpleDateFormat` pattern, quoted `T` and all, where a JSON Schema would say `format: date-time`. The implementation language is in the specification.

**And a response field documents request behaviour.** `id`: "The unique identifier of the agreement.If provided in POST, it will simply be ignored" — a missing space after the full stop, and a note about what a create does, on a field in a read schema.

**The identifier is sent twice more than it needs to be.** `latestVersionId` is a second opaque identifier on the record and `groupId` a third, both in the same prefixed base64 family as the agreement's own, so four fields on one record look alike and address different things.

## Modelling limits

- **One route.** Listing agreements. Members, documents, signing URLs, form fields, reminders, webhooks and the whole widget surface each want their own evidence.
- **The document's declared 401 is not served.** The wire wins: `code` and `message` are what the API sends, and serving `reason` and `status` because the specification says so would teach the mistake this Recipe exists to record.
- **Nothing is mapped in detection.** Adobe's published clients are Java, .NET and Python; none resolves on npm, Packagist or the Go module proxy under a name identifying this API.
