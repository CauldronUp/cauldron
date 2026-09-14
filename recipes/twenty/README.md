# twenty

Emulates the Twenty published OpenAPI document and the REST surface behind it, for local development and tests.

**9 conformance cases, 6 checked against the live API on 2026-09-13.**

Read from the document Twenty serves without a credential at [`api.twenty.com/open-api/core`](https://api.twenty.com/open-api/core), and struck live on 2026-09-13 with no header, with a wrong bearer, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**The contract describes one operation, and it is fetching the contract.** `paths` has a single entry: `GET /open-api/core`, `operationId GetOpenApiSchema`. Every other path — people, companies, notes, tasks — is generated per workspace and appears only once a token is presented, so an anonymous reader of the published document learns how to fetch the document and nothing else.

**And the document's only operation contradicts the document's own security.** The top level says `security: [{bearerAuth: []}]`, so every operation requires a bearer. `GET /open-api/core` answers 200 to a request with no header at all.

**It also has to override the server to reach itself.** The document's server is `https://api.twenty.com/rest/`; the one operation declares its own `servers: [{url: "https://api.twenty.com"}]`, because the path it lives at is not under the base URL the document declares. `cauldron drift` makes the same point from the other side: it reports `GET /open-api/core` as undeclared, because joining the document's server to the document's path gives `/rest/open-api/core`, which is not where the document lives.

**A missing token is 403 and a wrong token is 401.**

```
(no header)      403  {"statusCode":403,"messages":["Missing authentication token"],"error":"FORBIDDEN_EXCEPTION"}
Bearer notreal   401  {"statusCode":401,"messages":["Token invalid."],"error":"UNAUTHENTICATED"}
```

The two statuses are the wrong way round. 401 is for a caller who has not authenticated; 403 is for one who has and may not. Twenty answers 403 to the request carrying nothing and 401 to the one carrying something.

**And the two error codes come from different worlds.** `FORBIDDEN_EXCEPTION` is a class name with a suffix; `UNAUTHENTICATED` is a gRPC status name. One API, two vocabularies, in the field a client switches on — with `messages`, a plural array, carrying one string beside them.

**The security scheme's own instructions produce an invalid header.** `bearerAuth.description` reads: ``Enter the token with the `Bearer: ` prefix, e.g. "Bearer abcde12345".`` A colon after `Bearer` is not what the header looks like, and the example beside it does not have one.

**The base URL is labelled "Production Development".** Two environment words in one server description, in a document with one server.

**Filtering is a language you build by string concatenation.** From the document's own prose: `field[COMPARATOR]:value`, joined with commas, with `field.subField` for composite fields and `%` as a wildcard for `like` and `ilike`. The comparator goes in square brackets inside the *value* of a query parameter, so every filter is assembled by hand and escaped by nobody.

**`components` holds nothing but the security scheme.** No schemas, no parameters, no responses — and the one operation's 200 is an inline, hand-written description of what an OpenAPI document looks like.

**The credential is checked before the path**, so an unrouted path and a wrong method both answer "Missing authentication token". And the document's `info.description` carries a paragraph warning "Never put your token in a URL", with its reasoning about access logs, browser history and `Referer` headers — in the prose field a code generator drops.

## Sources

- [`api.twenty.com/open-api/core`](https://api.twenty.com/open-api/core) — served without a credential; recorded by `cauldron drift`, `spec_seen: 2026-09-13`.
- Live: `api.twenty.com`, struck 2026-09-13 with no header, a wrong bearer, an unrouted path, and a wrong method.

## Modelling limits

- **Two routes.** The document, and one REST path kept for its failures.
- **`/rest/people` holds no key.** Twenty's object paths are generated per workspace and are not in the published document, so this Recipe declines to invent their records: the route accepts no credential, every request to it is refused, and the two refusals are what it is there for.
- **The document in the fixture is abridged.** The top-level shape, the single path, the server, the security scheme and the external docs are the live document's own; `info.description` is its first line rather than its full text.
- **Nothing is mapped in detection.** Twenty is reached through a plain HTTP call carrying a bearer token, and that does not resolve to this host through a dependency file. Checked 2026-09-13.
