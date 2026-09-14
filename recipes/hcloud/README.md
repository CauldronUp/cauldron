# hcloud

Emulates the Hetzner Cloud server listing for local development and tests.

**12 conformance cases, 4 checked against the live API on 2026-09-13.**

Read from the OpenAPI document Hetzner publishes at [`docs.hetzner.cloud/cloud.spec.json`](https://docs.hetzner.cloud/cloud.spec.json), and struck live on 2026-09-13 with no token, with a wrong token, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**A required field named `deprecated` is deprecated, and the replacement it names is deprecated too.** On `server_type`:

| field | type | what it says |
|---|---|---|
| `deprecated` | `boolean` | "This field is deprecated. Use the deprecation object instead." |
| `deprecation` | `object \| null`, marked `deprecated: true` | "This field is deprecated. Use the `deprecation` object in the `locations` field instead (`.locations[].deprecation`)." |

The first points at the second, the second points at a third, and the one in the `required` array is the one telling you not to read it.

**The count of pages goes null on the last page.** `last_page` is "Page number of the last page available. Can be null if the current page is the last one." So the field that says how many pages exist is guaranteed present only while you still have pages left — and `total_entries` sits beside it as "Total number of entries that exist for this query. Can be null if unknown."

**No error in this document has a status code.** The listing declares three responses: `200`, `4xx` and `5xx`. Ranges, not codes. A client generated from it cannot tell 401 from 404 from 422, and `cauldron drift` says so of every error this Recipe declares:

```
not backed: no operation the Recipe routes to answers 401, which token_required declares
not backed: no operation the Recipe routes to answers 404, which unknown_route declares
```

The errors are real — they were struck live. The document simply has no status code to back them with.

**And the one example body under `4xx` is a sentence the API does not send.** The document's example is `{"error": {"code": "unauthorized", "message": "unable to authenticate", "details": null}}`. Live, the two ways to fail authentication are:

```
(no token)       401  {"error":{"code":"unauthorized","details":null,"message":"token is required"}}
Bearer notreal   401  {"error":{"code":"unauthorized","details":null,"message":"the token you have provided is invalid"}}
```

Neither says "unable to authenticate".

**The path is routed before the token is read; the method is not.**

```
GET /v1/cauldron-nope   (no token)  404  api route not found
GET /v1/servers         (no token)  401  token is required
PUT /v1/servers         (no token)  401  token is required
```

A mistyped path is answered without a credential. A mistyped verb is not.

**And the two answers come out of different encoders.** The 404 arrives pretty-printed over two-space indentation, `content-type: application/json; charset=UTF-8`, with a `content-length`, keys in the order `message`, `code`, `details`. The 401 arrives compact, `content-type: application/json` with no charset, `transfer-encoding: chunked`, keys in the order `code`, `details`, `message`. One envelope, two serialisers, and `details` — declared `["object", "null"]` and not in `required` — present and null in both.

**Cores are a floating-point number.** `cores`, `memory` and `disk` on `server_type`, and `primary_disk_size` on the server, are all `"type": "number"`. A CPU count that can carry a fraction, and a disk size in gigabytes that can too.

**A 3.4 MB document has eight reusable schemas.** 151 paths, and `components.schemas` holds exactly `ServiceTCP`, `ServiceHTTPProtocol`, `ServiceHTTPSProtocol`, `TargetTypeServer`, `TargetTypeLabelSelector`, `TargetTypeIP`, `ZonePrimary` and `ZoneSecondary`. Every other shape in the API, the server record included, is written out inline at each path that uses it.

Also pinned: `public_net.ipv4` is required and nullable, and its own `required` array is `["ip", "blocked", "dns_ptr"]` — so the Primary IP object may arrive without the id of the Primary IP it is; `backup_window` is a string like `"22-02"`, a time range with no date and no zone anywhere but the prose; and the three traffic counters are bytes "for the current billing period", on a record with no field naming which period that is.

## Sources

- [`docs.hetzner.cloud/cloud.spec.json`](https://docs.hetzner.cloud/cloud.spec.json) — OpenAPI 3.1.2, `Hetzner Cloud API 1.0.0`, served without a credential; recorded by `cauldron drift`, `spec_seen: 2026-09-13`.
- Live: `api.hetzner.cloud`, struck 2026-09-13 with no token, a wrong token, an unrouted path, and a wrong method.

## Modelling limits

- **Two routes of a hundred and fifty-one.** The server listing and one server. Actions, images, volumes, networks, firewalls, load balancers, primary IPs, certificates, ISOs, placement groups, datacenters, locations, pricing and SSH keys are the rest.
- **`last_page` is served as a number.** The document says it can be null on the last page; that claim needs a token and a real fleet to strike, so this Recipe serves the number and records the claim here rather than guessing when it disappears.
- **The pretty-printed 404 is described, not reproduced.** This Recipe serves one JSON encoding for both errors; the difference in whitespace, charset, transfer encoding and key order between the live 401 and the live 404 is a finding in prose.
- **The success fixture is document-derived.** Listing servers needs a real token; the records here are the listing's own inline schema with values of the declared types, and every case reading them is marked documentation-only.
- **A wrong method with a *valid* token was not struck.** With no credential it answers 401, which is what this Recipe serves. What it answers to an authenticated caller is unknown, and no case here claims it.
- **Nothing is mapped in detection.** Hetzner Cloud is reached through `hcloud-go`, `hcloud-python` or a plain bearer call, and none of them resolves to this host through a dependency file. Checked 2026-09-13.
