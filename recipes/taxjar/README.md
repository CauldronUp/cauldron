# TaxJar

Emulates the TaxJar API (v2), for local development and tests.

**13 conformance cases, 3 checked against the live API on 2026-09-05.**

TaxJar has a real sandbox, so the rate and tax cases still cite documentation rather than a token of their own. The credential shape needed no token at all, and checking it live found this Recipe's message incomplete.

## What this Recipe found

Every rate is a decimal between zero and one -- `combined_rate` of 0.08 means eight percent, not eight, and a display that multiplies by a hundred twice shows 800%. Tax also isn't one number: it's broken into state, county, city, and special-district amounts that sum to the total, with a jurisdiction that doesn't apply present at zero rather than absent, so summing the visible parts plus the total double-counts.

And `has_nexus` is the whole question: without nexus in the destination state the amount is zero and the whole breakdown is simply absent, which is a correct answer, not an error.

## What checking it live found

No credential and a garbage bearer token are byte-identical -- backwards from many providers, which tell the two apart -- and both name the exact route that was asked for: `"Not authorized for route 'GET /v2/rates/90210'"`. A path nothing declares gets its own 404, `"No such route '...'"`, resolved before any credential is judged at all.

## Sources

- Documentation: https://developers.taxjar.com/api/reference/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve taxjar     # run it
cauldron verify taxjar -v # check every claim
```
