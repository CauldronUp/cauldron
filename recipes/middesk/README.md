# middesk

Emulates the Middesk business-verification API for local development and tests.

**11 conformance cases, 4 checked against the live API on 2026-09-09.**

Written against Middesk's published OpenAPI document at `docs.middesk.com` and struck live against `api.middesk.com` on 2026-09-09 with no credential, with a deliberately invalid one, with a credential carrying no scheme prefix, and on a path that does not exist.

## What this Recipe found

**The document declares one response, and it is the success.** `GET /v1/businesses` lists exactly `200`. No 401, no 403, no 404, no 429 — nothing.

The live API answers 401 to every request arriving without a credential, so **the one response every new integration meets first is the one the specification does not mention**. A client generated from this document has no failure type at all: every method returns a business list, and the 401 arrives as a parse error or as a list-shaped object with no `data` in it, depending on which generator you used.

It is the far end of a spectrum this collection now covers completely:

| Provider | Failure statuses declared | Body described |
| --- | --- | --- |
| [attentive](../attentive) | five | none of them |
| [close](../close) | three | typed `"Any type"` |
| middesk | none | — |

All three answer a real, specific JSON body on the wire. In none of the three cases is the document silent because the API is.

**The security scheme is defined and never applied** — the same malpractice [close](../close) has, found on a different provider on the same day. `components.securitySchemes` defines `bearer_auth`. Top-level `security` is null. No operation declares one.

What every operation declares instead is a header parameter literally named `Authorization`, `required: true`, typed string. Two providers, two documents, one mistake: the machinery OpenAPI has for describing authentication sits populated and unreferenced beside a hand-rolled string, so a generator emits a client that takes an `authorization` argument and does no auth handling.

**`data` is not required and `object` is.** The envelope's `required` list is `["object", "has_more"]`.

So a conforming listing must say that it is a list, and whether there is more of it, and need not contain any records. The wrapper is guaranteed; the contents are optional. `total_count` is declared and not required either.

**The failure is an array of one.**

```json
401 {"errors":[{"message":"Unauthorized"}]}
```

An array shape is a promise that several things can go wrong at once and that a client should loop. Here the array has one entry, the entry has one field, and the field is the reason phrase. There is no code to branch on, no field name to blame, no identifier to quote.

And the same body answers three distinct mistakes: a missing credential, a wrong one, and a credential sent with no scheme in front of it. Whatever a caller got wrong, it is told the same thing.

**An unknown path declares HTML and sends nothing.** `404` with `Content-Type: text/html; charset=UTF-8` and a body of zero bytes — a header naming a media type for a document that is not there. [attentive](../attentive) sends an empty 401 with no `Content-Type` at all; this is the same emptiness with the header still attached.

**And the front door is a 204.** `GET /` answers `204 No Content` to a request with no credential whatsoever.

So the cheapest health check anybody writes — fetch the root, check the status is 2xx — reports this API as healthy *and* reports it identically whether or not your credential works. A monitor built that way goes green through an outage of everything behind the front door, and through a revoked key.

**There are two fields for the caller's own identifier.** `external_id` and `unique_external_id`, both nullable, neither required, one of them promising uniqueness in its name and nothing in the schema enforcing it.

**The status has six values and two of them read as "not started".** `open`, `pending`, `in_audit`, `in_review`, `approved`, `rejected`. A client deciding whether to let a business transact has to know which of the four non-terminal values it may proceed on, and the enum does not order them.

## Modelling limits

- **One route.** Listing businesses. Orders, liens, registrations, risk assessments, monitoring, documents, the timeline and the batch surface each want their own evidence — sixty-one paths in the document, one modelled here.
- **`object` is served as `list` and `business`.** The document declares the field required, types it a plain string, and gives no enum and no example, so it never says which string. The values here are the convention the envelope shape comes from rather than anything Middesk has written down, and that is a guess this Recipe is stating rather than hiding.
- **Nothing is mapped in detection.** No client for this API on npm, Packagist or the Go module proxy under an obvious name, checked 2026-09-09.
