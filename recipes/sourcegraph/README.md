# sourcegraph

Emulates the Sourcegraph GraphQL API for local development and tests.

**10 conformance cases, 7 checked against the live API on 2026-09-10.**

Struck live against `sourcegraph.com` on 2026-09-10 with no credential, with a deliberately invalid one, with a query naming a field that does not exist, and on a path that does not exist.

## What this Recipe found

**Not being logged in is data rather than an error.**

```
POST /.api/graphql   { currentUser { username } }   (no credential)
200 {"data":{"currentUser":null}}
```

No `errors` array at all. An anonymous caller is a valid session whose current user happens to be nobody — and a client checking `body.errors` before reading `body.data` finds nothing to check.

**And a wrong credential is a 401 in plain text.**

```
Authorization: token notreal
401 Content-Type: text/plain
Invalid access token.
```

So sending *nothing* succeeds and sending something *wrong* fails at the transport — the opposite way round from every other GraphQL API in this collection.

**That completes a set of three, and no two agree:**

| | no credential | wrong credential | field that does not exist |
| --- | --- | --- | --- |
| [railway](../railway) | 200, `data: null`, errors, code `INTERNAL_SERVER_ERROR` | same as none | **400**, `GRAPHQL_VALIDATION_FAILED` |
| [wave](../wave) | 200, field nulled, errors, code `UNAUTHENTICATED` | same as none | 200, `GRAPHQL_VALIDATION_FAILED` |
| sourcegraph | 200, field nulled, **no errors** | **401 `text/plain`** | 200, errors with **no code at all** |

Three GraphQL APIs, three different answers to each of three questions, and every one of them permitted by the specification. A client written against one reads nonsense from the next, and nothing it reads is a bug it can report.

**A query naming a field that does not exist has no machine-readable code.**

```json
200 {"errors":[{"message":"Cannot query field \"nope\" on type \"Query\". Did you mean \"node\"?",
                "locations":[{"line":1,"column":3}]}]}
```

No `extensions`, so no `code`. The message is the entire machine-readable content of the failure, and a client wanting to tell a validation error from anything else has to match on English prose.

It does carry a spelling suggestion — *Did you mean "node"?* — which is genuinely useful to a person, and is the part a machine cannot use.

**`totalCount` is 0 on a page holding two records and claiming more.** Struck live:

```json
"nodes": [ …two repositories… ],
"pageInfo": {"hasNextPage": true, "endCursor": "UmVwb3NpdG9yeUN1cnNvcjp7…"},
"totalCount": 0
```

The listing plainly has records, says there are further pages, and reports a total of zero. A client rendering "0 results" above two rows is reading the field correctly.

Counting a global code index is expensive and this connection declines to do it — a defensible decision, reported as a **wrong number** rather than as a null. Compare [matrix](../matrix), which faces the same cost and calls its field `total_room_count_estimate`.

**The connection is `nodes`, not `edges`.** Relay's shape without the edge wrapper, so a per-record cursor has nowhere to live and paging depends entirely on `pageInfo.endCursor`.

That makes four GraphQL providers here and three connection shapes: Railway and [braintree](../braintree) use `edges`/`node` with cursor paging, [wave](../wave) uses `edges`/`node` over *offset* paging, and this one drops edges altogether.

**The opaque identifier decodes to a type and a number.** `UmVwb3NpdG9yeTo2NDI0NTk0OA==` is base64 for `Repository:64245948`. The global id carries the type name and the underlying row id in plain sight, so anyone who base64-decodes one gets two more identifiers out of it.

**Introspection is open to anonymous callers**, as it is on [railway](../railway) — the whole schema, served to a request that cannot read a user's name.

**An unknown path is two lower-case words.** `404 text/plain`, `no route`.

**And a repository with no description sends an empty string** rather than null, so a client testing for null renders an empty caption instead of falling back.

## Modelling limits

- **One listing.** Repositories. Search, users, code intelligence, batch changes, insights and the whole admin surface each want their own evidence.
- **No GraphQL is parsed.** A document naming `repositories` gets the repositories fixture whatever fields it asked for; documents naming `currentUser` and `nope` are routed to their own fixed answers so the three shapes above can all be served from one path.
- **`totalCount` is served as the constant zero**, which is what the live connection returns. It is not a modelling shortcut: the finding is that the number is wrong, and serving a correct one would erase it.
- **No `spec:`.** The API is GraphQL, and its schema is served to anybody who asks.
- **Nothing is mapped in detection.** Sourcegraph is reached through its CLI, its browser extension and generic GraphQL clients pointed at whichever instance an organisation runs, so a dependency name says neither that this API is used nor which host it lives on.
