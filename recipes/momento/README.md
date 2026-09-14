# momento

Emulates the Momento HTTP cache read for local development and tests.

**8 conformance cases, 6 checked against the live API on 2026-09-13.**

Read from the HTTP API reference at [`docs.momentohq.com`](https://docs.momentohq.com/cache/develop/api-reference/http-api), and struck live on 2026-09-13 with no token, with a wrong token in the header, with a wrong token in the documented query parameter, with the wrong scheme, with no key, with two keys, without the required TTL, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**The documented query parameter is refused as a header format problem.** Momento documents two carriers for the credential — `Authorization: <token>` or `?token=<token>` — and its own refusal names only one:

```
GET /cache/mycache?key=k&token=notarealtoken
401 {"status":401,"title":"Unauthorized","detail":"Invalid authorization header format"}
```

There was no authorization header on that request. A caller who used the documented query parameter is told the header they did not send is malformed.

**And the message for no token at all names both carriers:**

```
401 {"status":401,"title":"No Auth Token Provided","detail":"You must provide a valid momento auth token via the 'Authorization' header or the 'token' query parameter."}
```

So the API can describe the query parameter when nothing is wrong with it, and forgets it the moment something is.

**The credential is documented to live in a URL.** `?token=<token>` is a first-class option for a cache credential — which puts the secret in every access log, proxy log, browser history and `Referer` along the way. The same host answers `access-control-allow-credentials: true`.

**The `Authorization` header takes a bare token with no scheme.** Not `Bearer`, not `Token`: the header value *is* the secret. Sending `Bearer <token>` draws the same "Invalid authorization header format" that a wrong token does, so the mistake every other API's habits will produce is indistinguishable from a revoked key.

**The error that fires when no key was given says "but not both".** Two neighbouring 400s, struck live:

```
no key at all:  {"title":"No Key Specified",       "detail":"You must specify either `key` or `key_base64, but not both."}
two keys:       {"title":"Multiple Keys Specified","detail":"You must specify either 'key' or 'key_base64', not both."}
```

The first fires when *neither* was sent and tells the caller not to send both. Its backtick is unbalanced and its comma sits inside the parameter name. The second is the same sentence with different quoting and a correct ending.

**A serde error reaches the caller verbatim.** A `PUT` without the required `ttl_seconds`:

```
400 {"status":400,"title":"Unable to Parse Query Parameters","detail":"Failed to deserialize query string: missing field `ttl_seconds`"}
```

The Rust deserialiser's own words, backtick and all, inside an RFC 9457 problem document.

**The 405 is the one failure with no problem document.** Every other refusal is `application/problem+json`. A `PATCH` answers `405`, `allow: GET,HEAD,PUT,DELETE` (no spaces after the commas), `content-length: 0`, and no `Content-Type` at all — so the one status whose semantics are best specified is the one served with nothing in it.

**An unrouted path reflects the path in plain text.** `No route for /cauldron-nope`, `text/plain`, on an API that speaks problem+json everywhere else. A client parsing failures as JSON throws on a typo rather than reporting it.

**Every problem document repeats the status in the body**, and none of them carries a `type` — so every problem is RFC 9457's default `about:blank`, and there are two numbers where there could be one.

Also pinned: the host is `api.cache.cell-us-east-1-1.prod.a.momentohq.com`, publishing the cell, the region, the environment and a single-letter segment in DNS, with the documented base given as the pattern `cell-{cell}-{region}-1`; and the HTTP API's own reference says most applications "should likely use [the] SDK clients" instead — documented and discouraged on the same page.

## Sources

- [Momento HTTP API reference](https://docs.momentohq.com/cache/develop/api-reference/http-api) — the two credential carriers, the `key`/`key_base64` choice, the required `ttl_seconds`, and the status code table.
- Live: `api.cache.cell-us-east-1-1.prod.a.momentohq.com`, struck 2026-09-13 on nine distinct requests.

## Modelling limits

- **No description is published.** Momento serves no OpenAPI document for the HTTP API, so there is no `spec` to fingerprint.
- **Three live errors are recorded here and not served.** "No Key Specified", "Multiple Keys Specified" and "Unable to Parse Query Parameters" are all request-shape failures, and this Recipe's one route answers with bytes rather than a record, so nothing in it can reach them. Their exact bodies are quoted above.
- **The hit's media type was not observed.** A successful read needs a real token and a real cache. This Recipe serves `application/octet-stream`, which is what "a scalar value" of bytes implies, and marks both success cases documentation-only.
- **One route of three.** The read. `PUT` (which needs `ttl_seconds`) and `DELETE` share the path.
- **Nothing is mapped in detection.** The Momento HTTP API is reached through a plain HTTP call carrying a bare `Authorization` value or a `token` query parameter, and neither resolves to this host through a dependency file. Checked 2026-09-13.
