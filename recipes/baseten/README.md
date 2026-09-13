# baseten

Emulates the Baseten management API's model listing for local development and tests.

**7 conformance cases, 3 checked against the live API on 2026-09-13.**

Read from the OpenAPI document Baseten serves at [`api.baseten.co/v1/spec`](https://api.baseten.co/v1/spec), and struck live on 2026-09-13 with no credential, with a wrong one, and against the address its own documentation gives for that document.

## What this Recipe found

**A sixty-three byte error carries a two-and-a-half kilobyte browser Content-Security-Policy.** The body is:

```json
{"code": "PERMISSION_DENIED", "message": "Authorization error"}
```

The headers arriving with it allow scripts from Google Tag Manager, HubSpot, Stripe, Segment, PostHog, New Relic and doubleclick; frames from YouTube, Vimeo and Wistia; and `'unsafe-eval'` and `'unsafe-inline'` besides. None of it can matter — nothing here renders, and a command-line client will never execute a script. Forty times the payload, spent on a browser that is not there.

**And the directives come back in a different order every time.** Three consecutive requests to the same path:

```
frame-src   worker-src       connect-src  form-action  default-src  style-src   ...
connect-src worker-src       default-src  script-src   form-action  child-src   ...
form-action frame-ancestors  style-src    font-src     connect-src  script-src  ...
```

Twelve directives, shuffled per response — a hash iteration order escaping into a security header. The policy is the same policy and the bytes are never the same bytes, so nothing downstream can cache it, diff it, or sign it.

**There is no way to be unauthenticated here.** No credential and a wrong credential both answer 403 with a byte-identical body. Forbidden means the request was understood and refused; a caller who has presented nothing at all receives the same sentence as one whose key was revoked, and `WWW-Authenticate` is on neither.

**The code is from another protocol.** `PERMISSION_DENIED` is `google.rpc.Code` 7 — a gRPC status name, in a JSON body, beside an HTTP status that already says the same thing.

**`Vary` names a cookie and an origin, and not the header the answer turns on.** `Vary: Origin, Cookie`, on an API whose every response depends on `Authorization`.

**The document's own address answers 202 Accepted with nothing in it.**

```
HTTP/1.1 202 Accepted
Content-Length: 0
Content-Type: text/html; charset=UTF-8
x-amzn-waf-action: challenge
Access-Control-Expose-Headers: x-amzn-waf-action
```

A firewall block wearing the status code for "your request has been accepted for processing", with the block deliberately exposed to browser script and nothing at all in the body.

**The listing declares one response and it is the 200.** `Gets all models` lists no 403, no 401 and no 429 — on an API that answered 403 to every request made here and publishes a page about handling 429. `cauldron drift` reports the 403, the 202 and the `/openapi.json` route as unbacked, and each report is right.

**A required field holds a sentence where a machine wants numbers.** `instance_type_name` is "Name of the instance type", required, and reads `1x2 - 1 vCPU, 2 GiB RAM`. Anything that wants the core count parses English.

**Two required identifiers are both nullable.** `production_deployment_id` and `development_deployment_id` are each `anyOf: [string, null]` and each `required`, so a model with neither deployment is a legal model — and the generated Python sample in the document sends `params={'name': None}`, a null query parameter, to fetch the listing.

**And the security scheme admits a second prefix in prose.** `BearerAuth` is `type: http, scheme: bearer`; its description adds "The legacy `Authorization: Api-Key <api_key>` scheme is also accepted". A generated client can only produce the first.

## Sources

- [`api.baseten.co/v1/spec`](https://api.baseten.co/v1/spec) — recorded by `cauldron drift`, `spec_seen: 2026-09-13`.
- Live: `api.baseten.co`, struck 2026-09-13 with no credential, a wrong credential, and against `/openapi.json`.

## Modelling limits

- **Two routes.** The model listing and the firewalled `/openapi.json`. Deployments, environments, secrets, teams, API keys, the billing and usage surfaces and the inference hosts are the rest of a 1.3 MB document.
- **The policy header is served in one fixed order.** Live it is shuffled per response; a deterministic sandbox serves one capture of it, and the cases assert a directive by regex rather than the whole string.
- **The 202 is that one path.** `/v1/models` is not challenged; `/openapi.json` was, twice, on separate probes. This Recipe pins the path that answered rather than claiming a rule about which paths the firewall picks.
- **The success fixture is document-derived.** Listing models needs a real key; the record here is `ModelV1`'s own required field set with the instance-type string the documentation prints, and every case reading it is marked documentation-only.
- **Nothing is mapped in detection.** Baseten is reached through `truss` on PyPI or a plain HTTP call carrying an `Authorization` header, and neither resolves to this host through a dependency file. Checked 2026-09-13.
