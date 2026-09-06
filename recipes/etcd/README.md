# etcd

Emulates the etcd v3 HTTP API (the gRPC gateway), for local development and tests.

**9 conformance cases, none checked against a live API.**

Written against etcd's own generated Swagger 2.0 document, published in its own repository — `etcd-io/etcd`, `Documentation/dev-guide/apispec/swagger/rpc.swagger.json`, 44 paths — read on 2026-09-06. etcd runs on the operator's own machines and there is no public instance, so every claim here is the provider's own document.

## What this Recipe found

**Keys and values are base64, and the fields are still called `key` and `value`.** `mvccpbKeyValue.key` is `{"type": "string", "format": "byte"}`, and in OpenAPI `format: byte` means base64. So the field named `key` does not hold the key — it holds the base64 of the key. Same for `value`, and for the `key` and `range_end` a caller *sends*.

`kv.key === "config/timeout"` never matches; `Buffer.from(kv.key, "base64")` does. And a client that sends a plain key gets a successful, empty range — a lookup for a key spelled in base64 that nobody has. A 200 with nothing in it, not an error.

This is the etcd v3 HTTP API's oldest trap, and the document states it in two words most readers slide past.

**Every 64-bit integer is a string:**

```
create_revision  mod_revision  version  lease  count
revision  cluster_id  member_id  raft_term      (on the header)
```

all `{"type": "string", "format": "int64"}` or `uint64`. This is protobuf's own JSON mapping and it exists for a good reason — a 64-bit integer does not survive JavaScript's 53-bit floats. The cost is that every number in the API is a string. `kv.version + 1` is `"41"`. `response.count === 0` is never true, because the value is `"0"`. Comparisons with `<` and `>` coerce and happen to work; equality does not, and neither does arithmetic.

The `header.revision` is the number a caller needs for a compare-and-swap — so the one value that must be compared exactly is the one that cannot be compared as a number.

**Nothing in the document is `required`, and protobuf omits zero values.** Not one schema declares a `required` list. Combined with protobuf's JSON mapping — which drops fields holding their type's zero value — an empty range is expected to answer the header and nothing else: no `kvs` (an empty array is a zero value), no `count` (`"0"` is), no `more` (`false` is). *Marked as reasoning from the encoding rather than an observation:* what is checkable in the document is that every one of those fields is optional.

**Every endpoint is a POST, including every read.** 44 paths, 44 POSTs. `/v3/kv/range` is a POST, as are `/v3/auth/role/list` and `/v3/maintenance/status`. Nothing is cacheable, nothing is idempotent by method, and a caller cannot tell a read from a write by the request line.

This is the second all-POST API in this collection after [outline](../outline), and they arrive at it from opposite directions: Outline chose RPC method names, etcd inherited gRPC.

**The credential goes in `Authorization` with no scheme** — an `apiKey` in the `Authorization` header with no prefix, so the header carries a bare token where RFC 7235 wants `<scheme> <credentials>`. Any proxy or library that parses `Authorization` into a scheme and a value sees a scheme with no credentials.

**A failure carries a gRPC code, not the HTTP status** — 16 is `UNAUTHENTICATED` beside a 401, 7 is `PERMISSION_DENIED` beside a 403. The same pairing as [argocd](../argocd) in this collection, from the same toolchain.

**The document ships the generator's placeholder version,** `"version not set"`, and its title is a filename: `api/etcdserverpb/rpc.proto`. Argo CD's document has the identical version string.

## Modelling limits

- **Nothing here is verified against a live API.** etcd runs on the operator's own machines.
- **The range read only.** 44 paths is a coordination store: put, delete, txn, compact, watch, lease, the auth surface and the maintenance surface each want their own evidence — and watch is a bidirectional stream this format cannot describe at all.
- **Base64 is served, not decoded.** A fixture holding `Y29uZmlnL3RpbWVvdXQ=` is holding what etcd sends, and a client that reads it as a key is making the mistake this Recipe exists to surface.
- **The zero-value omission is not served.** Reproducing protobuf's omission rules would mean modelling the encoder rather than the API; the finding is recorded in words.
