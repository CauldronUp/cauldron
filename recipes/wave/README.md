# wave

Emulates the Wave Accounting GraphQL API for local development and tests.

**12 conformance cases, 7 checked against the live API on 2026-09-07.**

Written against Wave's GraphQL reference at `developer.waveapps.com` and struck live against `gql.waveapps.com` on 2026-09-07 with no credential, with a deliberately invalid one, with a query that inlines an argument, and on a path that does not exist.

## What this Recipe found

**The code is right and the sentence is wrong twice.**

```
POST /graphql/public   { businesses { edges { node { id name } } } }   (no credential)
200 OK
{"errors":[{"extensions":{"code":"UNAUTHENTICATED",
                          "id":"803193d3-38e9-477d-860b-680e8d91eea7"},
            "message":"Invalid request, authentication expired.",
            "locations":[{"line":1,"column":3}],"path":["businesses"]}],
 "data":{"businesses":null}}
```

`UNAUTHENTICATED` is exactly the right code, which — after [railway](../railway) answering `INTERNAL_SERVER_ERROR` to the same question — is worth saying out loud.

The sentence beside it makes two claims and neither is true. The request was not invalid: it was a well-formed document. And nothing expired, because nothing was sent. A caller reading it goes looking for a token that has timed out when what they have is no token at all.

[adobesign](../adobesign) writes the honest version of the same sentence — "Access token provided is invalid **or has expired**" — and says it about a request that actually carried a token.

**`data` is `{"businesses": null}` rather than `null`.** The failing field is nulled and the document around it is not. [railway](../railway) nulls the whole document for the same failure at the same status, so `data.projects` is `undefined` there and `data.businesses` is `null` here.

Both are permitted by the GraphQL specification, which is what makes the difference expensive: **neither is a bug anyone can report**, and the same client code reads two different things depending on which vendor it is pointed at.

**Inline arguments are refused outright.**

```
POST /graphql/public   { __type(name: "Query") { fields { name } } }
200 OK
"Inline argument of type String is not allowed. Provide data using GraphQL variables."
```

Every literal in a query document is rejected. So every example anybody copies out of a GraphQL tutorial fails against this API, and introspection cannot be written the way the specification's own examples write it — the query has to be reshaped into `query Q($n: String!) { __type(name: $n) … }` before the server will look at it.

The policy is defensible: it makes injection through string interpolation impossible and makes query documents cacheable by hash. **No error message says any of that.** What a caller learns is that their query is not allowed.

And that refusal arrives as `GRAPHQL_VALIDATION_FAILED` at **HTTP 200**, where Railway sends the same code with a 400. So the two GraphQL APIs in this batch disagree about the status of a validation failure as well as about the shape of `data`.

**Every error carries a fresh UUID inside `extensions`.** Not beside the error, and not in a header: `extensions.id`, new on every request. It is the thing to quote to support, and it is in the last place a client would look for it — nested under the object whose documented purpose is vendor-specific metadata about the error, alongside the code.

**The Relay connection pages by number.** `businesses` returns `edges` and `node`, which is the cursor-connection shape, and its `pageInfo` is an `OffsetPageInfo`: `currentPage`, `totalPages`, `totalCount`. There is no cursor anywhere. So it looks like a cursor connection, is paged with `page` and `pageSize`, and a client written against the shape rather than the schema will look for `endCursor` and not find one.

`currentPage` is non-null in the schema and `totalPages` and `totalCount` are not — so a caller always knows where it is and may never learn how far the listing goes.

**An unrouted path is Express's own error page.** `<pre>Cannot POST /graphql/cauldron-nope</pre>` under `text/html`, so the framework beneath the GraphQL server answers when the GraphQL server does not, and the third failure on this host is neither of the first two shapes.

**The currency is an object rather than a code.** Reading it is two hops, and a query that names `currency` without naming a subfield is not a valid document at all — so the field cannot be fetched by accident.

## Modelling limits

- **One route.** Businesses. Customers, invoices, products, accounts, transactions and the whole money-in surface each want their own evidence.
- **No GraphQL is parsed.** A document naming `businesses` gets the businesses fixture whatever fields it asked for, and the inline-argument refusal is served for documents naming `__type` rather than by detecting literals.
- **Nothing is mapped in detection.** Wave is called through generic GraphQL clients and a URL, which no dependency name distinguishes.
- **No `spec:`.** The API is GraphQL, so its description is a schema.
