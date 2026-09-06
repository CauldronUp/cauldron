# shlink

Emulates the Shlink URL-shortener API (v3), for local development and tests.

**11 conformance cases, none checked against a live API.**

Written against Shlink's own OpenAPI 3.1 document, published in its own repository — `shlinkio/shlink`, `docs/swagger/`, 18 paths, version 3.0 — read on 2026-09-06. The document is split across files: every path, parameter and schema is a `$ref` to its own JSON file, so reading it means fetching about forty of them. Shlink is self-hosted and the public demo hosts did not answer, so every claim here is the provider's own document.

## What this Recipe found

**The API version is a path variable, and every path accepts every version.** The paths are `/rest/v{version}/short-urls`, and `version` is a path parameter:

```yaml
name: version
in: path
required: true
schema: {type: string, enum: ["3", "2", "1"]}
```

So `/rest/v1/short-urls` and `/rest/v3/short-urls` are the same operation with a different substitution, and nothing in the schema says which version introduced what.

The document *knows* — its own filenames record it: `v1_short-urls.json`, `v2_visits.json`, `v3_short-urls_{shortCode}_redirect-rules.json`. That knowledge is in the file layout and not in the API description, so `/rest/v1/short-urls/{shortCode}/redirect-rules` — a v3 feature at v1 — is a syntactically valid request a generated client will happily build.

**The path parameter is only half the identifier.**

```
GET    /rest/v{version}/short-urls/{shortCode}?domain=…
PATCH  /rest/v{version}/short-urls/{shortCode}?domain=…
DELETE /rest/v{version}/short-urls/{shortCode}?domain=…
```

A short code is unique *per domain*, not globally: the same code can exist on `example.com` and on the default domain at once. `domain` is a query parameter and it is `required: false`.

So **the URL of a DELETE does not say which record it removes.** Omitting the domain deletes the one on the default domain, which may not be the one the caller meant — and there is no error, because both are real short URLs. The identity of the thing being destroyed depends on a parameter you can forget.

**`domain: null` means the default domain, not the absence of one.** The field is `type: ["string", "null"]`, described as "Null if it belongs to default domain". The falsy value is a *specific* domain, and grouping short URLs by `domain` files every default-domain link under the key `null`.

**Query parameters with square brackets in their names.** `tags[]` and `excludeTags[]` — PHP's array convention, in the parameter's `name`. A conforming client percent-encodes them to `tags%5B%5D`, and the name a generator sees is `tags[]`, which is not an identifier in any language it targets. Every generator renames it, differently.

**The records are two levels down, under the resource's own name:** `{"shortUrls": {"data": [ … ], "pagination": { … }}}`. `body.data` is undefined; `body.shortUrls.data` is the list. A shared unwrapper has to know the resource before it can find the records.

**Five numbers in the pagination block, counting three different things:** `currentPage`, `pagesCount`, `itemsPerPage`, `itemsInCurrentPage`, `totalItems`. Two count items — this page and the whole set — one counts pages, and `itemsPerPage` is what was *asked for* rather than what arrived. Reading the wrong one of the two item counts sizes a progress bar to the page.

**A visit count that separates the robots.** `visitsSummary` is `{total, nonBots, bots}`, all three required, and `total` includes the bots. The obvious field to read is the flattering one; the honest one is `nonBots`, which a client only finds by looking.

**Three booleans that change what a redirect does.** `crawlable` puts it in robots.txt; `forwardQuery` decides whether the caller's query string survives; `hasRedirectRules` says the destination is *conditional* — and the record does not say on what. Only the third is a hint that the `longUrl` beside it may not be where anyone lands.

**Errors are RFC 7807 with real type URIs** — `https://shlink.io/api/error/invalid-api-key` — rather than `about:blank`, so unlike [peertube](../peertube)'s the `type` is actually a discriminator. Its `invalidElements` extension member is spelled without a hyphen, so it is a property access.

## Modelling limits

- **Nothing here is verified against a live API.** Shlink is self-hosted and the public demo hosts did not answer on 2026-09-06.
- **Short URLs, listed and fetched.** 18 paths is a URL shortener: tags, domains, visits, orphan visits, redirect rules, the Mercure surface and the redirect itself each want their own evidence — and `GET /{shortCode}`, the redirect, is the only endpoint that matters to an end user.
- **The version is served as a literal `v3`.** Routing all three would mean claiming this Recipe knows which endpoints each version really serves, which is exactly what the document does not say.
- **`domain` is declared and not routed.** Serving the same short code twice under two domains is the finding, and reproducing it would need two records with one identifier, which this format's store does not hold.
