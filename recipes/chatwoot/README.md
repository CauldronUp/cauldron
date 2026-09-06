# chatwoot

Emulates the Chatwoot help-desk API, for local development and tests.

**13 conformance cases, 2 checked against the live API on 2026-09-06.**

Written against Chatwoot's own generated OpenAPI 3.1 document, published in its own repository — `chatwoot/chatwoot`, `swagger/swagger.json`, 92 paths, version 1.1.0 — read on 2026-09-06. The two refusals were struck live against `app.chatwoot.com` the same day with no account: being refused needs no sign-up, which is the point of recording it.

## What this Recipe found

**Two listings, two envelope depths, one API version.** Contacts put the payload at the top:

```json
{"meta": {...}, "payload": [ ... ]}
```

Conversations put the same two keys one level further down:

```json
{"data": {"meta": {...}, "payload": [ ... ]}}
```

Adjacent paths under the same `/api/v1/accounts/{account_id}` prefix, the same document, the same version. `body.payload` reads a list of contacts and `undefined` for conversations — and a helper written to unwrap one of them returns `undefined` rather than throwing.

**And `meta` means two different things.** On contacts it is pagination: `{"count": 42, "current_page": 1}`. On conversations it is a set of assignment tallies — `mine_count`, `unassigned_count`, `assigned_count`, `all_count` — and carries no pagination at all. There is no total and no page number in a conversations listing, so a client cannot tell how many pages there are. `page` is the only parameter, and the only way to learn you have reached the end is to receive a short page. Reading `meta.count` there is `undefined`, which arithmetic turns into `NaN`.

Worse, `all_count` is a plausible-looking total that ignores every filter the request just sent — it counts the account, not the query.

**`current_page` is documented as a string *or* an integer.** The spec's own type is `["string", "integer"]`: one field, either type, per response. `meta.current_page + 1` is `2` when it arrives as a number and `"11"` when it arrives as a string.

**Refusals come in two shapes, and neither is the documented one.** Struck live:

```
GET /api/v1/accounts/1/contacts               (no credential)
HTTP 401  {"errors":["You need to sign in or sign up before continuing."]}

GET /api/v1/accounts/1/contacts               (api_access_token: bogus)
HTTP 401  {"error":"Invalid Access Token"}
```

Plural `errors` holding an array of strings when nothing was sent; singular `error` holding one string when something was sent and rejected. Same status, same endpoint, same second. A client with one error reader gets `undefined` for one of the two — and reports the reason as "undefined" exactly when the credential is wrong, the case it most needs to name. The spec documents a *third* shape for failures (`{"description": ..., "errors": [ ... ]}` with objects rather than strings); neither 401 uses it. Neither carries `WWW-Authenticate`.

**One header name, three kinds of credential.** The spec declares `userApiKey`, `agentBotApiKey` and `platformAppApiKey`. All three are `in: header` and all three are named `api_access_token`. Nothing in the request says which kind was sent — the server decides what you may reach by looking the token up. An integration holding the wrong kind is refused per endpoint rather than at the door, and the refusal reads as a permission problem rather than as the wrong credential.

**`message_type` is an integer with no name on the wire.** The enum is `0, 1, 2, 3` — incoming, outgoing, activity, template — and the response carries the digit. Filtering out activity messages means a literal `2` whose meaning lives in the documentation.

**Times are Unix numbers, and a conversation has two of them meaning the same thing.** `created_at` and `timestamp` carry identical descriptions in the spec — "The time at which conversation was created" — so which is safe to sort by cannot be read off the document.

**The status words are Chatwoot's own.** `open`, `resolved`, `pending`. Zendesk, Front and Help Scout all say *closed* for the state Chatwoot calls *resolved*, so a shared status mapper needs a Chatwoot branch.

## Modelling limits

- **Contacts and conversations on one account, listed and fetched.** 92 paths is a help desk: messages, inboxes, agents, teams, labels, canned responses, automation rules and the public per-inbox API each want their own evidence.
- **The three credential kinds are not told apart.** Modelling that would mean inventing which endpoints each may reach, and the spec's per-operation security names only one at a time.
- **The conversations listing is served without pagination metadata**, because that is what the provider does. It is paged here so the `page` parameter is real; the absence of a total is the finding, not a gap.
- **`meta.current_page` is served as an integer.** Serving it as sometimes-a-string would make the emulator nondeterministic, which is worse than stating that the provider's type is both.
