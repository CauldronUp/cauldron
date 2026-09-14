# checklyhq

Emulates the Checkly check listing for local development and tests.

**10 conformance cases, 4 checked against the live API on 2026-09-13.**

Read from the OpenAPI document Checkly serves without a credential at [`api.checklyhq.com/openapi.json`](https://api.checklyhq.com/openapi.json), and struck live on 2026-09-13 with no header, with a wrong bearer, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**The published cURL example does not send the header it is about.** From the security scheme's own description:

```
curl -H "Authorization: Bearer [apiKey]" "X-Checkly-Account: [accountId]"
```

There is no `-H` on the second string, so curl reads it as the URL. Copied and pasted, that command tries to fetch a host called `X-Checkly-Account` and never sends the account header at all — which is the header this API refuses requests for lacking.

**And that description is HTML in a machine-readable field.** It carries an `<a href=… target="_blank">`, a `<b>`, a `<code>`, a `<br>`, and five `</br>` — a closing tag for a void element, which is not HTML at all. A generated client prints the markup.

**A refusal repeats itself in a fourth key.** With a wrong token:

```json
{"statusCode":401,"error":"Unauthorized","message":"Bad token","attributes":{"error":"Bad token"}}
```

`statusCode` restates the status line, `error` restates its reason phrase, and `attributes.error` restates `message`. Four keys, one fact, and the fourth is the framework's internal bag surfacing on the wire.

**A wrong method is the same 404 as a wrong path.** Both answer `{"statusCode":404,"error":"Not Found","message":"Not Found"}` — where `error` and `message` are the same two words again. The operation declares a 400, a 401, a 403 and a 429 and no 404, so `cauldron drift` reports both as unbacked.

**The listing is an `anyOf` of thirteen shapes with no discriminator.** Heartbeat, Playwright, ICMP, traceroute, gRPC, SSL, API, TCP, URL, DNS, browser, multi-step and generic — and nothing in the schema tells a client which one it is holding. Each member carries its own `checkType` enum of one, so the discriminator exists thirteen times and is declared nowhere.

**Thirty-nine properties and one of them is required.** On an API check, only `name` must be present: not `id`, not `frequency`, not `locations`.

**Two of those thirty-nine are snake_case.** `created_at` and `updated_at`, among `doubleCheck`, `shouldFail`, `useGlobalAlertSettings`, `groupOrder`, `runtimeId`, `aiAutoRepairEnabled` and thirty-one more.

**`frequency` is a closed list of minutes that includes zero.** `[0, 1, 2, 5, 10, 15, 30, 60, 120, 180, 360, 720, 1440]`, default 10 — so "how often" has a legal value that is not a frequency.

**Three fields have three states each.** `member`, `pending` and `logicalId` are `nullable: true`, and their descriptions say what the null means: "Null when not managed by code." A field called `member` is documented as "True when the project owns this check" — membership standing in for ownership — and `pending` is "True when the binding is reserved by an import plan that has not been deployed yet", which is the deployment tool's internal state, on the public record.

**`bearerFormat` is `"Bearer"`**, repeating the scheme name where the format belongs; the account identifier travels as an ordinary header parameter on every operation rather than as part of the security scheme; and `intent.goal` is a two-thousand-character prose field for "the user or system outcome this check protects", on a monitoring check.

## Sources

- [`api.checklyhq.com/openapi.json`](https://api.checklyhq.com/openapi.json) — served without a credential; recorded by `cauldron drift`, `spec_seen: 2026-09-13`.
- Live: `api.checklyhq.com`, struck 2026-09-13 with no header, a wrong bearer, an unrouted path, and a wrong method.

## Modelling limits

- **One route of a hundred and sixty.** The check listing. Check groups, alert channels, dashboards, maintenance windows, snippets, environment variables, private locations, reporting, incidents and the status-page surface are the rest.
- **One shape of thirteen.** The listing's `anyOf` mixes thirteen check types; this Recipe serves records carrying `checkType` from two of them, with the field set of the API check, and records that the union has no discriminator rather than emulating all thirteen.
- **The account header is not required here.** Live it is checked alongside the bearer; the finding is that the example teaching you to send it does not.
- **The success fixture is document-derived.** Listing checks needs a real key; the records here are `ChecksV1ApiCheck`'s own properties with values from its declared enums, and every case reading them is marked documentation-only.
- **Nothing is mapped in detection.** Checkly is reached through `checkly` on npm or a plain HTTP call carrying a bearer token, and neither resolves to this host through a dependency file. Checked 2026-09-13.
