# gotify

Emulates the Gotify notification-server API, for local development and tests.

**11 conformance cases, none checked against a live API.**

Written against Gotify's own generated Swagger 2.0 document, published in its own repository — `gotify/server`, `docs/spec.json`, 30 paths, version 2.1.0 — read on 2026-09-06. Gotify is self-hosted and there is no public instance, so every claim here is the provider's own document.

## What this Recipe found

**`since` means less than.** The document's own words for the parameter on `GET /message`:

> since — "return all messages with an ID less than this value"

In every other Recipe in this collection `since` means *after*: after this time, after this position, going forward. Here it pages backwards, because messages come newest-first and paging means walking into the past. A client that assumes the usual sense sends the id it started from and receives the same page again, or sends the newest id it has seen and receives everything it already had.

Nothing errors. `since` is a valid parameter with a valid value, and the response is a correct page of the wrong end.

**Seven security definitions, three wire positions, and the type is the first letter of the token.** The document declares `appTokenAuthorizationHeader`, `appTokenHeader`, `appTokenQuery`, `clientTokenAuthorizationHeader`, `clientTokenHeader`, `clientTokenQuery` and `basicAuth`. The app and client pairs are the same place:

| app scheme | client scheme | where it goes |
|---|---|---|
| `appTokenAuthorizationHeader` | `clientTokenAuthorizationHeader` | `Authorization` |
| `appTokenHeader` | `clientTokenHeader` | `X-Gotify-Key` |
| `appTokenQuery` | `clientTokenQuery` | `?token=` |

So which kind of token you sent is **not visible in the request at all**. It is the token's first character: the document's own examples are `Axxxxxxxxxx` for an application token and `Cxxxxxxxxxx` for a client token. The whole access-control model is a prefix letter.

And the two kinds do not reach the same endpoints. `POST /message` accepts both. `GET /message`, `DELETE /message`, `/application`, `/client` and `/current/user` accept a client token only. An application token can write and cannot read — a good design, and one where using the wrong of two identical-looking credentials answers 401 on an endpoint that plainly exists and that the other token just read.

**This Recipe serves that split** rather than only describing it: the read route accepts the client token alone, the write route accepts both.

**The credential may travel in the query string** — `?token=`, a declared scheme. Gotify is a self-hosted box commonly sat behind a proxy that logs full URLs.

**The next page is a relative path, not a URL.** `paging.next` is documented as "the relative path for the next page … Should be combined with the gotify base url", example `/message?limit=50&since=123456`. Every other provider modelled here that sends a next link sends an absolute URL. A client passing this one to a URL constructor resolves it against whatever base it holds — in a browser, the page's own origin.

**There is no total anywhere.** `paging.size` is how many came back, `paging.limit` is what was asked for, `paging.since` is the last id seen. Nothing says how many messages exist.

**The error body carries the status a second time, and one of its fields is documented as another.** An `Error` has `error`, `errorCode` and `errorDescription`. `errorCode` is the HTTP status repeated in the body. And the published description of `errorDescription` reads, in full, "The http error code." — the description of the field above it, copied. A generated client's documentation says the wrong thing about the only field carrying the actual explanation.

**`extras` is namespaced with a double colon.** `client::display.contentType` — not a JSON convention, and awkward in every language offering dotted access to a decoded object.

**The shipped document's host is `localhost`,** with no `basePath`. A client generated from it points at the machine it was generated on.

## Modelling limits

- **Nothing here is verified against a live API.** Gotify is self-hosted and there is no public instance.
- **Messages, listed and created.** 30 paths is a notification server: applications, clients, users, plugins and the websocket stream each want their own evidence.
- **Per-application scoping is not modelled.** Upstream, an application token may post only as its own application; here any app token posts as the fixture's.
- **The credential is checked in the `Authorization` header.** `X-Gotify-Key` and the query-string form are recorded and not served — serving the query form would make a credential in a URL the easy path in local development.
- **`paging.next` is served as the relative path the provider sends**, which is the finding. Nothing here resolves it.
