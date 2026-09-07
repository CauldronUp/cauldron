# confluent

Emulates the Confluent Cloud API for local development and tests.

**11 conformance cases, 5 checked against the live API on 2026-09-07.**

Written against Confluent's API reference at `docs.confluent.io` and struck live against `api.confluent.cloud` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**The correlation id is on every failure except the one you would report.**

Struck live:

```
GET /iam/v2/service-accounts    (no credential)
401 {"errors":[{"id":"669d33df…","status":"401","detail":"Unauthorized","source":{}}]}

GET /iam/v2/service-accounts    notakey:notasecret
401 {"errors":[{"id":"aed98f10…","status":"401","detail":"invalid API key: …"}]}

GET /iam/v2/cauldron-nope       notakey:notasecret
404 {"errors":[{"status":"404","code":"route_not_found","detail":"Route not found","source":{}}]}
```

The two credential failures carry a fresh 32-hex `id`. The routing failure does not — and it gains a `code` the other two lack. **No two of these responses have the same key set**, three fields are effectively optional, and the one field a support ticket needs is missing from the failure most likely to produce one.

**`source` is always an empty object.** JSON:API defines `source` as a pointer to the part of the request that caused the problem — a JSON Pointer, a parameter name, a header name. Confluent sends `{}` on every failure a public probe can reach, so the field that exists to say *where* says nothing and still costs a key.

**There is no `title`.** JSON:API defines `title` as the short summary and `detail` as the specific explanation; only the long one arrives, so a client rendering failures in a list has to derive the short form from the sentence.

**The wrong-key message is a lesson about key types.** *"invalid API key: make sure you're using a Cloud or Global API Key, and not a Cluster API Key"* — which is genuinely the mistake people make, because Confluent issues three kinds of key that look identical and only one opens this surface. Serving that sentence verbatim is the useful half of this Recipe: a fake answering "Unauthorized" to a Cluster key would cost somebody an afternoon.

**Records carry their own kind and schema version.** `kind: "ServiceAccount"` and `api_version: "iam/v2"` on every record — Kubernetes' shape, so a client can dispatch on the record rather than on the route it came from.

**The cursor is nested under `metadata`.** Not `links.next`, which is where a client written against JSON:API looks — so it finds nothing and reads every page as the last one.

## Modelling limits

- **One route.** Service accounts. Clusters, environments, topics, connectors, schemas, RBAC bindings and the whole metrics surface each want their own evidence.
- **Nothing is mapped in detection.** Confluent's published clients are Kafka clients rather than clients of this control-plane API, and nothing on npm, Packagist or the Go module proxy calls `api.confluent.cloud` under an obvious name.
- **No `spec:`.** Confluent publishes per-service OpenAPI documents behind its documentation site rather than at a single stable URL that `cauldron drift` could fingerprint.
