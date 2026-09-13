# ngrok

Emulates the ngrok endpoint listing for local development and tests.

**11 conformance cases, 4 checked against the live API on 2026-09-13.**

Read from ngrok's own published Go client ([`ngrok-api-go/datatypes.go`](https://github.com/ngrok/ngrok-api-go)), and struck live against `api.ngrok.com` on 2026-09-13 with no header, with a wrong bearer, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**The refusal quotes your credential back in full.** Sending `Authorization: Bearer notarealtoken`:

```json
{"error_code":"ERR_NGROK_202","status_code":403,"msg":"The API authentication you specified does not look like a valid credential. Your credential: 'notarealtoken'. API keys and instructions are available on your dashboard: https://dashboard.ngrok.com/api-keys","details":{"operation_id":"op_…"}}
```

Not the first few characters, not a mask: the whole value, in single quotes, in the message field a client prints and a log keeps. Send a real key to the wrong account by mistake and the response contains it.

**The error code has an HTTP success in its name.** `ERR_NGROK_200`, on a 403, for a missing credential; `ERR_NGROK_202` for a wrong one. The digits are a vendor catalogue number and they read as status codes.

**And the status is 403 for having presented nothing.** Forbidden, to a caller the server has never heard from, with no `WWW-Authenticate` anywhere.

**The field a client switches on is missing from some failures.** The 403s carry `error_code`; the 404 does not:

```json
{"status_code":404,"msg":"Not Found","details":{"path":"/cauldron-nope"}}
```

So `error_code` is the discriminator, and it is absent exactly when the request went somewhere unexpected. `details` changes its keys to match — `{"operation_id": "op_…"}` on the 403s and `{"path": "/…"}` on the 404, one object with two disjoint shapes and nothing naming which to expect.

**A wrong method is reported as a path that is not there.** `PUT /endpoints` answers `{"status_code":404,"msg":"Not Found","details":{"path":"/endpoints"}}` — the path echoed back inside the failure that denies it.

**One address, five fields, two of them deprecated.** From the published type:

```
PublicURL   "deprecated [replaced by URL]: URL of the hostport served by this endpoint"
Hostport    "hostport served by this endpoint (hostname:port) -> soon to be deprecated"
Host, Port  the same thing, split
URL         "the url of the endpoint"
```

A deprecated field, a field announced as "soon to be deprecated" in a struct comment with no date on it, the split form of that one, and the replacement — all on the same record, all populated.

**`proto` is documented and `scheme` is not.** `Proto` is "protocol served by this endpoint. one of http, https, tcp, or tls"; `Scheme` beside it carries no comment at all and holds the same kind of value.

**Every field is `omitzero`.** The vendor's own struct tags mean a port of 0, an empty description and an empty metadata string vanish from the wire, so a client cannot tell "not set" from "set to nothing".

**And `metadata` and `traffic_policy` are strings.** "user-supplied metadata" typed `string`, and a policy document typed `string` beside it, on a record whose other references are `{id, uri}` objects.

## Sources

- [`ngrok/ngrok-api-go`](https://github.com/ngrok/ngrok-api-go) — `datatypes.go`, the generated wire types, and `endpoints/client.go`.
- Live: `api.ngrok.com`, struck 2026-09-13 with no header, a wrong bearer, an unrouted path, and a wrong method.

## Modelling limits

- **One route.** The endpoint listing. Reserved domains, reserved addresses, edges, tunnels, API keys, credentials, IP policies, event subscriptions and the whole edge-module surface are the rest.
- **`operation_id` is a constant.** Live it is a fresh `op_` identifier per request; the live case asserts its shape with a regex rather than its value.
- **`uri` is a constant.** The type documents it as "URI of the endpoints list API resource", which is this path; the forward cursor is served from the request.
- **No `spec:`.** ngrok publishes a generated Go client rather than the document it was generated from, and no OpenAPI description is served at any address this Recipe could find — so there is nothing for `cauldron drift` to record.
- **The success fixture is client-derived.** Listing endpoints needs a real key; the records here are `Endpoint`'s own fields with values of the declared shapes, and every case reading them is marked documentation-only.
- **Nothing is mapped in detection.** ngrok is reached through `ngrok-api-go`, `ngrok-api` on PyPI, or a plain HTTP call carrying a bearer token, and none of those resolve to this host through a dependency file. Checked 2026-09-13.
