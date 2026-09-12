# kestra

Emulates the Kestra open-source Flows API for local development and tests.

**10 conformance cases, none checked against a live API.**

Read from the OpenAPI document Kestra publishes at [`kestra.io/kestra.yml`](https://kestra.io/kestra.yml) — 462KB describing 194 paths — and contrasted against Kestra Cloud, struck live at `api.kestra.io` on 2026-09-11 with no credential and with a deliberately invalid one.

## What this Recipe found

**The document never says where the API is or how to authenticate to it.** 194 paths, and:

- no `servers` block
- no `security` block
- no `securitySchemes` under `components`

Meanwhile basic authentication has been mandatory on open-source instances since 0.24.0 — a username that is an email address and a password, set in configuration or through a setup page. And every path in the document is templated on `{tenant}` with nothing saying what to put there; the one place in the whole file that fills it in is a `ProblemDetail` example's `instance`: `/api/v1/main/flows`.

So a client generated from this document has no base URL, no credential, and a required path parameter it has to be told about out of band.

**The failure type it does define cannot be relied on for anything.** `ProblemDetail` is a careful RFC 9457 document — `type` described as "what clients branch on", `title` "stable for a given type and never parameterised", plus `status`, `detail`, `instance`, `errors` and `traceId`, each with a real description and example.

It declares **no `required` array at all**. Every member is optional, so `{}` is a conforming problem document, and the field a client is explicitly told to branch on is one the schema permits to be absent.

**And the vendor's own hosted API does not send that shape.** Struck live against Kestra Cloud, a different product on a different host:

```
api.kestra.io/v1/flows   401 application/json   {"message":"Unauthorized"}
```

Not `application/problem+json`. No `type`, no `title`, no `status` — one key the published contract does not contain, from the same company, about the same condition.

**The Cloud host also redirects API paths into a Google login.**

```
GET api.kestra.io/api/v1/main/flows/search   307 Location: /oauth/login/google
```

A programmatic caller using the open-source path shape gets a temporary redirect into an interactive sign-in flow. Every HTTP client follows 307 by default, so what actually arrives is a request for an HTML page, and what the caller finally sees is a 200 with a login screen in it.

**Three schema defects in the Flow type:**

| declaration | what it means |
| --- | --- |
| `variables: {type: object, additionalProperties: false}` | an object permitting no properties, on the field whose whole purpose is arbitrary user-defined keys |
| `tasks: {type: array, additionalProperties: true, minItems: 1}` | `additionalProperties` applies to objects; on an array it means nothing, and it sits beside a `minItems` that does |
| `labels: {oneOf: [Map]}` | the description says "a list of Label (key/value pairs) **or** a map of string to string"; the schema has one branch |

A validator enforcing the first one rejects every flow that uses variables at all.

**And the required search parameter has no wire format.** `GET /flows/search` marks `filters` required and types it as an array of `QueryFilter` *objects* in the query string, with no `style` and no `explode`. The one parameter a caller must send is the one the document does not say how to serialise.

## Sources

- [`kestra.yml`](https://kestra.io/kestra.yml) — the OpenAPI document, version 2.0.1.
- [Basic Authentication Now Required in Kestra OSS 0.24.0](https://kestra.io/docs/migration-guide/v0.24.0/basic-authentication) — the credential the document omits.
- Live, for contrast only: `api.kestra.io`, struck 2026-09-11.

## Modelling limits

- **One route.** Searching flows. Executions, logs, triggers, namespaces, blueprints, KV store, plugins, storage and the rest are 193 more paths.
- **Nothing here is struck live.** Kestra's open-source server is self-hosted, so there is no public instance of this API to strike. The live transcripts above are Kestra Cloud, which is a different API on a different host with a different error shape — quoted as the contrast, not served here.
- **The credential comes from the migration guide, not the description.** The OpenAPI document declares no security scheme at all, so the basic-auth modelling here rests on Kestra's own documentation of the requirement rather than on the contract.
- **An unrouted path is not modelled.** With no instance to strike and no 404 declared on this operation, there is no evidence for what the open-source server answers.
- **The tenant is `main`.** The document templates it and never says what a caller should send; `main` is what an open-source install starts with, and it is the value in the document's own example.
- **Nothing is mapped in detection.** Kestra ships SDKs for Java, Python, Go and JavaScript, and each is pointed at whichever self-hosted instance a deployment runs — so the dependency says the project talks to Kestra and never which host. Checked 2026-09-11.
