# browserstack

Emulates the BrowserStack Automate build listing for local development and tests.

**9 conformance cases, 4 checked against the live API on 2026-09-14.**

Read from the Automate API reference at [`browserstack.com/docs`](https://www.browserstack.com/docs/automate/api-reference/selenium/build), and struck live on 2026-09-14 with no credential, with a wrong basic credential, without the `.json` suffix, on a path that does not exist, with a method the path does not take, and on the root.

## What this Recipe found

**One host, three different 404 bodies, depending on which layer answered.**

| request | body |
|---|---|
| `GET /automate/cauldron-nope.json` | nginx's default page, 162 bytes |
| `PUT /automate/builds.json` | BrowserStack's marketing 404 page, **259,615 bytes** |
| `GET /` | the same marketing page |

A wrong verb on a real API path returns a quarter of a megabyte of HTML titled "Page not found | BrowserStack", complete with `<!--[if IE 8]>` and `<!--[if IE 9]>` conditional comments — markup for browsers retired in 2016 — while a path that does not exist gets nginx's twelve-line default instead.

**The 401 is Rack's default, and the realm is the word "Application".**

```
HTTP/1.1 401 Unauthorized
WWW-Authenticate: Basic realm="Application"
Content-Type: text/html; charset=utf-8

HTTP Basic: Access denied.
```

Twenty-seven bytes of plain prose, declared `text/html`, from a JSON API. A client calling `.json()` on a wrong access key throws rather than reporting it, and the realm a browser shows the user is the generic word the framework ships with.

**LambdaTest, a direct competitor, sends the identical sentence differently.** `api.lambdatest.com` answers `HTTP Basic: Access denied.` too — the same Rack middleware, the same words — but as `text/plain; charset=utf-8`, without the trailing newline BrowserStack sends (26 bytes against 27), under `www-authenticate: Basic realm="Authorization Required"`. Two vendors, one library, and a client cannot treat the two refusals the same way on content type, byte length or realm. See the `lambdatest` Recipe.

**The listing is a bare array of one-key wrappers.** Every entry is `{"automation_build": {…}}`, so reading a build's name is `body[0].automation_build.name` — a wrapper repeated once per record, inside an array that is already typed by the endpoint it came from.

**The identifier is named after how it was made.** `hashed_id`, forty hex characters, is documented as the "Build ID required for other endpoints" — a public identifier whose name describes the implementation that produced it.

**And a credentialled listing hands back public links.** `public_url` is a "Shareable dashboard link" to the build's results, on every record, so anything that logs or forwards this listing forwards a URL that needs no credential to open.

Also pinned: `limit` defaults to 10 and caps at 100, which is the one place either number is written down; `build_tag` is null on the reference's own example; `duration` is a bare number documented as milliseconds only in the prose beside it; and `/automate/builds` without the `.json` suffix answers the same 401, so the extension is not what routes the request.

## Sources

- [BrowserStack Automate build API reference](https://www.browserstack.com/docs/automate/api-reference/selenium/build) — the response example, the field list, and the `limit` / `offset` / `status` / `projectId` parameters.
- Live: `api.browserstack.com`, struck 2026-09-14 with no credential, a wrong Basic credential, with and without the `.json` suffix, an unrouted path, a wrong method, and the root.

## Modelling limits

- **No description is published.** BrowserStack serves no OpenAPI document for the Automate API, so there is no `spec` to fingerprint.
- **The marketing 404 is seven lines of a 259,615-byte page.** This Recipe serves a minimal document carrying the same status, the same content type, the same `<title>` and the same IE conditional comments. The size is the finding and is recorded here rather than shipped.
- **The success side is document-derived.** Listing builds needs a real username and access key; the records here are the reference's own example plus a second of the same shape, and every case reading them is marked documentation-only.
- **One route of many.** The build listing. Sessions, projects, the plan endpoint, App Automate and the screenshot APIs are the rest.
- **`projectId` is not modelled.** It filters by a project this Recipe has no resource for.
- **Nothing is mapped in detection.** BrowserStack is reached through `browserstack-local` or a plain Basic call, and neither resolves to this host through a dependency file. Checked 2026-09-14.
