# qonto

Emulates the Qonto business API for local development and tests.

**9 conformance cases, 5 checked against the live API on 2026-09-07.**

Written against Qonto's reference at `api-doc.qonto.com` and struck live against `thirdparty.qonto.com` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**The two credential failures have different statuses, different serialisers and different field sets.**

```
GET /v2/transactions    (no credential)
400  {
       "errors": [ { "code": "bad_request", "detail": "Authorization field missing" } ]
     }

GET /v2/transactions    Authorization: Bearer notreal
401  {"errors":[{"code":"unauthorized","detail":"Invalid credentials"}],"trace_id":"4e46fda1…"}
```

Three differences in one pair.

**The status moves from 400 to 401** — which is the right way round and almost unique in this collection. Most providers here answer one status to both; several answer 400 to neither correctly. Qonto treats "you sent nothing" as a malformed request and "you sent something wrong" as a refusal, which is what those statuses mean.

**The first is pretty-printed and the second is compact**, so two serialisers on one route.

**And the second carries a `trace_id` the first does not** — so the failure a caller is more likely to hit anonymously is the one with nothing to quote to support.

**An unknown path is JSON declared as `text/plain`.**

```
GET /v2/cauldron-nope
404  Content-Type: text/plain
{"errors":[{"code":"not_found","detail":"Not found"}]}
```

The same mislabel [wrike](../wrike) has, in the same direction: real JSON under a header that says text, so a client checking the type before parsing skips a parse it could have done. Here it is on the routes that do *not* exist rather than the ones that do, which makes it the quieter version.

**The envelope is consistent across all three.** `errors` is always an array of `{code, detail}` objects — on 400, 401 and 404 alike — so one parser reads every failure, and only the wrapper's extra fields and the content type move.

**Money arrives twice in two types.** `amount: 210.0` as a float and `amount_cents: 21000` as an integer, on the same record. The float is the one a careless client reads and the integer is the one that is exact.

**A pending transaction has left and not settled.** `emitted_at` present, `settled_at` absent — so a balance summed over settled transactions and one summed over emitted transactions disagree, and neither is wrong.

## Modelling limits

- **One route.** Transactions. Accounts, memberships, labels, attachments, beneficiaries, requests and the whole statement surface each want their own evidence.
- **`trace_id` is recorded, not served.** It sits beside the errors array rather than inside an entry, and it is on one failure and not the other — and for a list-shaped envelope this format places a named error's fields inside the entry. Serving it on both would erase the finding; serving it inside the entry would put it somewhere Qonto does not.
- **Nothing is mapped in detection.** No client for this API turned up on npm, Packagist or the Go module proxy under an obvious name on 2026-09-07.
- **No `spec:`.** Qonto publishes a rendered reference site.
