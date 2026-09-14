# lambdatest

Emulates the LambdaTest build listing for local development and tests.

**11 conformance cases, 5 checked against the live API on 2026-09-14.**

Read from the automation API reference, and struck live on 2026-09-14 with no credential, with a wrong basic credential, on a path that does not exist, with a method the path does not take, on a version that does not exist, and on the root.

## What this Recipe found

**BrowserStack and LambdaTest refuse with the same sentence from the same middleware, and no client can treat the two the same way.** Both answer `HTTP Basic: Access denied.` — Rack's default. They differ in every other respect:

| | LambdaTest | BrowserStack |
|---|---|---|
| content type | `text/plain; charset=utf-8` | `text/html; charset=utf-8` |
| body length | 26 bytes | 27 bytes (trailing newline) |
| realm | `"Authorization Required"` | `"Application"` |
| header spelling | `www-authenticate` | `WWW-Authenticate` |

Two competitors, one library, four differences. And LambdaTest's realm is the name of the status code, so the prompt a browser shows the user says "Authorization Required" where it should say what they are logging in to. See the [`browserstack`](../browserstack) Recipe for the other half.

**One media type, two spellings, on one host.** The 401 declares `Content-Type: text/plain; charset=utf-8`. The 404 declares `Content-Type: text/plain`, with no charset. Same server, same media type, different headers.

**A path that does not exist, a verb a path does not take, and a version that was never released are one answer.** `/automation/api/v1/cauldron-nope`, `PUT /automation/api/v1/builds` and `/automation/api/v2/builds` all answer Go's default `404 page not found`, before any credential is read.

**And the API's own root is blocked by its CDN.** `GET /` answers `403` with Cloudflare's interstitial — "Sorry, you have been blocked", `<title>Attention Required! | Cloudflare</title>`, 4,549 bytes, carrying conditional comments for Internet Explorer 6.

**One key in the envelope is PascalCase and everything else is snake_case.** `Meta.attributes.org_id` and `Meta.result_set.{count,limit,offset,total}` sit above a `data` array whose records are `build_id`, `status_ind`, `create_timestamp`. One capital letter, in the key a client reads first.

**A field is named for a database column convention.** `status_ind` — indicator — carrying `"completed"`, on a public API record.

**Asking for public URLs changes what another parameter may be.** `publicurl` is documented "limit cannot exceed 20 when true", so the maximum page size depends on the value of a different, unrelated parameter. A caller who sets `limit=50` and then flips `publicurl` on finds out at runtime.

**`sort` is a mini-language, and its own example sorts by a field the record does not have.** The documented value is `asc.user_id,desc.org_id`: direction and field joined by a dot, pairs joined by commas. `org_id` appears nowhere on a build — it is in `Meta.attributes`.

Also pinned: three identifiers on one record, with `build_id` and `user_id` numbers and `project_id` the string `"ML"`; `duration` is a bare number with no unit stated anywhere; and the reference for this API is published under a different brand — `lambdatest.com/support/api-doc/` serves TestMu AI, whose canonical pages are on `testmuai.com`, while the API host is still `api.lambdatest.com`.

## Sources

- [Fetch all builds of an account](https://www.testmuai.com/support/api-doc/selenium-automation-api/build/fetch-all-builds-of-an-account/) — the response example, the field list, and the `offset` / `limit` / `status` / `fromdate` / `todate` / `sort` / `publicurl` / `username` parameters.
- Live: `api.lambdatest.com`, struck 2026-09-14 with no credential, a wrong Basic credential, an unrouted path, a wrong method, `/automation/api/v2/builds`, and the root.

## Modelling limits

- **No description is published.** LambdaTest serves no OpenAPI document for the automation API, so there is no `spec` to fingerprint.
- **The success side is document-derived.** Listing builds needs a real username and access key; the records here are the reference's own example plus a second of the same shape, and every case reading them is marked documentation-only.
- **`Meta.result_set.offset` is not served.** This format names a page and a limit and a count; it has no name for an echoed row offset, and inventing a constant for it would be worse than leaving it out. The other three `result_set` fields are served.
- **`publicurl`, `sort`, `fromdate` and `todate` are not modelled.** The first two are the findings above; the date filters would need a record field this Recipe does not resolve dates against.
- **The Cloudflare 403 on the root is described, not served.** It is the CDN's page, not the API's, and reproducing someone's interstitial adds nothing the sentence above does not.
- **One route of many.** The build listing. Sessions, tests, platforms, screenshots and the App Automate surface are the rest.
- **Nothing is mapped in detection.** LambdaTest is reached through `lambdatest-node-tunnel` or a plain Basic call, and neither resolves to this host through a dependency file. Checked 2026-09-14.
