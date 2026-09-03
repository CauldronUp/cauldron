# Workday

Emulates the Workday API (v1), for local development and tests.

**5 conformance cases, all of them checked against the live API on 2026-09-01.**

## What this Recipe found

Its gateway **tells anybody which services exist**. A
recognised prefix answers 401 and anything else 404, with no credential ever
distinguishing itself, so the list of Workday services is public whether or not
you have a tenant.

## Sources

- Documentation: https://community.workday.com/sites/default/files/file-hosting/restapi/index.html
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve workday     # run it
cauldron verify workday -v # check every claim
```
