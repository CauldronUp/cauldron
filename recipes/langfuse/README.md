# Langfuse

Emulates the Langfuse API (1), for local development and tests.

**9 conformance cases, 2 checked against the live API.**

Struck live 2026-09-05 against cloud.langfuse.com, no account and no key -- and found that auth had never actually been enforced. This file declared `scheme: basic` with no `keys` at all, which this format treats as "route first, tighten auth later"; every one of its own cases sent a credential and every one would have been accepted regardless of what it sent, because nothing was ever compared against anything. A real key is declared now, along with the two real refusal sentences: "No authorization header" for a missing credential and "Invalid credentials. Confirm that you've configured the correct host." for a wrong one.

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
