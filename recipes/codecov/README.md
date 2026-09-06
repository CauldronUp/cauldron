# codecov

Emulates the Codecov API (v2), for local development and tests.

**13 conformance cases, 12 checked against the live API on 2026-09-06.**

Written against Codecov's own OpenAPI 3.0.3 document, served by the API itself at `api.codecov.io/api/v2/schema/` — 34 paths, version 2.0.0. Every observation below was made **with no credential at all**: coverage for public repositories is public, so the whole read surface could be checked without an account.

## What this Recipe found

**The next-page link downgrades to `http://`.** Struck live, over HTTPS:

```
GET https://api.codecov.io/api/v2/github/codecov/repos/?page_size=1
HTTP 200
{"count": 146,
 "next": "http://api.codecov.io/api/v2/github/codecov/repos/?page=2&page_size=1",
 "previous": null,
 "results": [ ... ]}
```

The request was `https`. The link handed back is `http`, and so is `previous` on every page after the first. Following it live:

```
GET http://api.codecov.io/api/v2/github/codecov/repos/?page=2&page_size=1
HTTP 301
Location: https://api.codecov.io:443/api/v2/github/codecov/repos/?page=2&page_size=1
```

The redirect is there and the transport does get fixed — *after the request has already been made in cleartext*. A client following `next` on an authenticated call sends `Authorization: Bearer <token>` over plain HTTP on every page but the first, which is exactly what `next` exists to be followed for. Nothing in the response warns about it.

The `Location` is also written with an explicit `:443`, which is unusual enough to break exact-URL comparison, request signing, and allow-lists built from the URL the client asked for.

**`totals` is null, not empty, on repositories with no coverage.** Struck live: of the first six repositories in codecov's own account, three answered `"totals": null` — not `{}`, not zeroes. So `repo.totals.coverage` throws on half a listing, and the half it throws on is the half a dashboard most wants to draw as "no data".

**Adjacent numbers in one object have different JSON types.** In the same `totals` object, `complexity` is `0.0` and `complexity_ratio` is `0`.

**Three different 404s, and one of them is GitHub's.** Struck live:

```
/github/no-such-owner-xyz/repos/  -> 404 {"detail":"No Owner matches the given query."}
/github/codecov/repos/no-such/    -> 404 {"detail":"Github API: Not Found"}
/notaservice/codecov/repos/       -> 404 {"detail":"Invalid service for Owner."}
```

The middle one is the upstream's own words passed through: a client cannot tell "this repository is not on Codecov" from "GitHub refused the lookup", which are a configuration problem and an outage.

**A rejected credential does carry `WWW-Authenticate`.** 401 with `www-authenticate: Bearer` and `{"detail":"Invalid token."}` — worth recording because most providers modelled here send a 401 with nothing for a client's retry-with-credentials path to catch. Every failure uses the same key, `detail`, holding a single string: Django REST Framework's default, consistent across all four.

**`page_size` is not capped.** Struck live: `?page_size=9999` on an account with 146 repositories returned all 146 in one response. A client passing a page size through from user input can ask for an unbounded response.

**The service enum ships a value called `to_be_deleted`.** The `{service}` path parameter enumerates bitbucket, bitbucket_server, github, github_enterprise, gitlab, gitlab_enterprise — and `to_be_deleted`. It is in the published document, so a client generating a picker or a validator from the schema offers it.

**Nearly every path ends in a slash, and six do not.** `/repos/`, `/commits/`, `/pulls/` do; `/compare/flags`, `/compare/components`, `/compare/impacted_files`, `/report/tree` and the two `{file_path}` paths do not. Django redirects the missing case, so getting it wrong is a 301 rather than a 404 — and a 301 is where headers and request bodies get dropped.

**A repository has no numeric identifier.** The name is the only key a sync can hold, so a rename upstream is indistinguishable from a deletion plus a creation.

## Modelling limits

- **Repositories on one owner, listed and fetched.** 34 paths is a coverage service: commits, pulls, branches, flags, components, the comparison endpoints and the file-level reports each want their own evidence.
- **The `http://` downgrade is recorded and not reproduced.** This emulator is served over `http` on a loopback address, so a next link cannot be a downgrade relative to the request — there is no `https` to fall from. What the case asserts is that `next` is an absolute URL carrying the page forward; the scheme finding lives here in words.
- **`totals: null` *is* served** on the records that have no coverage, and is asserted.
- **Coverage numbers are the ones the live API answered on 2026-09-06** for public repositories. They move; the shape is the claim.
