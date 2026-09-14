# squadcast

Emulates the Squadcast service listing for local development and tests.

**10 conformance cases, 4 checked against the live API on 2026-09-13.**

Read from the OpenAPI document Squadcast publishes for its own reference site, and struck live against `api.squadcast.com` on 2026-09-13 with no header, with a wrong bearer, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**Every service hands out its own ingestion key.** `api_key` is a `required` property of the service record, so `GET /v3/services` returns the alert-intake key for every service in the organisation to anyone who can read the listing. A credential scoped to read configuration ends up holding the keys that write alerts.

**Two refusals, two content types, one shape.**

```
(no header)      401  application/json  {"meta":{"status":401,"error_message":"No credentials sent"}}
Bearer notreal   401  text/plain        {"meta":{"error_message":"invalid access token","status":401}}
```

The same object, the same status, and one of them announced as plain text — with the keys in the other order, which is two code paths rather than one.

**And the status inside is declared as either a string or a number.** `Common.V3.ErrorMeta.status` is `anyOf: [{type: string}, {type: integer}]`, so a client switching on it has to handle `401` and `"401"`. The description of that error type reads "Represents a single response containing data of type T." — a generic envelope's docstring, left on the failure shape.

**Everything is under `meta` or under `data`, and never both.** A success is `{"data": [...]}` and a failure is `{"meta": {...}}`, so the two halves of the envelope never appear together and a client cannot read a status off a successful response.

**The listing declares twelve responses.** 200, 400, 401, 402, 403, 404, 409, 422, 500, 502, 503 and 504 — including Payment Required, described as "Client error", and Gateway Timeout, described as "Server error", on an operation that lists services.

**Four configuration fields, three named after features and one named nothing.** A service `required`s `auto_pause_transient_alerts_config`, `intelligent_alerts_grouping_config` and `delay_notification_config` — and `config`, a bare object beside them.

**And the record names two specific third-party products.** `slack` and `jira_cloud` are typed fields on the service, so the shape of a Squadcast service depends on which integrations Squadcast has built.

**The credential is checked before the path**, so an unrouted path and a wrong method both answer "No credentials sent"; a service carries an `email` address of its own, for alerts that arrive by mail; and the same failure comes from `api.squadcast.com` and `api.eu.squadcast.com` alike — two servers in one document, with no field on any record saying which one answered.

## Sources

- [`Squadcast-API-Spec.json`](https://openapi.gitbook.com/o/rNxiyZGOoRfKnq7j5RX6/spec/Squadcast-API-Spec.json) — the document behind [`apidocs.squadcast.com`](https://apidocs.squadcast.com); recorded by `cauldron drift`, `spec_seen: 2026-09-13`.
- Live: `api.squadcast.com` and `api.eu.squadcast.com`, struck 2026-09-13 with no header, a wrong bearer, an unrouted path, and a wrong method.

## Modelling limits

- **One route of a hundred and thirty-eight.** The service listing. Incidents, escalation policies, squads, schedules, runbooks, global event rules, status pages, audit logs and the analytics surface are the rest.
- **One region.** The document declares a US server and an EU one; this Recipe serves one host, and the finding is that no record says which answered.
- **The success fixture is document-derived.** Listing services needs a real token; the records here are `V3.Services.Service`'s own required fields with values of the declared shapes — including an `api_key`, because the schema requires one — and every case reading them is marked documentation-only.
- **Nothing is mapped in detection.** Squadcast is reached through its Terraform provider or a plain HTTP call carrying a bearer token, and neither resolves to this host through a dependency file. Checked 2026-09-13.
