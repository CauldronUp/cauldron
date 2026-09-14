# speechmatics

Emulates the Speechmatics batch job listing for local development and tests.

**11 conformance cases, 6 checked against the live API on 2026-09-13.**

Read from the Speechmatics documentation and from `speechmatics-python`'s own batch client, and struck live on 2026-09-13 with no key, with a wrong bearer token, on a path that does not exist, with a method the path does not take, on the previous major version, and on the root.

## What this Recipe found

**The refusal never reaches the application.** Every 401 is nginx's default error page:

```html
<html>
<head><title>401 Authorization Required</title></head>
<body>
<center><h1>401 Authorization Required</h1></center>
<hr><center>nginx</center>
</body>
</html>
```

`Content-Type: text/html`, from a JSON API. The server software is named in the body. There is no `WWW-Authenticate` header, which RFC 9110 says a 401 MUST carry. And "Authorization Required" is nginx's wording for the status the specification calls *Unauthorized* — so the phrase a developer searches for is not the phrase in any standard. A client calling `.json()` on a wrong key throws instead of reporting it.

**The 404 is JSON, and it is served as HTML.**

```
GET /v2/cauldron-nope
404  Content-Type: text/html

{"code":404, "details": "path not found"}
```

Valid JSON, declared `text/html`. A client that trusts the content type will not parse it; one that trusts the body will. The spacing is hand-written: a space after the comma and after the second colon, none after the first.

**So one API answers two failures in two formats, and the header is wrong for one of them.** The 401 is HTML declared HTML. The 404 is JSON declared HTML. The success is JSON declared JSON.

**The previous major version is a missing path.** `/v1/jobs` answers the same `path not found` a route that never existed does — no `410 Gone`, no upgrade sentence, nothing naming v2.

**A wrong verb is answered by the credential gate.** `PUT /v2/jobs` gets the nginx 401, while an unrouted path gets the JSON 404 with no key sent — so the path is resolved before the credential and the method is not.

**The listing takes no parameters, and silently truncates.** `speechmatics-python`'s `list_jobs()` takes no arguments, and its docstring reads:

> "Lists last 100 jobs within 7 days associated with auth_token for the SaaS or all of the jobs for the batch appliance."

The same call returns a capped list from the cloud and an uncapped one from the on-premise appliance. Nothing in the response says which happened, there is no parameter to widen it, and there is no field saying anything was left out — so code written against an appliance and moved to the SaaS starts missing jobs with no error and no signal.

**A single job and a list of jobs arrive under different keys.** The listing is `{"jobs": [...]}`; one job is `{"job": {...}}`. The wrapper changes with the cardinality, so code cannot unwrap both the same way.

Also pinned: `duration` is a bare number, with seconds stated nowhere on the record; `id` is ten characters of lowercase alphanumeric with no prefix and no declared format; the root of the API host answers `200` with a styled HTML page titled "Speechmatics API v2"; and the SDK's own note about concurrency beside the cap warns of "rate-limitting at the API end for over-use".

## Sources

- [Speechmatics batch documentation](https://docs.speechmatics.com/speech-to-text/batch) — the job record and its fields.
- [`speechmatics/speechmatics-python`](https://github.com/speechmatics/speechmatics-python) `speechmatics/batch_client.py` — `list_jobs()` and its docstring, and the `["jobs"]` / `["job"]` unwrapping.
- Live: `asr.api.speechmatics.com`, struck 2026-09-13 with no key, a wrong bearer token, an unrouted path, a wrong method, `/v1/jobs`, and `/`.

## Modelling limits

- **No description is published.** Speechmatics serves no OpenAPI document at any of the usual paths, so there is no `spec` to fingerprint.
- **The success side is document-derived.** Listing jobs needs a real key; the records here are the documentation's own example job plus a second of the same shape, and every case reading them is marked documentation-only.
- **The not-found body for a job is not observed.** Reaching it needs a valid key. This Recipe serves the provider's own `{"code":404, "details": "…"}` shape, declared `text/html` as the real one is, with a sentence about a job rather than a path.
- **The truncation is described, not served.** The 100-job, 7-day cap is in the SDK's docstring, not in any response; this Recipe serves what it holds and records the cap here.
- **The root page is one line of a real one.** Live it is a small styled page; this Recipe serves a minimal document carrying the same title, status and content type.
- **Three routes of many.** The listing, one job, and the root. Submit, delete, transcript, alignment and the notification surface are the rest.
- **Nothing is mapped in detection.** Speechmatics is reached through `speechmatics-python` or a plain bearer call, and neither resolves to this host through a dependency file. Checked 2026-09-13.
