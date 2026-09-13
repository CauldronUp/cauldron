# cockroach

Emulates the CockroachDB Cloud cluster listing for local development and tests.

**10 conformance cases, 3 checked against the live API on 2026-09-13.**

Read from the OpenAPI document Cockroach Labs publishes in its own Go SDK ([`internal/spec/openapi.json`](https://raw.githubusercontent.com/cockroachdb/cockroach-cloud-sdk-go/master/internal/spec/openapi.json)), and struck live against `cockroachlabs.cloud` on 2026-09-13 with no header, with a wrong bearer, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**The error code is a bare 16.**

```json
{
  "code": 16,
  "message": "missing authorization header; see https://www.cockroachlabs.com/docs/cockroachcloud/cloud-api.html#call-the-api"
}
```

Sixteen is `UNAUTHENTICATED` in `google.rpc.Code`, and nothing in the response says so. A client switching on `code` has an integer from a gRPC enum, in a REST body, beside an HTTP status that already carried the same meaning. And the message carries a URL with a file extension and an anchor, so the sentence a client prints links to a section of a web page.

**A wrong bearer names the two things it would have accepted.** "invalid secret provided in authorization header; expected either an API key or a JWT" — two credential kinds, disclosed by the failure, in an API whose document declares one `securityScheme`, `{"type": "http", "scheme": "bearer"}`, and never mentions either.

**A path that does not exist is answered by the credential check.** `GET /api/v1/cauldron-nope` with no header answers the same 401, so a typo in a path is reported as a missing header.

**Five declared failures, and the document describes the body of none of them.** `CockroachCloud_ListClusters` lists 400, 401, 403, 404 and 500, each with `"schema": {}` — an empty schema, which constrains nothing. Only `default`, the response nobody looks up, points at a real type. And that type, `Status`, is `{code: int32, details: [Any], message}`: the shape the live failures actually have, declared once, referenced by no status a caller will meet.

**Three status fields, and the type names cross over the field names.** A cluster carries `state` (typed `ClusterState.Type`), `operation_status` (typed `ClusterStatus.Type`) and `upgrade_status` (typed `ClusterUpgradeStatus.Type`). So the field called status is typed Status, the field called state is typed State, and the two words mean different things that nothing on the record explains.

**`DELETED` is a state a cluster can be in.** `ClusterState.Type` is `["CREATING", "CREATED", "CREATION_FAILED", "DELETED", "LOCKED"]`, with `deleted_at` beside it — so a deleted cluster is a cluster you can still be handed.

**And `LOCKED` is an instruction to the caller.** Its description:

> An exclusive operation is being performed on this cluster. Other operations should not proceed if they did not set a cluster into the LOCKED state.

A mutual-exclusion protocol, written into an enum member's description, and enforced by whoever reads it.

**The paging parameters have dots in their names.** `pagination.page`, `pagination.limit`, `pagination.as_of_time`, `pagination.sort_order`, `pagination.sort_by` — a nested message flattened into a query string by a gRPC gateway, so the dot is a literal character in the parameter name. The response carries `next_page` and `previous_page` and no total, so an empty page is how a caller learns it has finished.

**And one field names a single cloud vendor.** `azure_cluster_identity_client_id`, thirty-two characters, on the record that also carries `cloud_provider` with `GCP`, `AWS` and `AZURE` in it.

**`info.version` is `2024-09-16`** — a date, where a version belongs. The enum descriptions begin with a space and a hyphen because they are protobuf comments rendered as a bulleted list of members. And a `PUT` answers 411 Length Required with an HTML page whose text reads "POST requests require a `Content-length` header" — a load balancer describing a PUT as a POST.

## Sources

- [`cockroach-cloud-sdk-go/internal/spec/openapi.json`](https://raw.githubusercontent.com/cockroachdb/cockroach-cloud-sdk-go/master/internal/spec/openapi.json) — recorded by `cauldron drift`, `spec_seen: 2026-09-13`.
- Live: `cockroachlabs.cloud`, struck 2026-09-13 with no header, a wrong bearer, an unrouted path, and a wrong method.

## Modelling limits

- **One route.** The cluster listing. Folders, backups, private endpoints, log export, SQL users, databases, invoices and the version-deferral surface are the rest of a 690 KB document.
- **The 411 is not served.** A `PUT` never reaches the API: a load balancer answers an HTML page before it does. This Recipe records that rather than emulating a proxy's error page, so a wrong method here gets Cauldron's own 405.
- **`previous_page` is not sent.** The schema declares it beside `next_page`, and a cursor cannot be arithmetic'd backwards — which is why Cauldron does not invent one. The forward token is served and the backward one is recorded.
- **The success fixture is document-derived.** Listing clusters needs a real key; the records here are `Cluster`'s own required fields with values of the declared shapes, and every case reading them is marked documentation-only.
- **Nothing is mapped in detection.** CockroachDB Cloud is reached through `cockroach-cloud-sdk-go`, the `ccloud` CLI, or a plain HTTP call carrying a bearer token, and none of those resolve to this host through a dependency file. Checked 2026-09-13.
