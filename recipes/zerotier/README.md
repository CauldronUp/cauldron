# zerotier

Emulates the ZeroTier Central network listing and status endpoint for local development and tests.

**14 conformance cases, 9 checked against the live API on 2026-09-13.**

Read from the OpenAPI document ZeroTier vendors into its own Go client at [`zerotier/go-ztcentral/spec.json`](https://raw.githubusercontent.com/zerotier/go-ztcentral/main/spec.json), and struck live on 2026-09-13 with no credential, with a wrong credential, with the wrong scheme, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**One endpoint refuses in three content types.**

```
GET /api/v1/network  (no Authorization)     403  text/html         <a href="/login">Forbidden</a>.
GET /api/v1/network  Authorization: Bearer  403  text/plain        forbidden
GET /api/v1/network  Authorization: token   401  application/json  {"type":"internal","message":"Access Denied"}
GET /api/v1/cauldron-nope                   404  text/plain        404 page not found
PUT /api/v1/network                         404  text/plain        404 page not found
```

HTML, plain text and JSON, from one path, at three statuses. The document declares no error body anywhere: its four shared responses — `BadRequest`, `AccessDeniedError`, `NotFound`, `UnauthorizedError` — each carry a description and no `content`.

**The credential-less refusal is a redirect that is not one.** It is a 403, and it carries `location: /login` in the headers with Go's redirect body — an `<a>` element — in the payload. A client following `Location` on a 3xx never sees it, because this is not a 3xx; a client reading the body for an error message gets an anchor tag.

**A wrong verb is reported as a missing page.** `PUT` on a path that exists answers the same `404 page not found` as a path that does not, so a client cannot tell a typo in the URL from a typo in the method. Both are `net/http`'s default `NotFound` handler.

**And a caller's bad token is categorised as `internal`.** `{"type":"internal","message":"Access Denied"}` — the one JSON refusal of the five, with the fault classified as the server's.

**The status endpoint ignores the credential it is documented to describe.** `/status` is summarised "Obtain the overall status of the account tied to the API token in use". It is covered by the document's global `security: [{"tokenAuth": []}]` with no per-operation override. It answers `200` to a request with no `Authorization` header, to a wrong token, and to the wrong scheme — and `user` comes back `null` in all three.

**And it breaks its own schema on the way.** `version` is declared `{"type": "string"}` with no `nullable`, and every live response sends `"version": null`. `apiVersion` carries the example `"4"` and sends `"1"`.

**Five fields of the live response are not in the schema at all.** `cauldron drift` reports them:

```
not backed: GET /status: the success schema does not declare online
not backed: GET /status: the success schema does not declare secondFactor
not backed: GET /status: the success schema does not declare stripePublishableKey
not backed: GET /status: the success schema does not declare supportEmbedCode
not backed: GET /status: the success schema does not declare clusterNode
```

All five arrive on the first unauthenticated request. Among them: a live `pk_live_` Stripe publishable key, a Kubernetes pod name, and a `<script>` element delivered as the value of a JSON string.

**`uptime` is nanoseconds, per replica, and its example is a wall clock.** Four consecutive unauthenticated calls:

```
ztc-central-69848ffd6d-gsjpz  443946516007781
ztc-central-69848ffd6d-dht4v  443915929187436
ztc-central-69848ffd6d-gsjpz  443946774393450
ztc-central-69848ffd6d-fw7d7  444149191473658
```

Three pods, and uptimes 2.3×10¹¹ apart between them — so "Uptime on server" means whichever replica the load balancer picked. The two calls that landed on `gsjpz`, 259 ms apart, differ by 258,385,669, which is 10⁹ per second: nanoseconds, a unit the document does not state. Its example value for the field is `1613067920454` — the same number it gives as the example for `clock`, which is a Unix timestamp in milliseconds.

**A member count is documented to be wrong on the listing.** `onlineMemberCount` carries the note "May be 0 on endpoints returning lists of Networks". The count in a listing is not the count; the only way to learn it is to fetch every network individually. This Recipe serves both: the listing zeroes the field, the single fetch of the same network carries 9.

**Nothing on a network is required.** `Network` has no `required` array, so `id` is optional in the schema that describes it. Neither does `NetworkConfig`.

**The security scheme is an HTTP scheme that does not exist.** `{"type": "http", "scheme": "token"}`. `token` is not in the IANA HTTP Authentication Scheme Registry — and live proves the string is load-bearing: `Authorization: token …` reaches the authenticator and returns 401, `Authorization: Bearer …` does not and returns 403.

**The same document ships under two names.** `zerotier/zerotier-rust-api` holds `zerotier-central-api/openapi.json` and `zerotier-one-api/openapi.json`, byte for byte identical (md5 `03047936668a49ff3e32a581b30c1e7e`), both titled "ZeroTier Legacy Central API", both declaring `servers: [{"url": "https://api.zerotier.com/api/v1"}]`. The crate named for ZeroTier One — the local node service, reached on `127.0.0.1:9993` — is generated from the cloud API's document, so its client is pointed at the internet. And that copy is not the one in the Go client either: 15 paths against 14, a different `Network.permissions`, and the same version string `v1` on both.

Also pinned: the listing declares no parameters at all — no page, no limit, no filter, every network in one bare array; `id` appears on the network and again inside `config`; `private` is documented "If false, members will \*NOT\* need to be authorized to join", with the emphasis in asterisks inside a JSON description; `multicastLimit` warns that "Setting this to 0 will disable IPv4 communication"; and `x-envoy-decorator-operation: ztc-central-service.default.svc.cluster.local:80/*` rides on every refusal.

## Sources

- [`zerotier/go-ztcentral/spec.json`](https://raw.githubusercontent.com/zerotier/go-ztcentral/main/spec.json) — OpenAPI 3.0.0, `ZeroTier Central API v1`, the document ZeroTier's own Go client is generated from; recorded by `cauldron drift`, `spec_seen: 2026-09-13`.
- [`zerotier/zerotier-rust-api`](https://github.com/zerotier/zerotier-rust-api) — the two identical copies under `zerotier-central-api/` and `zerotier-one-api/`.
- Live: `api.zerotier.com`, struck 2026-09-13 with no credential, a wrong token, the wrong scheme, an unrouted path, a wrong method, and four consecutive status calls.

## Modelling limits

- **Three routes of fourteen.** Status, the network listing, and one network. Members, users, tokens, organisations, invitations and `/randomToken` are the rest.
- **The Stripe key in the fixture is not the real one.** Live, `/status` returns a genuine `pk_live_` publishable key without a credential. That is the finding; serving the key itself is not, so the fixture carries a placeholder and the case asserts the prefix.
- **The pod name is one that answered on 2026-09-13.** It is a Kubernetes pod name and will not exist by the time you read this; the finding is that the field is there at all, and that it changes between consecutive requests.
- **The network fixture is document-derived.** Listing networks needs a real token; the records here are `Network` and `NetworkConfig`'s own properties with values of the declared types.
- **`x-envoy-decorator-operation` is served on one refusal, not all five.** Live it rides on every one; this Recipe declares it where it is asserted rather than repeating it into every error.
- **Nothing is mapped in detection.** ZeroTier Central is reached through `go-ztcentral`, the `zerotier-central-api` crate, or a plain `Authorization: token` call, and none of them resolves to this host through a dependency file. Checked 2026-09-13.
