# pylon

Emulates the Pylon account listing for local development and tests.

**10 conformance cases, 4 checked against the live API on 2026-09-13.**

Read from the OpenAPI document Pylon serves without a credential at [`api.usepylon.com/openapi.json`](https://api.usepylon.com/openapi.json), and struck live on 2026-09-13 with no credential, with a wrong token, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**A refusal cites a section of an RFC.**

```json
{"errors":["Token must follow Bearer authorization scheme: https://www.rfc-editor.org/rfc/rfc6750#section-2.1."],
 "request_id":"4ba417ae-…","code":"invalid_authorization_header"}
```

A URL with an anchor, pointing at RFC 6750 §2.1, inside the sentence a client prints — said to a request that sent no header to be wrong about.

**And two of the four failures end in an exclamation mark.** `Invalid API token!` and `Not a valid API URL!` do; `Method not allowed` and the RFC one do not. One API, one error array, two registers.

**The messages are bare strings in an array.** `errors` is a list of sentences, not of objects, so there is nothing machine-readable inside it — and the thing a client would branch on sits outside, in `code`, beside a `request_id` that every response carries, failures and successes alike.

**The document names the vendor's Go packages.** Every schema carries `x-go-name` and `x-go-package`: `Account` is `"x-go-package": "pylon/api/apiserver/apitypes"` and the response body is `"pylon/api/apiserver/endpoints"`. The published contract states the internal module layout and the Go identifier for each field.

**Three fields hold the domain and two of them are documented identically.** `domain` is "The primary domain of the account." `primary_domain` is "The primary domain of the account." — the same sentence, word for word, on two fields, with `domains` beside them holding the list.

**Nothing on an account is required.** `Account` has no `required` array at all, so `id` is optional in the contract that describes it.

**One resource, two names for its identifier.** Six paths take `{account_id}` — highlights, relationships, notebook blocks — and three take `{id}`, for the same account. And `/accounts/{account_id}/relationship` is a POST while `/accounts/{account_id}/relationships` is a GET, so the singular creates and the plural lists.

**And the listing declares 200, 400, 403, 404 and 500.** Not the 401 every anonymous caller meets, and not the 405 a wrong method gets — so `cauldron drift` reports both as unbacked. The document also declares a US server and an EU one, with no field on any record saying which answered.

## Sources

- [`api.usepylon.com/openapi.json`](https://api.usepylon.com/openapi.json) — served without a credential; recorded by `cauldron drift`, `spec_seen: 2026-09-13`.
- Live: `api.usepylon.com`, struck 2026-09-13 with no credential, a wrong token, an unrouted path, and a wrong method.

## Modelling limits

- **One route of a hundred and nine.** The account listing. Issues, contacts, teams, tags, knowledge base, macros, custom fields, attachments and the notebook surface are the rest.
- **`request_id` is a constant.** Live it is a fresh uuid on every response; the live case asserts its shape with a regex.
- **One region.** The document declares a US host and an EU one; this Recipe serves one, and the finding is that no record says which answered.
- **The success fixture is document-derived.** Listing accounts needs a real token; the records here are `Account`'s own properties with values of the declared types, and every case reading them is marked documentation-only.
- **Nothing is mapped in detection.** Pylon is reached through a plain HTTP call carrying a bearer token, which does not resolve to this host through a dependency file. Checked 2026-09-13.
