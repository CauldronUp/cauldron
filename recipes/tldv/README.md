# tldv

Emulates the tl;dv meeting listing for local development and tests.

**12 conformance cases, 6 checked against the live API on 2026-09-13.**

Read from the REST documentation at [`doc.tldv.io`](https://doc.tldv.io/), and struck live on 2026-09-13 with no key, with a wrong key, on a path that does not exist, with a method the path does not take, and on the health endpoint.

## What this Recipe found

**A caller who sends a key is told authorization is required.**

```
GET /v1alpha1/meetings   (no header)              401
GET /v1alpha1/meetings   x-api-key: notarealkey   401

{"name":"AuthorizationRequiredError","message":"Authorization is required for request on GET /v1alpha1/meetings"}
```

Byte for byte the same body. A revoked key, a typo'd key and a forgotten key all read a sentence about the header being absent, so nothing in the answer tells a caller which of the three they did.

**And the sentence echoes the request back.** The message is composed from the method and the full path:

```
GET /v1alpha1/meetings/abc123
→ "Authorization is required for request on GET /v1alpha1/meetings/abc123"
```

The record id the caller asked for comes back inside an error about their credential, before the credential was accepted — so the path is reflected to an unauthenticated caller, and anything that logs or forwards these messages carries whatever was in the URL.

**The error's type is a JavaScript class name.** `AuthorizationRequiredError`, under a field called `name` — which is where a JavaScript `Error` keeps its constructor name. The implementation's class hierarchy is the API's error vocabulary.

**A wrong path and a wrong method are an HTML page.**

```html
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>Error</title>
</head>
<body>
<pre>Cannot GET /v1alpha1/cauldron-nope</pre>
</body>
</html>
```

`text/html`, 404, from Express's default error handler, on an API that is JSON everywhere else. A client calling `.json()` on a typo throws rather than reporting it. And `PUT /v1alpha1/meetings` — a real path with the wrong verb — gets `Cannot PUT /v1alpha1/meetings` at 404 rather than a 405, so a wrong method is reported as a missing page.

**The version is alpha, and the documentation says so.** Every path is under `/v1alpha1`, on a key-gated commercial API, and the reference's own note is "Expect upcoming changes." One documented endpoint, `/v1alpha1/meetings/{meetingId}/highlights`, is already marked deprecated — retired before the version it lives in has shipped.

**The API host is called pasta.** `pasta.tldv.io` is the documented base URL and the only host these routes answer on.

**The listing describes its paging in the response and not in the request.** The documented envelope carries `page`, `pages`, `total` and `pageSize`. The parameters that set them are not documented on the reference at all.

Also pinned: one record carries two free-form bags — `extraProperties` ("Additional meeting metadata like conferenceId") and `metadata` ("Custom fields for integration") — with nothing saying which is meant for whom; `/v1alpha1/health` answers `{"status":"ok"}` to anyone with no key, on the same host and under the same alpha prefix as everything else; and `duration` is a number of seconds beside `happenedAt`, a string, with no unit declared for the first and no format for the second.

## Sources

- [tl;dv API reference](https://doc.tldv.io/) — the documented endpoints, the `x-api-key` header, and the listing and single-meeting response shapes.
- Live: `pasta.tldv.io`, struck 2026-09-13 with no key, a wrong key, an unrouted path, a wrong method, a single-meeting path, and `/v1alpha1/health`.

## Modelling limits

- **No description is published.** tl;dv serves no OpenAPI document at any of the usual paths, so there is no `spec` to fingerprint and `cauldron drift` has nothing to compare. The shapes here come from the reference's own examples.
- **The paging parameter names are assumed.** The response names `page` and `pageSize`; the reference does not name the parameters that set them. This Recipe assumes they match and says so here rather than silently.
- **The not-found body is not observed.** Fetching a meeting that is not there needs a real key, and every unauthenticated attempt is answered by the credential check first. The 404 served here is the provider's own envelope with a plausible sentence in it, marked documentation-only.
- **Record ids are 24 characters of opaque text.** The reference gives no example id, so the length is a guess shaped like the ids this kind of service mints; nothing here claims it was observed.
- **Three routes of eight.** The listing, one meeting, and health. Import, download, transcript, notes and the deprecated highlights endpoint are the rest.
- **Nothing is mapped in detection.** tl;dv is reached through a plain HTTP call carrying an `x-api-key` header, which does not resolve to this host through a dependency file. Checked 2026-09-13.
