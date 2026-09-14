# docraptor

Emulates the DocRaptor asynchronous document job for local development and tests.

**8 conformance cases, 5 checked against the live API on 2026-09-14.**

Read from the API documentation at [`docraptor.com/documentation/api`](https://docraptor.com/documentation/api), and struck live on 2026-09-14 with no credential, with a wrong basic credential, on a path that does not exist, with a method the path does not take, and on a path one underscore away from a real one.

## What this Recipe found

**The failures are XML and the successes are JSON.**

```
Content-Type: application/xml; charset=utf-8

<?xml version="1.0" encoding="UTF-8"?>
<errors>
  <error>Please provide a valid API key.</error>
</errors>
```

while `GET /status/{id}` answers `{"status":"completed", …}`. One API, two serialisations, split by whether the request worked — so a client calling `.json()` on its own mistake throws. The XML declaration, the wrapper element and the charset are all there to carry nine words.

**A missing key and a wrong key are the same sentence.** "Please provide a valid API key." — to a caller who provided one, and to a caller who provided none.

**A mistyped path gets the marketing site.** `GET /cauldron-nope`, `PUT /docs` and `GET /doc_status/abc123` — one underscore from the real `/status/abc123` — all answer `404` with `text/html` and DocRaptor's own homepage markup, `<meta name="description" content="Simple Excel and PDF document generation webservice.">` included.

**A list of errors is one string with newlines in it.**

```json
{"status":"failed",
 "validation_errors":"Name can't be blank\nName is too long (maximum is 200 characters)"}
```

Two errors, joined by a newline, in a field whose name is plural. A caller who wants them separately splits on `\n` and hopes no message contains one.

**And the completion time is an English sentence.** `"message": "Completed at Mon Jun 06 18:33:17 +0000 2011"` — Ruby's default `Time#to_s`, prefixed with a word, in a field called `message`. There is no machine-readable finish time anywhere on the record.

**The success and the failure share no fields but one.** A completed job carries `download_url`, `message`, `number_of_pages` and `status`. A failed one carries `status` and `validation_errors`. The only key on both is `status`, so every other field is optional in practice and none is documented as such.

Also pinned: the asynchronous create answers `{"status_id":"123454321"}` and nothing else, so the only thing a caller gets back from submitting a document is a numeric string; `status` moves through `queued`, `working`, `completed` and `failed`, with the first two carrying none of the completion fields; a job has 600 seconds to finish; and a `download_url` may be good for exactly one download depending on the account's retention policy, which nothing on the record says.

## Sources

- [DocRaptor API documentation](https://docraptor.com/documentation/api) — the async flow, the status fields, and the failure body.
- Live: `docraptor.com`, struck 2026-09-14 with no credential, a wrong Basic credential, an unrouted path, a wrong method, `/status/abc123`, and `/doc_status/abc123`.

## Modelling limits

- **No description is published.** DocRaptor serves no OpenAPI document, so there is no `spec` to fingerprint.
- **The marketing 404 is abridged.** Live it is DocRaptor's full homepage; this Recipe serves a minimal document carrying the same status, the same content type and the same `description` meta tag, which is what identifies it.
- **The synchronous document is not served.** `POST /docs` without `async` returns the PDF bytes themselves; this Recipe models the asynchronous flow, where the same path answers a JSON status id.
- **The status side is document-derived.** Reading a job needs a real key; the records here are the documentation's own completed and failed examples, and the cases reading them are marked documentation-only.
- **A status id that is not there is served as the credential failure.** That is what live answers to an unauthenticated caller; what it answers to an authenticated one asking for a job that does not exist was not observed.
- **Two routes of a handful.** Submitting a document and reading its status. The download, the synchronous path and the callback are the rest.
- **Nothing is mapped in detection.** DocRaptor is reached through `docraptor` on RubyGems, npm or PyPI, and none of them resolves to this host through a dependency file. Checked 2026-09-14.
