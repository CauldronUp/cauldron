# affirm

Emulates the Affirm payments API for local development and tests.

**13 conformance cases, 6 checked against the live API on 2026-09-07.**

Written against Affirm's reference at `docs.affirm.com` and struck live against `api.affirm.com` on 2026-09-07 with no credential, with a deliberately invalid one, with and without the required query parameter, and on a path that does not exist.

## What this Recipe found

**The parameter check runs before the credential check.** Struck live with no `Authorization` header at all:

```
GET /api/v1/transactions
400 application/problem+json
{
  "detail": "Missing query parameter 'transaction_type'",
  "status": 400,
  "title": "Bad Request",
  "type": "about:blank"
}
```

So the first thing an unauthenticated caller learns from this API is the name of a query parameter it did not send. Add the parameter and the same credential-less request becomes a 401. Send a wrong credential without the parameter and it is the same 400, byte for byte.

**And the parameter it is telling you about has exactly one legal value.** Affirm's own description of `transaction_type` is `string enum required`, `Allowed: charge` — one member. It is mandatory, it can only be sent one way, and it therefore distinguishes nothing except a request that remembered it from a request that did not. Which is the request above, answered before authentication.

**One endpoint, two error models, and `type` means something different in each.**

```
400  {"detail": …,           "status": 400,       "title": "Bad Request", "type": "about:blank"}
401  {"status_code": 401,    "type": "unauthorized",
      "code": "api-key-pair-not-found",
      "message": "You did not provide an API key pair."}
```

The 400 is RFC 9457, and `type` is a URI. The 401 four bytes away has a `type` that is a bare word, a `code` that is a slug, `status_code` where the other has `status`, and `message` where the other has `detail`. **Not one field name is shared.** A client that wants the reason out of a failed request to this one endpoint needs two parsers and a branch on the status to choose between them.

**`about:blank` is the specification's own value for "nothing beyond the status code".** It is the fifth problem document in this collection and the only one that declares itself empty. [postman](../postman) points `type` at a URI naming a Postman problem, [triggerdev](../triggerdev) at MDN's page about the status, [agicap](../agicap) at the RFC clause defining 500, and [revai](../revai) omits it. So the five uses are: right, a documentation link, a specification link, absent, and explicitly nothing.

**The two credential failures are told apart, and both sentences are true.** `api-key-pair-not-found` with "You did not provide an API key pair." for a request carrying none, and `api-key-pair-invalid` with "You have provided an invalid API key pair." for one carrying a wrong one. After seven providers here that answer "invalid key" to a request holding no key, this is one of the few that gets both halves right — and it names the credential's actual shape, a pair rather than a token.

**An unknown path is a styled web page.** `<title>404 - Wrong Page</title>`, with a navbar and a webfont, under `text/html`, on the API host. So one API answers three content types across three failures: `application/problem+json`, `application/json`, and a marketing 404.

**The listing says more exists and gives you nothing to ask with.** The schema behind the reference declares `ListTransactions` as `has_prev`, `has_next` and `transactions` — two booleans and an array, and **no cursor anywhere**. Paging means taking the last record's own identifier and sending it back as `before`, so the envelope cannot be paged from; the page has to be read. The sibling endpoint, `ListTransactionEvents`, sends `prev_page` and `next_page` as URL-encoded parameter strings, so two listings on one API page two different ways.

The direction words invert too: the schema glosses `has_prev` as "previous (**newer**) results" and `has_next` as "next (**older**) results", so `next` means back in time.

**The reference and the schema disagree about a field's type.** The attribute table on the Transaction Object page calls `provider_id` a `string`; the schema behind the endpoint page calls it an `integer`; the wire sends a number. It is a two-value enum — 1 is Affirm, 2 is Katapult — so a merchant's listing spans two lenders and the field that says which is a bare integer with no label on the record.

**The example JSON on the reference page does not parse.** `"id: "A1B2-C3D4",` — the key is missing its closing quote. And the identifier in it is nine characters with a hyphen, where the prose two paragraphs above says an id is "either a 12 or 16 alphanumeric character string".

**The schema documents one of its own fields as unreliable.** `fee` on a transaction event: "The correctness of this value may be inaccurate for certain loan types, refer to the settlement reports for the accurate merchant fee information for accounting and reconciliation purposes." A field in a payments API, in the machine-readable description, telling you not to use it for accounting.

## Modelling limits

- **One route.** Listing transactions. Authorize, capture, refund, void, update, settlement events, the whole Checkout, Cards and Disputes surfaces each want their own evidence.
- **`has_prev` is recorded rather than served.** This format has one boolean for "more exists" and Affirm sends two, one in each direction. Serving a constant for the other would be wrong on every page but the first.
- **No `spec:`.** There is an OpenAPI document and it is embedded inside the rendered reference page rather than served at an address, so there is nothing stable to record a hash of.
- **Nothing is mapped in detection.** The packages naming Affirm are merchant-platform plugins rather than clients of this API.
