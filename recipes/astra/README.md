# astra

Emulates the DataStax Astra DevOps database listing for local development and tests.

**10 conformance cases, 4 checked against the live API on 2026-09-13.**

Read from DataStax's own Java SDK ([`astra-sdk-devops`](https://github.com/datastax/astra-sdk-java)), and struck live against `api.astra.datastax.com` on 2026-09-13 with no header, with a malformed token, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**A format string that was never formatted.** Sending `Authorization: Bearer AstraCS:notarealtoken`:

```json
{"errors":[{"ID":340016,"message":"Resource, token, was malformed. Error: auth-manager returned 400 for call %s/v2/admin/tokens?%s"}]}
```

Two `%s` verbs, unexpanded, in a message delivered to a caller. The arguments that would have filled them are gone, the internal service that failed is named, and its own HTTP status is quoted — so what arrives is a log line from inside the platform with the interesting parts missing.

**And the key is `ID`, in capitals, beside a lowercase `message`.** `{"ID": 340002, "message": "no bearer token in request"}` — a Go struct field exported without a tag, so the initialism's casing became the wire name, in an object whose other key is lower case.

**An unrouted path is Go's own four words.** `GET /v2/cauldron-nope` answers `404 page not found` as `text/plain`, and so does `PUT /v2/databases` — the `net/http` default mux, on a JSON API, giving one answer for a path that is not there and a method that is not allowed.

**The record says what it is and what it was observed to be.** `status` and `observedStatus` are both `DatabaseStatusType`, on the same object, and nothing on the record says which one to believe. In the fixture here the second database is `PARKED` and observed `HIBERNATING` — two words for a stopped database, disagreeing on one record.

**That type has nineteen members and one of them is `UNKNOWN`.** `ACTIVE, ERROR, DECOMMISSIONING, DEGRADED, HIBERNATED, HIBERNATING, INITIALIZING, MAINTENANCE, PARKED, PARKING, PENDING, PREPARED, PREPARING, RESIZING, RESUMING, TERMINATED, TERMINATING, UNKNOWN, UNPARKING`. Four transitions are spelled as an `-ing`/`-ed` pair and four are not: there is a `RESUMING` and no `RESUMED`, a `RESIZING` and no `RESIZED`, a `DECOMMISSIONING` and no `DECOMMISSIONED`.

**Three fields for one concept.** A database carries `keyspace`, `keyspaces` and `additionalKeyspaces` — a singular, a set, and a second set for the ones that are extra. A client asking which keyspaces a database has picks one of three answers.

**And five addresses, two of them web pages.** `studioUrl`, `grafanaUrl`, `cqlshUrl`, `graphqlUrl` and `dataEndpointUrl` sit on the database record, so a JSON listing carries links to a Grafana dashboard and a browser IDE beside the endpoints a program would call.

## Sources

- [`datastax/astra-sdk-java`](https://github.com/datastax/astra-sdk-java) — `astra-sdk-devops`'s `Database`, `DatabaseInfo` and `DatabaseStatusType`.
- Live: `api.astra.datastax.com`, struck 2026-09-13 with no header, a malformed token, an unrouted path, and a wrong method.

## Modelling limits

- **One route.** The database listing. Keyspaces, datacentres, access lists, private endpoints, streaming, roles and tokens are the rest of the DevOps surface.
- **No `spec:`.** DataStax publishes this API as a reference page and a generated Java SDK; `api.astra.datastax.com/openapi.json` answers 403 from CloudFront, and no machine-readable description is served at any address this Recipe could reach — so there is nothing for `cauldron drift` to record.
- **`cost` and `metrics` are not modelled.** The SDK carries both on the record; nothing in a published source fixes their field names closely enough to serve them.
- **The success fixture is SDK-derived.** Listing databases needs a real token; the records here are the SDK's own fields with values of the declared shapes, and every case reading them is marked documentation-only.
- **Nothing is mapped in detection.** Astra is reached through `astrapy` on PyPI, `astra-db-ts` on npm, or a plain HTTP call carrying an `AstraCS:` bearer, and none of those resolve to this host through a dependency file. Checked 2026-09-13.
