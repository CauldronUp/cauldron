# prismic

Emulates the Prismic repository document for local development and tests.

**8 conformance cases, 7 checked against the live API on 2026-09-13.**

The repository document — the thing every Prismic client fetches before it can do anything else — struck live against a repository's CDN host on 2026-09-13, and against a repository that does not exist.

## What this Recipe found

**A repository that does not exist answers zero bytes labelled as JSON.**

```
HTTP/1.1 404 Not Found
Content-Type: application/json
Content-Length: 0
```

The label promises a parseable body and there is nothing to parse. `.json()` throws on an empty string while the header says otherwise, so a client that trusts the content type gets a parse error rather than a 404 it can report.

**The API describes itself as HTML forms.** The document's `forms` object holds three of them, and each is an HTML `<form>` element serialised into JSON:

```json
"everything": {
  "method": "GET",
  "enctype": "application/x-www-form-urlencoded",
  "action": "https://<repo>.cdn.prismic.io/api/v2/documents/search",
  "fields": { "ref": {…}, "q": {…}, "page": {…}, … }
}
```

A method, an encoding type, an action and a field list — what a browser reads out of markup. A client is meant to fetch this document, find the form, and build the request from its fields: hypermedia by way of the 1999 web form.

**The credential is one of the form's fields.** `access_token` sits in the same `fields` map as `q`, `page` and `orderings`, typed `String`. The secret is declared as a query parameter, in a machine-readable description handed to anyone who asks without one.

**And everything else about the repository comes with it.** Every content type and its label, every language, the bookmarks, the tags, the A/B experiments, the OAuth endpoints and the master ref — to a caller who has presented nothing. Knowing the subdomain is the whole of the access control on this document.

**Integer fields have string defaults.** `page` and `pageSize` are declared `"type": "Integer"` with `"default": "1"` and `"default": "20"`. The type says integer and the default is quoted, so a client that reads the default and uses it without converting sends a string where the same form says a number belongs.

**Every query needs a `ref`, and the only way to get one is this document.** The master ref is a token that changes on every publish, so a read is always two requests: one for the ref and one for the content. A client that caches the ref serves stale content; one that does not doubles its request count.

**The OAuth endpoints are on a different host.** `oauth_initiate` points at `<repo>.prismic.io` while the document itself is served from `<repo>.cdn.prismic.io`. The CDN hands out the origin's addresses.

**And the API's version is a commit.** `"version": "482083b"` — seven hex characters, a git short hash, in the field a client would compare against a release.

## Sources

- Live: a Prismic repository's CDN host and one that does not exist, struck 2026-09-13.

## Modelling limits

- **One route.** The repository document. The search endpoint it describes, the GraphQL endpoint, the OAuth flow and the Migration API each want their own evidence.
- **The empty 404 is triggered differently here.** Live it is what an unknown *repository* answers, which is a subdomain rather than a path; Cauldron serves one host, so the same shape — 404, `application/json`, zero bytes — is what an unrouted path answers here. The shape is struck live and the trigger is not.
- **Three forms exist and one is served.** `everything`, `collections` and `tags` are all in the live document; the fixture carries `everything`, which is the one every client reads.
- **The fixture's repository name is a placeholder.** The document's `action` and OAuth URLs name a repository, and the ones here are shaped like real Prismic addresses without being any account's.
- **No `spec:`.** Prismic describes this API in prose. The machine-readable description it does publish is the `forms` object inside the document itself, which describes the *next* request rather than this one.
- **Nothing is mapped in detection.** The `@prismicio/client` packages are real clients, but every one of them is pointed at a repository subdomain that lives in configuration rather than in the dependency — and the subdomain is the only thing an emulator would need. Checked 2026-09-13.
