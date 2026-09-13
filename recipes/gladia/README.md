# gladia

Emulates the Gladia transcription jobs API for local development and tests.

**11 conformance cases, 5 checked against the live API on 2026-09-13.**

Read from the OpenAPI document Gladia serves without a credential at [`api.gladia.io/openapi.json`](https://api.gladia.io/openapi.json), and struck live on 2026-09-13 with no key, with a wrong key, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**A method the route does take is answered as a path that is not there.**

```
GET  /v2/cauldron-nope   404  "Cannot GET /v2/cauldron-nope"
PUT  /v2/transcription   404  "Cannot PUT /v2/transcription"
```

The second path is a real route with a documented GET and POST on it, and the API answers that it cannot be found. Express's default handler, reached because nothing declares a 405 — so a client cannot tell a typo in the path from a typo in the verb. And that 404 appears in no operation's response list: `TranscriptionController_list_v2` declares 200 and 401 and nothing else, so the status a caller is most likely to meet first is outside the contract entirely.

**One error object, two naming conventions.** Every failure carries `statusCode` beside `request_id`, `timestamp` and `path` — camelCase from NestJS's built-in exception filter, snake_case from the fields Gladia added around it, in the same five-key object.

**The debug id is the front of the resource id.** The spec's own examples put `id: "45463597-20b7-4af7-b3b3-f5fb778203ab"` and `request_id: "G-45463597"` on one record: the first block of the uuid with a `G-` in front. Two identifiers, one a substring of the other, both on the thing they name.

**A failed job carries an HTTP status as a record field.** `error_code` is `integer, nullable, example: 500`, on a record delivered inside a 200. The response status says the listing worked; the number that says the job did not is in the body, in the space HTTP already had a field for.

**The record says it is version 2, on the path that already says v2.** `version: integer, example: 2`, on `/v2/transcription`.

**Three links and not one of them goes back.** The envelope is required to carry `first`, `current` and `next`, with `next` nullable on the last page. There is no `previous`, no `last` and no count — so a client can always return to the beginning, can never step back one page, and never learns how many pages there are.

**Three record fields are `object, nullable` with no properties at all.** `file`, `request_params` and `result` — and `result` is what the API is for. A fourth, `post_session_metadata`, is `required` and typed the same way, so the contract promises the key will be present and says nothing about what is under it.

**And the documented message for a 401 is one nobody sends.** The spec's example is "gladia key not found". Live, no header answers "no gladia key provided" and a wrong key answers "gladia user not found" — a user, missing, because a key was wrong, on a path with no users in it.

## Sources

- [`api.gladia.io/openapi.json`](https://api.gladia.io/openapi.json) — served without a credential; recorded by `cauldron drift`, `spec_seen: 2026-09-13`.
- Live: `api.gladia.io`, struck 2026-09-13 with no key, a wrong key, an unrouted path, and a method the route does not take.

## Modelling limits

- **Two routes.** Listing transcription jobs and creating one. Upload, the `pre-recorded` and `live` aliases beside them, the per-job file downloads, `/v1/history` and the public `/v1/models` catalogue are the rest.
- **One kind of job.** `items` is a `oneOf` discriminated on `kind`, mixing `PreRecordedResponse` and `StreamingResponse` in one array. This Recipe serves the pre-recorded shape; the heterogeneous array is recorded here rather than emulated.
- **`request_id`, `timestamp` and `path` are constants in the failures.** Live they are per-request: a fresh `G-` id, the moment of the failure, and the path that was asked for. A declarative Recipe fills them from its declaration, so the live cases assert the shape with a regex where the value cannot be fixed, and `path` is right for the failures these cases provoke.
- **`unknown_route` assumes a GET.** The live message is `Cannot <METHOD> <path>` and the runtime hands an unrouted path only the path, so the Recipe's wording covers the GET case it asserts. The wrong-method failure, which does receive both, is exact.
- **Nothing is mapped in detection.** Gladia is reached through `@gladiaio/sdk` or a plain HTTP call with an `x-gladia-key` header, and neither resolves to this host through a dependency file. Checked 2026-09-13.
