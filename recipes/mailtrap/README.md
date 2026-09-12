# mailtrap

Emulates the Mailtrap account-management API for local development and tests.

**8 conformance cases, 4 checked against the live API on 2026-09-11.**

Written against Mailtrap's own documentation — which is served as Markdown with an OpenAPI document embedded in each page — and struck live against `mailtrap.io` and `send.api.mailtrap.io` on 2026-09-11 with no credential, with a deliberately invalid one, and on a path that does not exist.

## What this Recipe found

**Three failure shapes, two hosts, one vendor.**

```
mailtrap.io/api/accounts        401 {"error":"Incorrect API token"}
mailtrap.io/api/cauldron-nope   404 {"errors":"API endpoint not found. …"}
send.api.mailtrap.io/api/send   401 {"success":false,"errors":["Unauthorized"]}
```

`error` holding a string. `errors` holding a string. `errors` holding an array.

The key changes between singular and plural, and **the plural key still carries one message** — so the name does not track the shape, and a client that reads `body.errors[0]` gets the first *character* of the 404's sentence.

The sending host adds a `success` boolean beside its array, which the other host never sends. A caller talking to both — the normal case, since one manages the account and the other sends the mail — needs two parsers for one vendor's failures.

**The 404 is addressed to a language model.**

```json
{"errors":"API endpoint not found. See Mailtrap documentation:
  https://docs.mailtrap.io/developers or https://docs.mailtrap.io/llms.txt for LLMs."}
```

Two links: one for a person and one, named as such, for an LLM. It is the first error message in this collection to route its reader by species.

The `llms.txt` it points at is real — 40KB of Markdown index — and it is how the record shape below was read. The API's own 404 gave better directions than its documentation site did: `api-docs.mailtrap.io/openapi.json` answers 1.6MB of HTML, which is exactly the case Cauldron's own parser has a guard for.

**And the documentation is built to be machine-read.** Every page answers Markdown when `.md` is appended to its URL, and each endpoint page embeds a complete OpenAPI document *for that one endpoint* as a fenced JSON block.

So there is no single description to fingerprint, and there is a precise description of every endpoint. That is an unusual trade and an apparently deliberate one — it costs `cauldron drift` its grip and buys a reader an exact answer per page.

**Nothing on an account is guaranteed.** The embedded schema has no `required` array at all, so `id`, `name` and `access_levels` are each optional on every record in a listing whose purpose is identifying accounts.

**`access_levels` is an array of numbers with no names.** Mailtrap's permission levels arrive as integers — `1000`, `100` — so a client deciding what a caller may do is comparing magic numbers it has to look up somewhere other than this response.

## Modelling limits

- **One route.** Accounts. Inboxes, messages, sending domains, contacts, campaigns and the whole sandbox surface each want their own evidence.
- **The sending host is quoted, not served.** `send.api.mailtrap.io` is a second host with its own error shape; it appears above as the contrast and is not modelled here.
- **No `spec:`.** Mailtrap publishes one OpenAPI document per endpoint, inside the Markdown of that endpoint's page. There is nothing single to fingerprint.
- **Nothing is mapped in detection.** The `mailtrap` packages on npm and Packagist are sending clients for the transactional API rather than clients of this account-management surface, so a dependency would offer the wrong host. Checked 2026-09-11.
