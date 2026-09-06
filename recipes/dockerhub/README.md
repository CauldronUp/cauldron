# Docker Hub

Emulates the Docker Hub API (v2), for local development and tests.

**12 conformance cases, 7 checked against the live API on 2026-08-23.**

## What this Recipe found

Hub's, which answers for a public repository without a token:
a repository nobody can see is a 404 carrying `{"message": "object not
found"}` and no code, `count` is of everything rather than of the page, and
one image wears several tags -- on one page of `library/nginx`, seventeen
digests were shared.

Three more of its cases were settled the same day. `latest` is an ordinary tag
and is routinely the older one: `library/mongo` last updated `latest` on
2026-07-23 under digest `sha256:e0ce8c35124d...`, while its `8.3.8` tag was
updated on 2026-08-18 under a different digest. Sixty-one of the hundred tags
on that page are newer than `latest`, and the same holds on `library/python`,
`library/nginx` and `library/redis`, so it is the ordinary condition rather
than a repository somebody forgot. And a tag can be marked `inactive` and
still be answered like any other: `library/registry` returns its `2.5.2` with
`tag_status: inactive`, one of fifteen in that state out of seventy, `2.5` and
`2.7.0` among them. Plain version tags, which is what makes it worth having --
somebody pinned to `registry:2.5.2` is pinned to something the registry has
already decided is a candidate for removal, and it says so in a field nobody
reads while answering 200. Fifteen are GitHub's, which answers for a public repository
without a token: its errors repeat the HTTP status in the body as a string,
an issue carries no owner and no repo, and an issue is addressed by its
number rather than by its id -- `golang/go` issue 81026 has id 5222669952,
and only one of those two works in a path. And the fourth corrected this
Recipe: it claimed the last page of a listing carries no `Link` header, and
GitHub sends one holding `rel="prev"`. A single page carries none, which is
where the wrong claim came from -- the first page and the last page are the
same page, and a single page is not a last page.

The difference is the whole reason a client reads that header. One that stops
when the header is missing never stops against GitHub, because the header is
there; it is `rel="next"` that is not. The case was teaching exactly the loop
its own comment said the header exists to prevent.

A `rel="prev"` is opt in rather than implied, because providers disagree:
Basecamp's own README describes `rel="next"` alone, so its last page really
does carry no header. Only offset and page numbering can have one at all, a
cursor being a position the caller was handed rather than a number to count
backwards from.

Asking GitHub for the rest of that Recipe's claims found one more thing it had
never modelled: `GET /issues` returns pull requests as well as issues, because
every pull request is an issue in GitHub's data model, and the only thing
telling them apart is whether a `pull_request` key is present. Not null --
absent. Twelve of the hundred open "issues" on one page of `golang/go` are
pull requests, so an open-issue count taken from that endpoint is wrong by
that much, a sync mirroring issues into a tracker imports pull requests, and
neither errors. The Recipe returned issues only, which meant code filtering
pull requests out had nothing to filter and code forgetting to filter them
looked right.

Two of that Recipe's sixteen cases still carry no date, and both are refusals
rather than oversights. Checking that a bad credential is rejected means
sending a credential-shaped header to somebody else's authentication endpoint
to watch it fail. Checking the rate-limit response means exhausting a rate
limit on purpose, which spends capacity belonging to everyone sharing this
address. The other fourteen needed nothing but a GET anybody can make.

## Sources

- Documentation: https://docs.docker.com/docker-hub/api/latest/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve dockerhub     # run it
cauldron verify dockerhub -v # check every claim
```
