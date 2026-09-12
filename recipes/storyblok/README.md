# storyblok

Emulates the Storyblok Content Delivery API for local development and tests.

**10 conformance cases, 3 checked against the live API on 2026-09-11.**

Read from Storyblok's own documentation — every page of which answers Markdown when `.md` is appended to its URL — and struck live against `api.storyblok.com` on 2026-09-11 with no credential, with a deliberately invalid one, and on a path that does not exist.

## What this Recipe found

**One API answers an object for one failure and a bare array for the other.**

```
/v2/cdn/stories                 401 {"error":"Unauthorized"}
/v2/cdn/stories?token=notreal   401 {"error":"Unauthorized"}
/v2/cdn/cauldron-nope           404 ["This record could not be found"]
```

The 401 is an object with a key. The 404 is a top-level JSON array holding one string. A client reading `body.error` gets the sentence on one and `undefined` on the other; a client looping over the array gets one entry from the 404 and nothing at all from the 401.

This is the first bare array of *strings* in the collection — Salesforce answers a bare array of objects, and [datadog](../datadog) answers the same strings inside an envelope. Storyblok is the pair of both, in one API, chosen by which failure it is.

The same sentence answers a missing token and a wrong one, so there is no way to tell a caller who forgot the credential from one whose credential is bad.

**The credential is a query parameter, and only a query parameter.** The documentation's authentication section says it plainly — the access token is provided as a query parameter. So the secret is in the URL on every request: in access logs, in proxy logs, in `Referer` headers on any absolute link the content carries, and in the browser history of anyone who opens one by hand.

**The rate limit falls as the page grows, and three of its four boundaries land in two tiers at once.**

| listing size | limit |
| --- | --- |
| ≤ 25 entries | 50 per second |
| 25 to 50 entries | 15 per second |
| 50 to 75 entries | 10 per second |
| 75 to 100 entries | 6 per second |

25 is in the first row and the second. 50 is in the second and the third. 75 is in the third and the fourth. Each of those requests has two published limits, and 25 and 50 are the two page sizes anybody actually picks — 25 because it is the default, 50 because it is the round number above it.

**Pagination travels outside the JSON.** `total` and `per_page` come back as response *headers*, on an endpoint whose body is an envelope with room for them. The body's only envelope field is `cv`, a cache-version Unix timestamp. So the two numbers a client needs in order to page are the two it cannot read from the parsed response.

**Every story carries two identifiers, and its content carries a third.** An integer `id` and a `uuid` sit side by side on the record, and every block inside `content` has a `_uid` of its own — three identifier schemes on one document, with other endpoints keying by different ones.

**And one record spells "nothing" two ways.** The homepage in Storyblok's own example has `parent_id: 0` — zero meaning no parent — beside six fields saying null for the same idea: `sort_by_date`, `meta_data`, `release_id`, `path`, `default_full_slug` and `translated_slugs`. A client testing for absence needs both checks, and `parent_id` is the one where the falsy test and the null test agree by accident.

**`lang` is the string `default`.** Not a language tag, not null, not the space's actual default language — the literal word, where an IETF tag belongs.

**And there are nine base URLs.** Five regional hosts, four more for enterprise plans, and the Chinese one on a different domain entirely (`app.storyblokchina.cn`). Nothing in a token says which one serves it, so the host is configuration a client has to be told separately from its credential.

## Sources

- Live: `api.storyblok.com`, struck 2026-09-11.
- [Introduction to the Content Delivery API](https://www.storyblok.com/docs/api/content-delivery/v2) — base URLs, authentication, rate limits, pagination.
- [Retrieve Multiple Stories](https://www.storyblok.com/docs/api/content-delivery/v2/stories/retrieve-multiple-stories) — the parameters and the record shape.

## What it added to Cauldron

`responses.error.key: "-"` now works for the `string_list` style, meaning the array of sentences is the whole body rather than the contents of an envelope. The `list` style already spelled a bare array that way for Salesforce; this gives `string_list` the same spelling rather than inventing a second one.

## Modelling limits

- **One route.** Stories. Links, datasources, datasource entries, tags, spaces, assets and the whole Management API each want their own evidence.
- **`per_page` is not echoed.** Storyblok sends the requested page size back as a response header beside `total`. Cauldron can declare a header carrying the total (`count_header`) and one carrying the page count (`pages_header`), and has no key for a header echoing the page size. One provider is a note; if a second turns up, that is the field name.
- **The regional hosts are named, not served.** This Recipe models the European base URL. The others differ only in host, except the Chinese one, which is on another domain.
- **No `spec:`.** Storyblok publishes this API as prose, one page per operation — each also answering Markdown at the same address with `.md` appended. There is no OpenAPI document at any address this Recipe could find.
- **Nothing is mapped in detection.** The clients on npm and Packagist — `storyblok-js-client`, the framework SDKs, the PHP client — are pointed at whichever of the nine base URLs a space lives on, and the token does not say which. So a dependency name says the project uses Storyblok and not which host to stand in for. Checked 2026-09-11.
