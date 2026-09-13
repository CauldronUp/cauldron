# datocms

Emulates the DatoCMS Content Management API for local development and tests.

**8 conformance cases, 3 checked against the live API on 2026-09-13.**

Read from DatoCMS's own documentation — every page of which answers Markdown at the same address with `.md` appended — and struck live against `site-api.datocms.com` on 2026-09-13 with no `Accept` header, with no credential, with a deliberately invalid one, with the header the documentation's own example sends, and on a path that does not exist.

## What this Recipe found

**Failures arrive under `data`.** All four, struck live:

```
(no Accept)           406 {"data":[{"id":"29801c","type":"api_error","attributes":{"code":"INVALID_ACCEPT_HEADER",…}}]}
Accept, no token      401 {"data":[{"id":"02f293","type":"api_error","attributes":{"code":"INVALID_SITE",…}}]}
Accept, wrong token   401 {"data":[{"id":"5413bc","type":"api_error","attributes":{"code":"INVALID_AUTHORIZATION_HEADER",…}}]}
/cauldron-nope        404 {"data":[{"id":"cf3acb","type":"api_error","attributes":{"code":"INVALID_ENDPOINT",…}}]}
```

This is JSON:API, where failures belong in `errors` and the specification says `data` and `errors` "MUST NOT coexist in the same document". DatoCMS puts the errors *in* `data` — the same key, holding the same kind of thing, an array of resource objects.

So `body.data[0]` is a record on success and an error on failure, and the only way to tell is to read `type`.

**There is no message anywhere in a failure.** No `title`, no `detail`, no sentence. A code, an always-empty `details: {}`, and a `doc_url`. The only human-readable thing in the response is a URL a reader has to go and fetch.

**And the `id` is a second error code.** JSON:API defines it as "a unique identifier for this particular occurrence of the problem". Struck twice, the same failure answers `02f293` both times — it is stable per code, sits beside an actual `code` field, and is six hex characters that appear nowhere in the documentation. The record carries two codes, and the useless one is in the field a JSON:API client is told to treat as an occurrence id.

**A missing credential is reported as a problem with the project.** With an `Accept` header and no token at all, the code is `INVALID_SITE` — not a missing header, not an invalid token, but the *site*. A wrong token gets `INVALID_AUTHORIZATION_HEADER`. The failure meaning "you sent no credential" names a different noun than the one meaning "your credential is wrong".

**The documentation's own example sends the wrong header.** From the pagination page, verbatim:

```bash
curl \
  -H 'Accept: application/json' \
  -H 'Authentication: Bearer <YOUR-API-TOKEN>' \
  https://site-api.datocms.com/items?page[limit]=5&page[offset]=5
```

`Authentication`, not `Authorization`. Struck live, that header is ignored and the answer is `INVALID_SITE` — so a reader who copies the worked example is told their project is wrong, and nothing anywhere in the response points at the credential.

**Content negotiation fails before anything else.** Without an `Accept` header the answer is `406 INVALID_ACCEPT_HEADER`, whatever credential is attached. The first failure a naive client meets is about a header it did not know it had to send.

**The product renamed these and the API did not.** The documentation's first sentence:

> DatoCMS stores the individual pieces of content you create from a model as records (for backwards compatibility the API calls these `item`).

**And a field is a string or an object depending on how you asked.** A link field arrives as `"featured_block": "dhVR2HqgRVCTGFi_0bWqLqA"` or as `"featured_block": {id, type, attributes, relationships}`, chosen by the `nested` request parameter. A client has to test the type of a field before it can read it, and the type depends on the request rather than on the schema.

## Sources

- Live: `site-api.datocms.com`, struck 2026-09-13.
- [Record](https://www.datocms.com/docs/content-management-api/resources/item) — the resource object, and the rename.
- [List all records](https://www.datocms.com/docs/content-management-api/resources/item/instances) — the listing.
- [Pagination](https://www.datocms.com/docs/content-management-api/pagination) — `page[limit]`, `page[offset]`, `meta.total_count`, and the `Authentication` header.

## Modelling limits

- **One route.** Listing records. Models, fields, fieldsets, uploads, roles, webhooks, environments, build triggers and the rest of the management surface each want their own evidence.
- **The 406 is recorded, not served.** Cauldron answers one shape per route, and DatoCMS's depends on whether an `Accept` header arrived. Every case here sends one; the 406 is quoted above because it is the finding.
- **Nothing declares a page size.** The documentation says both the default `page[limit]` and its maximum "vary depending on the specific endpoint", and the listing page names neither. The 30 here is a sandbox default, not a claim about DatoCMS.
- **The record's attributes are the site's, not the API's.** A record's fields come from the model it was made from, so the two seeded here share a shape only because they share a model. Nothing about `title` or `publication_date` is part of this API.
- **`nested` is not modelled.** The parameter that turns a link field from an id into a whole object also drops the maximum page size from 500 to 30, and both halves want their own evidence.
- **No `spec:`.** DatoCMS documents this API as prose with worked examples. There is no OpenAPI document at any address this Recipe could find.
- **Nothing is mapped in detection.** The `@datocms/cma-client` packages are real clients of this API, but a project holding one is pointed at whichever site its token belongs to — and DatoCMS's read traffic goes to a different host entirely, `graphql.datocms.com`, which this Recipe does not model. Checked 2026-09-13.
