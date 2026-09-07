# payu

Emulates the PayU orders API for local development and tests.

**9 conformance cases, 3 checked against the live API on 2026-09-07.**

Written against PayU's reference at `developers.payu.com` and struck live against `secure.payu.com` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**Four fields for one failure, and three of them say the same thing.**

```json
{"status":{"statusCode":"UNAUTHORIZED",
           "code":"401",
           "codeLiteral":"INVALID_OR_MISSING_ACCESS_TOKEN",
           "statusDesc":"Incorrect authentication"}}
```

`statusCode` is a word. `code` is the HTTP status **as a string**. `codeLiteral` is a screaming-snake constant. `statusDesc` is prose. Three vocabularies for one fact plus a restatement of the status line — inside an object called `status` whose `statusCode` is not a code and whose `code` is not a status.

The one field carrying information no other has is `codeLiteral`: `INVALID_OR_MISSING_ACCESS_TOKEN` names **both** cases, which is the same honest ambiguity [pulumi](../pulumi) writes as a sentence. Everything else is derivable from the status line.

**A missing credential and a wrong one are byte-identical**, so the field naming both cases is the accurate one and nothing separates them.

**An unknown path is a 302.** No body, no JSON, a redirect to PayU's merchant panel login. A client following redirects by default — which most do — ends up with an HTML sign-in page carrying a **200**, and has to work out that its URL was wrong from the content type.

That is the worst not-found in this collection so far. [openphone](../openphone) and [triggerdev](../triggerdev) answer HTML with a 404, which at least keeps the status honest. Here the status a client finally sees is a success.

**The `status` key carries the success too**, which is the good half: one field to check on both paths rather than a success shape and a failure shape that share nothing.

**Money is minor units in a string.** `"21000"` is 210 złoty — exact, because it is text; a hundred times what a person would say, because it is minor units; and the only thing that says which is the currency code beside it.

**There are three ways of not being finished.** `NEW`, `PENDING` and `WAITING_FOR_CONFIRMATION` are all "not yet" and only the last means somebody has to act, so a client testing `status !== COMPLETED` cannot tell which it has.

## Modelling limits

- **One route.** Orders. Refunds, payouts, payment methods, tokenised cards and the whole notification surface each want their own evidence.
- **The 302 is recorded, not served.** Reproducing a redirect to a login page would mean serving a sign-in flow, and the finding is that the redirect exists.
- **Nothing is mapped in detection.** No client for this API turned up on npm, Packagist or the Go module proxy under an obvious name on 2026-09-07.
- **No `spec:`.** PayU publishes a rendered reference site.
