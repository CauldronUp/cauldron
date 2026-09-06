# Firecrawl

Emulates the Firecrawl API (v1), for local development and tests.

**18 conformance cases, 3 checked against the live API on 2026-09-05.**

The resource cases cite documentation rather than an observation on a real account, because observing them would mean crawling somebody else's website. The three refusal cases at the end were struck live against api.firecrawl.dev.

## What this Recipe found

A running crawl already has results. While `status` is `scraping`, `data` holds whatever pages have been fetched so far -- nothing in the response distinguishes "some of the pages" from "all of the pages" except that one status field, so a read taken early is a real, incomplete answer rather than an empty one. `completed` and `total` also move independently, and `total` is only an estimate that rises as the crawler discovers more links, so a progress bar computed from the two can go backwards and a loop waiting for them to become equal can wait forever.

**The live probe found something it should not have.** A POST to `/v1/scrape` with no Authorization header at all, sent to check what an unauthenticated request receives, answered 200 and actually fetched `https://example.com` -- real content, real metadata, `"creditsUsed":1`, against an account this task has no relationship to. That contradicts this file's own documented claim that every request needs a valid key, and it happened by accident: no second attempt was made, at this route or at the multi-page `/v1/crawl`, specifically to avoid making the same mistake twice. The full account is in the Recipe's own header rather than smoothed over here. The read-only half of the same probe is uneventful: an unrouted path and a wrong method on a real path both answer Express's own default 404 before authentication is ever consulted, and an invented token is refused exactly as documented.

Starting a crawl returns just an identifier -- no status, no counts, no data -- so code that reads a page count off the start response is reading a field that doesn't exist yet. And a failed page never shows up as a failure: the crawl completes, the completed count is quietly lower than the number of pages discovered, and nothing in the response lists which ones dropped out.

## Sources

- Documentation: https://docs.firecrawl.dev/api-reference/introduction
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve firecrawl     # run it
cauldron verify firecrawl -v # check every claim
```
