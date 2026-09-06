# SerpApi

Emulates the SerpApi API (unversioned), for local development and tests.

**9 conformance cases, 8 checked against the live API on 2026-09-03.** The unchecked one is the paging case: it sends the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

It **times its own scrape and not the page** -- two
clocks about the request and none about the result.

## Sources

- Documentation: https://serpapi.com/search-api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve serpapi     # run it
cauldron verify serpapi -v # check every claim
```
