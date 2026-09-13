# fillout

Emulates the Fillout forms API for local development and tests.

**5 conformance cases, 2 checked against the live API on 2026-09-13.**

Read from the OpenAPI document Fillout embeds in the documentation page for the endpoint, and struck live against `api.fillout.com` on 2026-09-13 with no credential and with three differently-wrong ones.

## What this Recipe found

**The failures walk a caller through the credential's format.** All three are 400:

```
(no header)          {"statusCode":400,"error":"Bad Request","message":"API authorization header missing"}
Bearer notarealkey   {"statusCode":400,"error":"Bad Request","message":"API key missing underscore"}
Bearer abc_def       {"statusCode":400,"error":"Bad Request","message":"Incorrectly formatted API key was passed"}
```

Three sentences, each narrowing what a valid key looks like. The second says a key contains an underscore. The third says an underscore is not enough. A caller who has never seen a real key learns its shape by sending wrong ones — which is the opposite of what a credential failure is for.

**And none of them is a 401.** Every credential state — absent, malformed, wrong — is `400 Bad Request`. A client branching on 401 never fires, and a retry-on-401 interceptor never runs.

**The document declares no failure at all.** `GET /forms` lists exactly one response, `200`, under a `security: [bearerAuth]` that applies to every operation. The three answers above are the entire failure surface, and the description mentions none of them.

**A form is two fields.** `FormSummary` is `name` and `formId`, both required, and nothing else — no created date, no status, no URL, no submission count. Nothing in a record can sort, filter or date the listing, and the operation takes no parameters with which to ask.

**`formId` is documented as "the public identifier of the form".** Which says there is another one. The listing returns the public id and never the other, and nothing in the document says what the difference is or where the other lives.

**The base URL is a data-residency choice.** Two servers are declared — `api.fillout.com` and `eu-api.fillout.com` — with nothing in a key or a response saying which one holds a given account. The host is configuration a caller has to be told separately from the credential.

**And the envelope is Fastify's, for the third time in this collection.** `{statusCode, error, message}` also arrives from [instantly](../instantly), where the key order reverses between two layers, and from [valtown](../valtown), where the application's own `error` carries a whole sentence and Fastify's carries a reason phrase. Same three keys, three different meanings of `error` across the three.

## Sources

- Live: `api.fillout.com`, struck 2026-09-13.
- [Get forms](https://fillout.com/help/api-reference/get-forms) — the page, and the OpenAPI document embedded in it.
- [Fillout REST API](https://fillout.com/help/fillout-rest-api) — the credential and where it comes from.

## Modelling limits

- **One route.** Listing forms. Form metadata, submissions, single submissions, deletions and webhooks each want their own evidence.
- **A wrong key with an underscore in it gets a third sentence live.** `Incorrectly formatted API key was passed`, struck live with `Bearer abc_def`. Cauldron declares one failure for a rejected credential and cannot branch on the shape of what was sent, so this Recipe serves the underscore message — the one that discloses the format — and the third is quoted above.
- **The EU host is not modelled.** `eu-api.fillout.com` is declared beside the US one and differs only in host.
- **No `spec:`.** Fillout embeds a complete OpenAPI document per endpoint inside that endpoint's page, the same arrangement [mailtrap](../mailtrap) and [lemlist](../lemlist) use. There is nothing single to fingerprint.
- **Nothing is mapped in detection.** Fillout publishes no first-party client library, and a project using it holds an embed script or a webhook handler rather than a client of this API. Checked 2026-09-13.
