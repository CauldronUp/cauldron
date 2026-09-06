# canny

Emulates the Canny API (v1), for local development and tests.

**10 conformance cases, none checked against a live API.**

Written against Canny's own API reference — `developers.canny.io/api-reference` — read on 2026-09-06.

## What this Recipe found

**Reading is a POST, and the credential is part of the body.** Listing posts is `POST /api/v1/posts/list` with the key as a field called `apiKey` alongside `boardID` and `limit`. Two consequences are worth having a fake for. A client cannot set the credential once — there is no default header to hang it on, so every call site carries the secret and the one that forgets sends a well-formed request that is refused. And the secret lands in request bodies, so anything that logs them — a proxy, an APM agent, a retry queue that serialises the request — has the key, and body logging is exactly what gets switched on mid-incident by someone not thinking about credentials.

It also means a `GET` is refused. Every read is a write-shaped request, so an HTTP cache in front of it caches nothing, and a `curl` that omits `-X POST` gets a refusal that says nothing about whether the key was right.

**The paging signal is a boolean with no total beside it.** A listing answers `{"hasMore": true, "posts": […]}` — no count, no next cursor, no `Link` header. The client pages with `skip` and stops when `hasMore` goes false. The default `limit` is **ten**, small enough that code tested against a board of eight posts never pages at all and meets the second page for the first time in production.

**Two versions disagree about the envelope key.** v1 wraps its array in `posts`; v2 wraps it in `items`. Same vendor, same kind of answer, different key — and the version is a path segment rather than a header, so following a v2 example while calling a v1 path reads `undefined` off a response that arrived perfectly well.

**A post can arrive with no author.** `author` is documented as nullable and is served present-and-null, so `author.name` throws rather than reading `undefined`. `statusChangedAt` is null on a post whose status has never moved, which is what a changelog reads.

## Modelling limits

- **The failure shape is Canny's gap, not a choice this Recipe made.** The reference documents `{"error": "…"}` and does not state the status for a bad or missing key. This Recipe serves 401. If Canny answers 200 with an error inside — which some POST-only APIs do — this is wrong, and it is written down so whoever finds out has somewhere to correct.
- **Posts only.** Boards, comments, votes, users and companies each want their own evidence.
- **Embedded objects are served, their targets are not.** An embedded board here is not a board this Recipe serves.
- **v2 is not served.** Its `items` key is recorded because the disagreement is the point.
