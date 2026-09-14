# depot

Emulates the Depot project listing for local development and tests.

**9 conformance cases, 5 checked against the live API on 2026-09-13.**

Read from Depot's own published CLI ([`depot/cli`](https://github.com/depot/cli)'s `pkg/cmd/projects/output.go`), and struck live against `api.depot.dev` on 2026-09-13 with no credential, with a wrong token, with a method the path does not take, with no content type, on a path that does not exist, and at the root.

## What this Recipe found

**Four failure envelopes on one host, and `code` means two different things.**

```
POST …/ListProjects, no header   401  {"code":"unauthenticated","message":"Missing authorization header"}
POST …/Nope                      404  {"message":"Route POST:/depot.core.v1.ProjectService/Nope not found","error":"Not Found","statusCode":404}
POST …/ListProjects, no type     415  {"statusCode":415,"code":"FST_ERR_CTP_INVALID_MEDIA_TYPE","error":"Unsupported Media Type","message":"Unsupported Media Type"}
GET  …/ListProjects              405  (no body at all)
```

On the first, `code` is a gRPC status name. On the third, `code` is a Fastify internal error constant. Two vocabularies, one key, one hostname — and the key a client switches on is the one that changes meaning.

**The 415 says the same thing twice and the 404 says it three ways.** `"error":"Unsupported Media Type","message":"Unsupported Media Type"` is the reason phrase in both fields; the routing failure carries `message`, `error` and `statusCode` where the RPC failure carries none of them.

**And a wrong method is a status line and nothing else.** `GET` on a real RPC path answers 405 with an empty body, so the one failure with no shape at all is the one a browser produces by accident.

**The root is open and says `{"ok":true}`.** No credential, 200, two characters of information.

**The vendor's CLI renames every field before printing it.** `output.go` reads `project.GetProjectId()`, `GetOrganizationId()`, `GetRegionId()`, `GetCreatedAt()` and `GetCachePolicy()` off the protobuf message, and then writes them out under `project_id`, `organization_id`, `region_id`, `created_at` and `cache_policy`. A script reading `depot projects get --output json` and one calling the API see different keys for the same fields.

**A cache budget is bytes and a retention is days, on one object.** `cache_policy` is `{keep_bytes, keep_days}`, and the CLI divides the first by `1024 * 1024 * 1024` to print gigabytes — so the units are in neither the field names nor the response.

**The timestamp is RFC 3339 with nanoseconds**, formatted from a protobuf `Timestamp` and marked `omitempty`, so a project with no creation time prints no key at all. And `newProjectOutput` returns the error "API returned no project" when the message is nil — the CLI saying what it expects, rather than what the contract promises.

## Sources

- [`depot/cli`](https://github.com/depot/cli) — `pkg/cmd/projects/output.go`: the accessors it reads and the names it prints.
- Live: `api.depot.dev`, struck 2026-09-13 with no credential, a wrong token, a wrong method, a missing content type, an unrouted path, and the root.

## Modelling limits

- **Two routes.** The project listing and the open root. Builds, build tokens, organisations, tokens and the registry surface are the rest.
- **The 415 is recorded, not served.** It is produced by omitting the request's content type, which a conformance case sets by construction.
- **The wire names are protobuf JSON's default.** Depot serves Connect RPC; the field names the CLI reads are `project_id`, `organization_id`, `region_id`, `created_at` and `cache_policy` in the proto, and this Recipe serves them under protobuf JSON's lowerCamelCase mapping. The finding is the renaming the CLI does on top, which is in its own source.
- **No `spec:`.** Depot's surface is described by protobuf definitions on a schema registry rather than by an OpenAPI document at any address this Recipe could reach, so there is nothing for `cauldron drift` to record.
- **The success fixture is CLI-derived.** Listing projects needs a real token; the records here are the fields `output.go` reads, with values of the shapes it formats, and every case reading them is marked documentation-only.
- **Nothing is mapped in detection.** Depot is reached through the `depot` CLI or a plain HTTP call carrying a bearer token, and neither resolves to this host through a dependency file. Checked 2026-09-13.
