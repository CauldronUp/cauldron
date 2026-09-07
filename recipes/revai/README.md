# revai

Emulates the Rev.ai speech-to-text API for local development and tests.

**10 conformance cases, 4 checked against the live API on 2026-09-07.**

Written against Rev.ai's reference at `docs.rev.ai` and struck live against `api.rev.ai` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**A problem document with no problem in it.**

```
GET /speechtotext/v1/jobs    (no credential)
GET /speechtotext/v1/jobs    Authorization: Bearer notreal
401  Content-Type: application/problem+json
{"title":"Authorization has been denied for this request.","status":401}
```

RFC 9457 defines five members: `type`, `title`, `status`, `detail` and `instance`. **Two arrive.** `title` is a sentence and `status` restates the status line, so the document carries one fact the caller did not already have — and it is prose.

No `type` means, per the specification, that the type is `about:blank`: *"the problem has no additional semantics beyond that of the HTTP status code"*. Which is exactly true here. Stating it by omission is the specification working as designed, and it is also the most standardised way in existence of saying nothing.

**This collection now has three `problem+json` documents and they fail three different ways:**

| | `type` |
|---|---|
| [triggerdev](../triggerdev) | MDN's page about the HTTP 401 status |
| [agicap](../agicap) | the RFC clause that *defines* 500 |
| **revai** | absent |

Three providers reaching for the same specification, and not one uses that field for what it is for: naming a problem this API can have.

**A missing credential and a wrong one are identical**, so the two failures needing different fixes are one document — and the document has no field that could have separated them even in principle.

**The sentence is passive and names no actor.** "Authorization has been denied for this request." Not "your token is invalid", not "no Authorization header": something denied something. It is ASP.NET's own default for a failed authorisation filter.

**A failed job puts its reason in a second field.** `status` is one word for every kind of failure and `failure` says which — so a client switching on status cannot tell a bad download from a bad audio file, and the field that can is absent on every job that worked.

**An unfinished job has no completion time and no duration**, absent rather than null — so "finished" and "succeeded" are different questions and only one of them has a field.

## Modelling limits

- **One route.** The job listing. Transcripts, captions, custom vocabularies, the streaming surface and the language-identification and sentiment jobs each want their own evidence.
- **Nothing is mapped in detection.** Rev.ai's SDKs are named for the company and cover several products, so a dependency on one does not identify this surface.
- **No `spec:`.** Rev.ai publishes a rendered reference site.
