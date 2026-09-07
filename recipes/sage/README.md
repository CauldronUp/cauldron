# sage

Emulates the Sage Accounting API for local development and tests.

**10 conformance cases, 5 checked against the live API on 2026-09-07.**

Written against Sage's reference at `developer.sage.com` and struck live against `api.accounting.sage.com` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**Two failures, two shapes, and one of them is a bare array with dollar-prefixed keys.**

```
GET /v3.1/contacts    (no credential)
401 { "error": "Authorization missing. Please, add an authorization token at the 'Authorization' HTTP header."}

GET /v3.1/contacts    Authorization: Bearer notreal
401 [{"$severity":"error","$dataCode":"AccessTokenVerificationError",
      "$message":"Access token is invalid (JWT malformed or missing).","$source":"Authorization"}]
```

The first is an object with one field. The second is a **top-level array** whose every key begins with a dollar sign.

A bare array is the shape [salesforce](../salesforce) has: `body.error` is `undefined`, and a client has to index before it can read anything. The `$` prefixes are rarer still — legal JSON, and awkward in every language that maps JSON onto objects, because `$severity` is not an identifier and has to be reached by string.

**So one API answers two documents that share no field, no shape and no naming convention, on one route, at one status.**

**`$source` names the header** — `"Authorization"` — which is the useful thing, and it is on the shape that is hardest to read.

**The first message has a stray comma and a leading space.** `{ "error":` with a space after the brace, and "Please, add an authorization token" with a comma after *Please*. Small, and exactly the kind of thing that survives into every fixture anybody records by hand.

**An unknown path answers the token failure**, so routing is judged after the credential.

**The listing envelope is dollar-prefixed too** — `$items`, `$total`, `$page` — so the one thing the success and failure shapes share is a sigil. And the request parameters are snake case (`items_per_page`) while the response fields are dollar-prefixed camel case (`$itemsPerPage`), so the name a client sends and the name it reads back are spelled two different ways.

**The display string folds the reference into the name.** `"Northwind Traders (NORTH01)"` when there is a reference and the bare name when there is not — so a client parsing `displayed_as` to recover the reference cannot tell an absent one from a parse failure.

## Modelling limits

- **One route.** Contacts. Invoices, credit notes, bank accounts, ledger entries, tax rates and the whole journal surface each want their own evidence.
- **Nothing is mapped in detection.** Sage publishes SDKs per product and per country, and none resolves under a name that identifies this surface.
- **No `spec:`.** Sage publishes a rendered reference site.
