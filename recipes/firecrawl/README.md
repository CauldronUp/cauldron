# Firecrawl

Emulates the Firecrawl API (v1), for local development and tests.

**13 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A running crawl already has results. While `status` is `scraping`, `data` holds whatever pages have been fetched so far -- nothing in the response distinguishes "some of the pages" from "all of the pages" except that one status field, so a read taken early is a real, incomplete answer rather than an empty one. `completed` and `total` also move independently, and `total` is only an estimate that rises as the crawler discovers more links, so a progress bar computed from the two can go backwards and a loop waiting for them to become equal can wait forever.

Starting a crawl returns just an identifier -- no status, no counts, no data -- so code that reads a page count off the start response is reading a field that doesn't exist yet. And a failed page never shows up as a failure: the crawl completes, the completed count is quietly lower than the number of pages discovered, and nothing in the response lists which ones dropped out.

## Sources

- Documentation: https://docs.firecrawl.dev/api-reference/introduction
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve firecrawl     # run it
cauldron verify firecrawl -v # check every claim
```
