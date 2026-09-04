# Langfuse

Emulates the Langfuse API (1), for local development and tests.

**7 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A deprecated endpoint answers 200 with everything requested, plus a `_deprecation` object naming the replacement endpoint and a `sunsetAt` date -- nothing about the status code or the data itself signals a problem, so a client reads the payload, works perfectly, and stops working on a date nobody read. Three endpoints here are already in that state: `GET /traces`, `GET /observations` and `GET /v2/scores` are all deprecated in favor of `/v2/observations` and `/v3/scores`, which carry no `_deprecation` marker at all -- the same resource lives at several versioned paths simultaneously, and only the response body says which generation you're looking at.

A trace's `observations` field is an array of ids, not objects, so code iterating it looking for names, models, or costs just iterates strings. And the paging envelope's `totalItems` and `totalPages` are genuinely different numbers under the same `meta` object -- a loop written against the wrong one either requests pages that don't exist or stops before reaching the end.

## Sources

- Documentation: https://cloud.langfuse.com/generated/api/openapi.yml
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve langfuse     # run it
cauldron verify langfuse -v # check every claim
```
