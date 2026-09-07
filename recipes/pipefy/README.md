# pipefy

Emulates the Pipefy GraphQL API for local development and tests.

**10 conformance cases, 5 checked against the live API on 2026-09-07.**

Written against Pipefy's reference at `developers.pipefy.com` and struck live against `api.pipefy.com` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**Two specifications, one endpoint, and which you get depends on which mistake you made.**

```
POST /graphql    (no credential)
401 {"errors":[{"title":"Unauthorized","detail":"You are not authorized to access this page"}]}

POST /graphql    Authorization: Bearer notreal
401 {"error":"invalid_token","error_description":"The access token is invalid","state":"unauthorized"}
```

The first is **JSON:API** — an `errors` array of objects with `title` and `detail`. The second is **OAuth 2.0's bearer-token error response** from RFC 6750: `error`, `error_description`, and a `state` that RFC 6749 defines for the authorization *request* rather than for an error.

One route, two documents, two unrelated standards bodies. A client cannot write one parser for a 401 from it. [triggerdev](../triggerdev) has the same split and at least keeps both halves in-house.

**The OAuth half is missing its header.** RFC 6750 puts a bearer-token error in `WWW-Authenticate` and makes the body optional. Pipefy sends the body and no header — so the machine-readable channel the specification defines is empty and the optional one carries everything.

**"You are not authorized to access this page."** A *page* — on a GraphQL endpoint that has never rendered one. Rails' own phrasing for a failed `authorize` filter, reaching an API client through shared middleware. The same shape [grafanacloud](../grafanacloud)'s "Login Required" has.

**A GraphQL API that answers 401.** GraphQL's own convention is 200 with an `errors` array, because a response can be partly successful. Pipefy uses HTTP statuses for authentication and the GraphQL envelope for everything else — the pragmatic choice, and it means a client has three shapes to handle rather than one.

**A done card can still be overdue.** `done: true` beside a due date that has passed, so a board counting overdue cards checks two fields and only one is called `done`.

**The phase is a name rather than an identifier.** Renaming a phase changes every record that reports it, and two pipes with a phase called "Done" are indistinguishable.

## Modelling limits

- **One query.** Cards. Pipes, phases, tables, records, organizations, webhooks and the whole automation surface each want their own evidence — and in GraphQL they are all the same route, so what varies is the query rather than the path.
- **The response key is the caller's choice.** GraphQL names each result after the field in the query, so `data.cards` is what *this* query produces and a different query produces a different envelope. Served as the documented field name.
- **Nothing is mapped in detection.** Pipefy is called through a generic GraphQL client and a URL, which no dependency name distinguishes.
- **No `spec:`.** The schema is introspectable rather than published, and introspection needs a credential.
