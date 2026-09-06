# streamchat

Emulates the Stream Chat API for local development and tests.

**14 conformance cases, all of them checked against the live API on 2026-09-06.**

Written against Stream's own OpenAPI 3.0.3 document, published in its protocol repository — `GetStream/protocol`, `openapi/chat-openapi-clientside.yaml`, 76 paths, 348 schemas, version `v237.27.1` — and struck live against `chat.stream-io-api.com` on 2026-09-06. Every live observation was made with no credential and then with a deliberately invalid one, which is the whole of what a public probe can see here.

## What this Recipe found

**The refusal names the internal handler, so the same failure has a different sentence on every route.**

Three routes, no credential, struck live:

```
GET /app                      "GetApp failed with error: \"api_key or app_id not provided\""
GET /members                  "QueryMembers failed with error: \"api_key or app_id not provided\""
GET /messages/{id}/reactions  "GetReactions failed with error: \"api_key or app_id not provided\""
```

Those prefixes are the document's own `operationId`s. So the message on a 401 is a function of the route, and a client matching on it has to know all 76. The sentence is a Go error chain printed onto the wire — an outer wrap, a colon, and the inner error in escaped quotes — which is why both routes in this Recipe declare their own refusals rather than sharing one.

**A wrong credential and a missing one differ only inside the quotes.** Both are 401, both carry `code: 2`, and the only thing that moves is `api_key or app_id not provided` → `api_key not valid`. A client separating "I sent nothing" from "I sent the wrong thing" has to parse the sentence.

**`StatusCode` is PascalCase and everything beside it is not.** Struck live:

```json
{"StatusCode":401,"code":2,"details":[],"duration":"0",
 "message":"GetReactions failed with error: \"api_key not valid\"",
 "more_info":"https://getstream.io/chat/docs/api_errors_response"}
```

`code`, `details`, `duration`, `message`, `more_info` — and `StatusCode`. One field out of six, capitalised, in the object every client parses. It is capitalised in the document too, so it is the wire name rather than a slip in one handler.

**There are two numbers on every failure and neither is the other.** `StatusCode` is the HTTP status; `code` is Stream's own, and it is `2` for a credential failure and `16` for an unknown route. Small integers with no relation to the status beside them.

**A request carries three credentials in three different places.**

```yaml
security:
  - JWT: [], api_key: [], stream-auth-type: []
  - api_key: [], stream-auth-type: []
```

Both entries are ANDs. The first wants all three at once — `Authorization: <jwt>` as a header, `?api_key=` in the query, and `Stream-Auth-Type` as a second header. The second drops the JWT, which is how anonymous access is expressed.

And **`JWT` is declared `type: apiKey`**, `in: header, name: Authorization` — not `type: http, scheme: bearer`. So the token goes in the Authorization header **with no `Bearer ` prefix**, and every client library's built-in bearer helper produces a header this API refuses.

`Stream-Auth-Type` is a header whose job is to say *which* credential you are using. It is declared as a security scheme, so a generator asks the developer for a third secret, and it is not a secret at all — its value is a word.

Struck live, the `api_key` is checked first and alone: a wrong key is refused before the other two matter, which is the credential this Recipe models.

**A query parameter carries a JSON document.** Six GET operations take a parameter literally named `payload`:

```
GET /members?payload={"type":"messaging","id":"general","limit":30,"offset":0,...}
```

It is declared with OpenAPI 3's `content:` form rather than `schema:` — legal, and rare enough that most generators fall back to `payload: string` and hand the encoding back to the caller. The filter, the sort, the limit and the offset all live inside it, so **a listing's paging controls sit inside a URL-encoded JSON blob**: every proxy logs it in full, no cache can key on it sensibly, and it counts against URL length. `/users`, `/search`, `/query_banned_users`, `/query_future_channel_bans` and `/moderation/flags/message` are the others.

`/messages/{id}/reactions` — the listing this Recipe pages — is the exception, with real named `limit` and `offset` parameters.

**`sort` is an array with `maxItems: 1`.** An array shape that promises composable ordering and a constraint that removes it. Its entries are `{field, direction}` with direction as the integers `1` and `-1`.

**The real filter grammar lives in vendor extensions.** `filter_conditions` is declared `type: object` with `additionalProperties: {}` — which says "anything at all" — while a sibling `x-stream-filter-fields` names each filterable field and the operators it accepts, and `x-stream-sort-fields` names the sortable ones. The constraints exist and are machine-readable, in keys no generator reads. So every generated client accepts an object the server will reject, and the document had the answer the whole time.

`x-stream-index` appears **2548 times**, a dotted ordering string on almost every field in the document.

**`duration` is required on every response and is the string `"0"` on every failure.** 85 response schemas declare it, all `type: string`, described as "Duration of the request in milliseconds" — and **not one of the 85 declares a format**. The document uses `format: duration` exactly once, and not on any of them. Being required everywhere is the one thing this field gets right: a client that reads it never has to check whether it is there.

**A wrong method is a 404.** Struck live, `DELETE /app` answers `{"StatusCode":404,"code":16,…,"message":"Request failed with error: \"Not Found\""}` — byte for byte what an unknown path gets. The router matches the pair and reports only that the pair is unknown, so a client cannot tell a typo in a URL from a verb it should not have used.

**One `deprecated: true` in 348 schemas, and it is not on the deprecated field a caller will meet.** The single marker sits on `UpdateUsersResponse.membership_deletion_task_id`, described as "Deprecated: always empty". Meanwhile `ChannelMemberResponse.role` — the field a client reads to find out what someone is allowed to do — says "(DEPRECATED: use `channel_role` instead)" in its description and carries no marker at all.

That is the exact inverse of [firebaseauth](../firebaseauth), one Recipe away, which marks 49 properties deprecated and describes none of them. Two documents, two vocabularies, and a tool that reads only the machine-readable half is wrong about both — in opposite directions.

**20 enumerations are prose.** "One of: member, moderator, admin, owner", written into a description, against 37 real `enum:` keys elsewhere in the same document. A third of this API's vocabularies are invisible to a generator.

**A channel member has no id.** `ChannelMemberResponse`'s required list is `custom`, `created_at`, `updated_at`, `banned`, `shadow_banned`, `channel_role`, `notifications_muted` — and there is no member identifier anywhere on it. A member *is* the pair (channel, user), so `user_id` is the only handle there is.

Five separate fields describe bans on that record: `banned`, `shadow_banned`, `ban_expires`, `ban_from_future_channels` and `future_channel_ban_expires`. A moderation screen showing only `banned` shows a shadow-banned member as fine.

**And the version is `v237.27.1`.** Major version two hundred and thirty-seven, on a document whose paths carry no version at all. Nothing in a URL changes when it moves, so a client has no way to pin one and no way to notice it moved.

## Detection

Three clients, each checked by fetching its published archive and grepping for the host this Recipe serves:

| Package | Ecosystem | Names `chat.stream-io-api.com` |
|---|---|---|
| `stream-chat` | npm | yes, 7 times |
| `get-stream/stream-chat` | Packagist | yes |
| `github.com/GetStream/stream-chat-go/v8` | Go | yes |

**`getstream` — the package named after the company — is left out, and it is the sharpest near-miss here.** It is official, it is maintained, and its own description reads "the official low-level GetStream.io client". It names `api.stream-io-api.com` four times and `chat.stream-io-api.com` **zero** times, because it is the client for Activity Feeds: a different product of the same company, on a different host, with a different API. A project depending on it is not calling anything this Recipe serves.

`stream-chat-react` is left out for the ordinary reason: it names neither host, because it reaches this API through `stream-chat` rather than itself. The [firebaseauth](../firebaseauth) mapping applies the same test in the other direction — `firebase` is a meta-package and is mapped, because its archive does name the host.

## Modelling limits

- **Only the `api_key` is checked.** The document requires three credentials at once and the wire refuses on the query parameter before the other two are looked at, so that is the one this Recipe can model faithfully. Serving the JWT and `Stream-Auth-Type` as well would make the fake stricter than the provider, and this format has no way to require several distinct credentials together — `also:` is for one credential in several places, which is a different thing.
- **`/members` is not paged.** Its `limit` and `offset` live inside the `payload` JSON, which this format matches from named query parameters only. Declaring page-style paging here would put the controls somewhere the real API does not read them. The route is served unpaged and the finding is recorded rather than approximated.
- **The `payload` parameter is not parsed.** Filters, sorts and the `x-stream-filter-fields` operator grammar are recorded above and not routed.
- **Two routes of 76.** Channels, messages, threads, polls, moderation, devices, blocklists and the long-poll transport each want their own evidence. What is here is the reaction listing and the member listing, and the failure path they share.
- **A Go major bump changes the dependency string.** Detection maps `.../stream-chat-go/v8`; a project pinned to `v7` declares a different module path and is not matched. That is how Go module majors work rather than anything Stream did.
