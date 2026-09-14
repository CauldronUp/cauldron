# uploadthing

Emulates the UploadThing file listing for local development and tests.

**9 conformance cases, 5 checked against the live API on 2026-09-14.**

Read from the types in UploadThing's own SDK, and struck live on 2026-09-14 with no key, with a wrong key, on a version that does not exist, on a path that does not exist, with methods the path does not take, and on the root.

## What this Recipe found

**A missing key is a 400 and a wrong key is a 401.**

```
(no header)                     400  {"error":"Missing API Key","data":400}
x-uploadthing-api-key: wrong    401  {"error":"Invalid API key"}
```

Not sending a credential is a *bad request*; only sending the wrong one is *unauthorized*. A client that retries on 401 and gives up on 400 has the two exactly backwards — the missing header is the recoverable one.

**And the status code is in a field called `data`.** `data` is where a successful response's payload belongs; here it holds the number that is already on the status line. The 401 omits the field entirely rather than sending `"data":401`, so the same envelope has two shapes one status apart.

**The two sentences do not agree about capitals.** "Missing API **K**ey" and "Invalid API **k**ey" — written by the same service, about the same header, one request apart.

**Every mistake made without a key is that 400, and every mistake made with a wrong one is that 401.** A path that does not exist, a version that does not exist (`/v7/listFiles`), a verb the path does not take, the root, and a POST with an empty body all answer according to the credential and nothing else. The credential is judged before the path in both directions, so this API cannot tell a caller their URL is wrong at all.

**The CDN calls a well-formed refusal an error.** `X-Cache: Error from cloudfront` rides on every one of these responses, beside an `X-Amzn-Trace-Id` whose `Sampled=0` says the trace was not recorded.

**A status value has a space in it.** The enum is `"Deletion Pending"`, `"Failed"`, `"Uploaded"`, `"Uploading"` — title-case English, and one of the four is two words, so anything switching on it is comparing display strings.

**The SDK's own example does not match its own return type.** The JSDoc above `listFiles` reads:

```js
const data = await listFiles({ limit: 1 });
console.log(data); // { key: "2e0fdb64-…_image.jpg", id: "2e0fdb64-…" }
```

The declared response is `{hasMore, files: [...]}`. The example prints a single file where the method returns an envelope around an array.

**And the key contains the id and the filename.** `key` is `<uuid>_<original name>`, so the identifier a caller stores carries the name of the file that was uploaded — and renaming the file does not change it.

Also pinned: `id` and `key` are two identifiers for one file, with the first embedded in the second, and `customId` is nullable beside them, making three; `uploadedAt` is a bare number; and the listing pages with `limit` and `offset` and reports only `hasMore`, so nothing says how many files there are.

## Sources

- [`pingdotgg/uploadthing`](https://github.com/pingdotgg/uploadthing) — `packages/uploadthing/src/sdk/index.ts` for the `ListFileResponse` schema and the JSDoc example, and `.../sdk/types.ts` for `ListFilesOptions`.
- [UploadThing API reference](https://docs.uploadthing.com/api-reference/ut-api).
- Live: `api.uploadthing.com`, struck 2026-09-14 with no key, a wrong key, `/v7/listFiles`, an unrouted path, `PUT` and `POST`, and the root.

## Modelling limits

- **No description is published.** UploadThing serves no OpenAPI document at any of the usual paths, so there is no `spec` to fingerprint.
- **The success side is SDK-derived.** Listing files needs a real key; the records here are the SDK's own `ListFileResponse` schema with values of the declared types, and every case reading them is marked documentation-only.
- **The default page size is a guess.** `ListFilesOptions` makes `limit` and `offset` optional and names no default; this Recipe serves 100, which is a choice rather than an observation.
- **No routing errors are declared.** Nothing can reach one: the credential decides every response, which is the finding.
- **One route of many.** The file listing. Uploads, deletes, renames, `getFileUrls`, signed URLs and the callback surface are the rest.
- **Nothing is mapped in detection.** UploadThing is reached through `uploadthing` on npm, which does not resolve to this host through a dependency file. Checked 2026-09-14.
