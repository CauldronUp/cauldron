# upsun

Emulates the Upsun organizations API for local development and tests.

**10 conformance cases, 3 checked against the live API on 2026-09-13.**

Read from the OpenAPI document Upsun serves without a credential at [`meta.upsun.com/openapi-spec`](https://meta.upsun.com/openapi-spec), and struck live against `api.upsun.com` on 2026-09-13 with no credential, with a deliberately invalid one, on the deprecated path, and on an organization that does not exist.

## What this Recipe found

**Four failure shapes, from one API.**

```
/projects                      400 application/problem+json
                               {"type":"http://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html#sec10.4.1",
                                "title":"Deprecated endpoint. See https://docs.upsun.com/api/",
                                "status":400,"detail":"Bad Request."}
/organizations                 401 {"error":"authentication_required","error_description":"Bearer authentication is required."}
/organizations (wrong token)   401 {"error":"invalid_token","error_description":"The access token is invalid."}
/organizations/<uuid>/projects 404 {"status":404,"title":"Not Found","message":"Organization not found."}
```

A problem document, an OAuth error, another OAuth error, and a fourth object that is half a problem document — `title` from RFC 9457 and `message` where `detail` belongs. Four shapes, two media types, and nothing in a response saying which to expect.

**The deprecation notice is typed with a dead specification.** RFC 9457 says `type` is a URI identifying the *problem kind*, and that dereferencing it should give human-readable documentation of that kind. This one points at section 10.4.1 of RFC 2616 — the 1999 HTTP/1.1 specification's definition of "400 Bad Request", obsoleted in 2014.

So the machine-readable identity of the problem is a link to a retired document's definition of the status code already sitting in the `status` field.

**And `title` and `detail` are the wrong way round.** RFC 9457 says `title` summarises the problem *type* and "SHOULD NOT change from occurrence to occurrence"; `detail` is the explanation specific to this occurrence. Here `title` carries the deprecation notice and a URL, and `detail` carries the generic phrase `Bad Request.` The stable field holds the specific message and the specific field holds the generic one.

**The notice points at the documentation of the API that sent it** — `docs.upsun.com/api/`, the whole reference, not a replacement endpoint. A client reading this has been told that something is deprecated and not what to call instead.

**And `/projects` is not in the description at all.** 187 paths, and the one answering the deprecation notice is not among them. `cauldron drift` says so:

```
upsun   not backed: no operation the Recipe routes to answers 401, which authentication_required declares
        not backed: no operation the Recipe routes to answers 401, which invalid_token declares
        not backed: the description does not declare GET /projects
```

**Whether you are asked to authenticate depends on whether the path has an id in it.** With no credential, `/organizations` answers 401 `authentication_required` and `/organizations/<a uuid>/projects` answers 404 `Organization not found.` — so identifier lookup runs before the credential check on the second, and an anonymous caller can tell a real organization id from a made-up one.

**The 401 is not declared.** `GET /organizations` lists 200, 400 and 403, and answers 401 to every request arriving without a token.

**`api.platform.sh` answers all of this byte for byte.** Two brands, one gateway, and a deprecation notice on one of them pointing at the other's documentation.

**The page size defaults to something its own type forbids.** `page[size]` is declared `{"type": "integer", "minimum": 1, "default": null}`. null is neither an integer nor at least 1.

**One record, two identifier formats.** `id` is a ULID and `owner_id` is a UUID, on the same object, for the organization and the person who owns it.

**Three fields are the organization's name.** `namespace` ("the namespace in which the organization name is unique"), `name` ("a unique machine name") and `label` ("the human-readable label"). Nothing says which a caller sends back.

**And nothing on the record is required.** Seventeen properties and no `required` array, `id` included.

## Sources

- Live: `api.upsun.com` and `api.platform.sh`, struck 2026-09-13.
- [`openapi-spec`](https://meta.upsun.com/openapi-spec) — served without a credential, 1.6MB, 187 paths.

## Modelling limits

- **Two routes.** Listing organizations, and the deprecated `/projects` that exists only to say it is gone. Projects, environments, deployments, domains, certificates, invoices, the observability surface and the rest are 185 more paths.
- **The organization-not-found ordering is recorded, not served.** `/organizations/<a uuid>/projects` answers 404 to an anonymous caller, which means the id is resolved before the credential is read. Modelling that would need a route answerable without a credential, and what a *real* organization id does anonymously was never tested.
- **The deprecated route is declared public** because that is how it behaves: a caller with a token gets the 401 instead, so the notice that the endpoint is gone is visible only to callers who are not signed in.
- **Nothing is mapped in detection.** Upsun and Platform.sh are driven by their own CLI and by a Terraform provider; the published client libraries are generated and rarely held directly, and either brand's CLI can be pointed at either host. Checked 2026-09-13.
