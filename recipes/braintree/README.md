# braintree

Emulates the Braintree GraphQL API for local development and tests.

**12 conformance cases, 6 checked against the live API on 2026-09-09.**

Written against Braintree's GraphQL reference at `graphql.braintreepayments.com` and struck live against `payments.sandbox.braintree-api.com` on 2026-09-09 with no credential, with a deliberately invalid one, with a query naming a field that does not exist, and with the required version header left off.

## What this Recipe found

**`extensions` appears twice in one response and means two different things.**

```json
{"extensions":{"requestId":"45d995dd-84a3-40ca-ae78-d67cc7d02c43"},
 "data":null,
 "errors":[{"message":"Authentication credentials are missing. Authorization header is required and must contain a value.",
            "extensions":{"errorClass":"AUTHENTICATION","errorType":"developer_error"}}]}
```

The top-level `extensions` carries a request id. The `extensions` inside the error carries a taxonomy. Same key, two levels, and nothing in common between what they hold — so a client reading `body.extensions` and a client reading `body.errors[0].extensions` are reading unrelated things, and code that merges them loses one.

GraphQL's specification defines both — a top-level `extensions` for the response and an `extensions` on each error — so the collision is the specification's rather than Braintree's. It is still the place a reader gets lost, and it is the first response in this collection to carry both at once.

**There are two taxonomies inside the inner one.** `errorClass` is `AUTHENTICATION` and `errorType` is `developer_error`. One says what went wrong, the other says whose fault it is, both about the same failure, in two casings — screaming snake beside lower snake, in one object.

**Authentication runs before the query is parsed.** `{ nope }`, naming a field the schema does not have, answers the authentication error rather than a validation error. So an unauthenticated caller cannot tell a typo from a missing credential, and cannot learn anything about the schema by asking it questions.

[railway](../railway) is the exact opposite: its whole schema is introspectable with no credential at all — `apiTokens`, `auditLogs`, `bucketS3Credentials` and the admin fields included — and its authorisation failure is reported as `INTERNAL_SERVER_ERROR`. Braintree tells you nothing and is accurate about why; Railway tells you everything and is wrong about why.

**The version header is invisible until you authenticate too.** Braintree requires `Braintree-Version` on every request. Leaving it off answers the same authentication error, so the requirement cannot be discovered by omitting it.

**Both sentences are true, and the first names the header.** "Authentication credentials are missing. Authorization header is required and must contain a value." for a request carrying none, and "Authentication credentials are invalid." for one carrying a wrong one. Two branches, both accurate, and the missing case says where the credential goes — which seven providers in this collection do not manage.

**`edges` is nullable and `pageInfo` is not.** `TransactionConnection` types them `[TransactionConnectionEdge]` and `PageInfo!`, so a conforming response must carry the paging metadata and need not carry any transactions.

It is the same trade [middesk](../middesk) makes in REST, where `object` and `has_more` are required and `data` is not. Two providers, two protocols, one decision: **the envelope is guaranteed and the contents are optional.**

**Every transaction has two identifiers.** `id: ID!` is the GraphQL global one, `legacyId: ID!` is the old REST one. Both non-null, neither interchangeable, and only the word "legacy" distinguishes them.

**Money carries its currency twice.** `MonetaryAmount` is `{value, currencyIsoCode, currencyCode}` — two fields of the same type whose descriptions say the same thing in different words: "The ISO code for the money's currency" and "The currency code for the monetary amount".

And `value` is documented as "either a whole number or a number with up to 3 decimal places". So a money field in a payments API has **three** decimal places and two permitted shapes, and code that assumes two truncates a real value.

**The listing cannot be asked for everything.** `transactions` takes `input: TransactionSearchInput!` — non-null — so there is no form of this query meaning "all of them". A search filter is mandatory before anything comes back.

**Two more the schema states and this Recipe does not serve.** `paymentMethodSnapshot` is typed nullable and its description says "This will always be present", so the prose and the type disagree. And `paymentMethod` is "only present if a multi-use payment method was used to create the transaction **and it has not been deleted**" — so deleting a card retroactively changes the shape of transactions that already happened.

## Modelling limits

- **One route.** Transaction search. Customers, disputes, verifications, payments, the whole tokenisation and mutation surface each want their own evidence.
- **No GraphQL is parsed.** A document naming `transactions` gets the transactions fixture whatever fields it asked for, and paging is read from `variables.first` and `variables.after` rather than from arguments inlined in the document.
- **The mandatory `input` is not enforced.** The schema makes `TransactionSearchInput!` non-null; nothing here refuses a query that omits it, because nothing here parses the query.
- **No `spec:`.** The API is GraphQL, and this schema cannot be introspected without a credential — authentication runs before the document is parsed.
- **Nothing is mapped in detection.** Braintree's published SDKs target its older REST and server-to-server APIs rather than this GraphQL surface, so a dependency on one says nothing about whether a project speaks GraphQL to `payments.braintree-api.com`.
