# gitea

Emulates the Gitea API (v1), for local development and tests.

**11 conformance cases, 8 checked against the live API on 2026-09-06.**

The live cases were struck against `gitea.com` using public repositories and no account. Written against Gitea's own generated swagger — `gitea.com/swagger.v1.json`, 342 paths, version 1.27.0+dev.

Gitea is a self-hosted git forge whose API is deliberately GitHub-shaped: owner and repo in the path, an issues listing that includes pull requests, a `Link` header advertising the next page. That resemblance is the danger. The places it diverges are places a GitHub client has no reason to look.

## What this Recipe found

**`per_page` is not the parameter, and the `Link` header echoes it back anyway.** The page size is `limit`. Sending GitHub's `per_page` does not fail — it is ignored and the default of 10 is served. That much is an ordinary naming difference. Struck live against `gitea/tea`, a repository with 149 issues:

```
GET /api/v1/repos/gitea/tea/issues?per_page=1
HTTP 200
X-Total-Count: 149
Link: <...issues?page=2&per_page=1>; rel="next",
      <...issues?page=15&per_page=1>; rel="last"
(10 records)

GET /api/v1/repos/gitea/tea/issues?limit=1
HTTP 200
X-Total-Count: 149
Link: <...issues?limit=1&page=2>; rel="next",
      <...issues?limit=1&page=149>; rel="last"
(1 record)
```

The first response ignored `per_page`, served pages of 10, computed `rel="last"` as page 15 — and then wrote `per_page=1` into the very URLs it handed back. The pagination header confirms a belief the server did not act on. A client trusting `rel="last"` to say how many pages there are is reading a number derived from a page size it never received.

`X-Total-Count` is the only quantity in that response that is true under both spellings.

**The limit is clamped in silence at both ends.** `limit=9999` answers 50 records and no warning, so a client asking for everything at once gets 50 and a `Link` header — and one that reads the record count as "that was all of them" stops at 50 of 149.

**The issues listing returns pull requests, and links to them as pulls.** GitHub does this too, so a ported client is already wrong the same way — but Gitea's `pull_request` is *present and null* on a plain issue rather than absent, so the GitHub test for it (`"pull_request" in issue`) is true for everything. Struck live: both records on the first page of `gitea/tea` were pull requests, and each `html_url` pointed at `/pulls/1112`, not `/issues/1112`. The number space is shared, and `/issues/1112` fetches the pull request.

**A record has two identifiers and the path takes the small one.** `number` is per repository, `id` is instance-global. Struck live: issue number 1112 has id 495221, and the path takes 1112.

**The credential is prefixed `token `, and there is a documented way to put it in the URL.** The swagger's own `AuthorizationHeaderToken` says tokens must be prepended with `token` and a space — not `Bearer`. Beside it sits `AccessToken`, an apiKey `in: query` named `access_token`, described there as deprecated for removal in Gitea 1.23. A credential in a query string reaches access logs, proxy logs and `Referer` headers.

**Failures point at the docs rather than at the request.** Struck live:

```
GET /api/v1/repos/gitea/no-such-repo-xyz
HTTP 404  {"message":"not found","url":"https://gitea.com/api/swagger"}

GET /api/v1/user   (Authorization: Bearer bogus)
HTTP 401  {"message":"invalid username, password or token",
           "url":"https://gitea.com/api/swagger"}
```

`url` is the same swagger link on every failure — not the failing request, not documentation for this error, and not unique. An error reporter grouping by it has one group. The 401 carries no `WWW-Authenticate`.

**Nothing says how much budget is left.** No RateLimit headers and no `Retry-After` on any response struck. The only way to find an instance's limit is to cross it.

## Modelling limits

- **Issues on one repository, listed and fetched.** 342 paths is a forge: pulls, releases, orgs, admin, packages and Actions each want evidence of their own.
- **The credential is checked for presence, not validity**, and `token ` is served as the scheme rather than verified. Basic auth and the `Sudo` impersonation header are documented and not modelled. Distinguishing a *rejected* query credential from no credential at all needs a valid token, which this Recipe does not have and did not ask for — so `access_token` is recorded from the spec rather than struck.
- **`pull_request` is served as an object or null and its contents are not modelled.** Presence, not shape, is the claim.
- **The `Link` header here is computed from the page size actually served.** Reproducing the echo of an ignored parameter would mean serving a URL this emulator knows to be misleading; the finding is recorded in words instead.
