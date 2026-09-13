# bytescale

Emulates the Bytescale recent-jobs listing for local development and tests.

**10 conformance cases, 5 checked against the live API on 2026-09-13.**

Read from Bytescale's own generated TypeScript SDK ([`src/public/shared/generated/models/index.ts`](https://github.com/bytescale/bytescale-javascript-sdk)), and struck live against `api.bytescale.com` on 2026-09-13 with no header, with a bearer value that is not a key, with a bearer value shaped like one, with a Basic header, and on a path that does not exist.

## What this Recipe found

**A nine-letter string is diagnosed as a malformed JWT.** `Authorization: Bearer notarealkey` answers:

```json
{"error":{"message":"Unauthenticated: your JWT must include an 'accountId' field at the same level as the 'exp' field.","code":"jwt_missing_account_id_field","timestamp":"..."}}
```

Nothing in that request was a JWT. The API explains the placement of two claims relative to each other to a caller who sent eleven characters, because "not a key" and "a JWT missing a field" are the same branch.

**And a value that does look like a key gets part of itself back.** `Bearer public_notarealkey` answers `authentication_failed_account_not_found` with `details: {"accountId":"notarea"}` — the seven characters after the prefix, which is how Bytescale derives an account id from a key. The refusal quotes a slice of the credential it just refused, and that same slice is in every URL: `/v2/accounts/{accountId}/jobs`. The account id is a prefix of the secret, so every access log that records a path records it.

**A third header shape gets a lecture about usernames.** `Basic Zm9vOmJhcg==` answers "Unauthenticated: incorrect username. Please provide 'apikey' as your username and your API key as the password." So Basic is supported, with a constant username, documented only in the failure you get for not knowing it.

**An unrouted path is answered with an instruction that ends in "etc."**

```json
{"error":{"message":"Check the URL for this GET request. (Hint: Ensure the account ID in the URL is valid, etc.)","code":"check_path","details":{"path":"/v2/cauldron-nope"}}}
```

The message is a to-do addressed to the caller, the hint is parenthesised, and the list of things to check is left unfinished.

**The 401 sends both halves of a CORS pair browsers refuse together.** `access-control-allow-origin: *` and `access-control-allow-credentials: true` arrive on the same response. A browser rejects that combination outright, so the permissive one is exactly what makes the credentialed one useless.

**Every job record carries the same URL, and it is a link to the documentation.** `jobDocs` is typed `"https://www.bytescale.com/docs/job-api/GetJob"` — an enum with one member — and described as "Link to the documentation that describes how to get a job's status from its job ID". Documentation as a record field, identical on every record in every response.

**A timestamp's generated type is an object with a property called `tYPE`.** `created` and `lastUpdated` are `OpaqueNumberEpochMillis`, which the published SDK defines as:

```ts
export interface OpaqueNumberEpochMillis {
  tYPE: OpaqueNumberEpochMillisTYPEEnum;   // "EpochMillis"
}
```

A branded wrapper for a value that arrives on the wire as a number, whose single property is `TYPE` with its first letter lowercased by the generator. It is declared twice, identically, as `OpaqueNumberEpochMillis` and `OpaqueNumberEpochMillisAllOf`.

**Two listings on one API, two envelopes, and only one can be paged.** `ListRecentJobsResponse` is `{items}` and nothing else. `ListFolderResponse` beside it is `{cursor, folder, isPaginationComplete, items}` — and its `cursor` is documented as "Absolute path to a file or folder. Begins with a `/`", so the paging token is a filesystem path and the comment says so.

**And a failed job carries the whole failure envelope inside a 200.** `JobSummaryError` is `{timestamp, message, details, code}` — the same four keys the HTTP failures use, nested on a record in a successful response.

**`MultipartUploadProtocol` is `"1.0" | "1.1"`**, and the comment for `1.1` says it "fixes a known issue in the `2.0` protocol" — naming a version the type does not contain, for an issue the comment above attributed to `1.0`.

## Sources

- [`bytescale/bytescale-javascript-sdk`](https://github.com/bytescale/bytescale-javascript-sdk) — `src/public/shared/generated/models/index.ts`, the generated models.
- Live: `api.bytescale.com`, struck 2026-09-13 with no header, a non-key bearer, a key-shaped bearer, a Basic header, and an unrouted path.

## Modelling limits

- **One route.** The recent-jobs listing. Files, folders, the upload protocols, copy and delete batches and the transformation surface are the rest.
- **The echoed account id is a constant.** Live it is the seven characters after the prefix of whatever key was presented. A declarative Recipe fills it from its declaration, so the case that asserts it sends the key those seven characters come from.
- **`timestamp` is a constant in the failures.** Live it is the moment of the failure, to the millisecond.
- **No `spec:`.** Bytescale publishes a generated SDK rather than the document it was generated from, and no OpenAPI description is served at any address this Recipe could find — so there is nothing for `cauldron drift` to record.
- **The success fixture is SDK-derived.** Listing jobs needs a real key; the records here are `JobSummary`'s own fields with values of the declared shapes, and every case reading them is marked documentation-only.
- **Nothing is mapped in detection.** Bytescale is reached through `@bytescale/sdk` or a plain HTTP call carrying an `Authorization` header, and neither resolves to this host through a dependency file. Checked 2026-09-13.
