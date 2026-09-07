# aircall

Emulates the Aircall API for local development and tests.

**12 conformance cases, 5 checked against the live API on 2026-09-07.**

Written against Aircall's reference at `developers.aircall.io` and struck live against `api.aircall.io` on 2026-09-07 with no credential, with a deliberately invalid one, and on paths that do and do not exist.

## What this Recipe found

**Three statuses for three states, and every one of them is right.**

```
GET /v1/calls           (no credential)      401 {"message":"Unauthorized"}
GET /v1/calls           Basic notreal        403 {"message":"Forbidden"}
GET /v1/cauldron-nope   Basic notreal        404 {"message":"Not Found"}
```

401 for "I do not know who you are", 403 for "I know and the answer is no", 404 for "there is nothing here". That is what RFC 9110 says those three statuses mean, and Aircall is one of the very few providers in this collection that sends all three correctly from one endpoint.

Most of the collection sends one status for everything — [ironclad](../ironclad), [bird](../bird) and [acuity](../acuity) each answer three distinct mistakes with one byte-identical body — and several send the status that describes a different failure entirely.

**The body is the status line, copied.** `{"message":"Unauthorized"}`, `{"message":"Forbidden"}`, `{"message":"Not Found"}`: each is the reason phrase for the status it arrived with. There is no code, no correlation id and no link anywhere in any of them, so a caller who reads the body learns exactly what the first line already said.

Taken together those two findings are one trade, and it is the opposite of the trade [acuity](../acuity) makes. Acuity puts the status in the body three times, in three renderings, and answers everything with 401. Aircall puts the status in the body once and gets the status right. **Aircall's caller can act on the response and Acuity's cannot**, and the richer-looking body is the useless one.

**Routing runs before the credential**, which is what makes that 404 possible.

```
GET /v1/ping     (no credential)   401
GET /v1/         (no credential)   404
```

The two are separated by whether the path exists rather than by whether a key was sent, so a caller with a typo in the URL is told about the typo instead of being told to check its key. That is the failure four providers in this collection cannot report at all.

The cost is on the other side of the same coin: **the whole surface is enumerable anonymously.** Anyone can learn which paths this API serves without holding an account, one request at a time, because 401 and 404 partition them.

**A wrong scheme is a 403 rather than a 401.** `Authorization: Bearer notreal` on an API that takes Basic answers "Forbidden" — so a credential that could not authenticate anybody is reported as a caller who is known and refused. It is the one mistake this otherwise precise status set could have named exactly, and the one it rounds off.

**The timestamps are bare integers and one of them is a duration.** `started_at`, `answered_at` and `ended_at` are seconds since the epoch; `duration` is seconds of elapsed time. Four integer fields in the same unit, two of them meaning an instant and one meaning a length, and no name says which.

**A missed call has no `answered_at` at all.** The field is absent rather than null, so the field that says when the conversation began is the one a missed call does not have, and a client subtracting it from `ended_at` to get the ring time gets `NaN`. The reason lives in a second field, `missed_call_reason`, which is likewise absent on a call that was answered.

**The next page is a whole URL and the previous one is null.** `meta.next_page_link` is an absolute address, so a client concatenating it onto a base URL builds one that does not exist — the mistake eight Recipes in this collection had to be corrected for. `meta.previous_page_link` is present and null on the first page rather than absent, so a client testing for the key to detect the first page never detects it.

**Two counts one word apart.** `meta.count` is how many are on this page and `meta.total` is how many matched. The one called `count` is the smaller.

## Modelling limits

- **One route.** Calls. Users, numbers, contacts, teams, tags and the whole webhook surface each want their own evidence.
- **Nothing is mapped in detection.** Aircall publishes no client of its own on npm, Packagist or the Go module proxy under an obvious name, checked 2026-09-07.
- **No `spec:`.** Aircall publishes a rendered reference site.
