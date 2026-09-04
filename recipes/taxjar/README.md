# TaxJar

Emulates the TaxJar API (v2), for local development and tests.

**8 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Every rate is a decimal between zero and one -- `combined_rate` of 0.08 means eight percent, not eight, and a display that multiplies by a hundred twice shows 800%. Tax also isn't one number: it's broken into state, county, city, and special-district amounts that sum to the total, with a jurisdiction that doesn't apply present at zero rather than absent, so summing the visible parts plus the total double-counts.

And `has_nexus` is the whole question: without nexus in the destination state the amount is zero and the whole breakdown is simply absent, which is a correct answer, not an error.

## Sources

- Documentation: https://developers.taxjar.com/api/reference/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve taxjar     # run it
cauldron verify taxjar -v # check every claim
```
