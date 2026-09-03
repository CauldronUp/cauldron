# Companies House

Emulates the Companies House API (1.0), for local development and tests.

**12 conformance cases, all of them checked against the live API on 2026-09-03.**

## What this Recipe found

House's, where **the ignored password is refused**: empty
gets 401, non-empty gets 400, for a field the guide says is not read.

## Sources

- Documentation: https://developer-specs.company-information.service.gov.uk/companies-house-public-data-api/reference/company-profile
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve companieshouse     # run it
cauldron verify companieshouse -v # check every claim
```
