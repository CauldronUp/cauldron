# backblaze

Emulates the backblaze API (v3), for local development and tests.

**25 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

It **serves one bucket under two protocols that
disagree about credentials**. The native API collapses missing and wrong into one
`bad_auth_token` with an empty message; the S3 gateway splits the same pair three
ways.

## Sources

- Documentation: https://www.backblaze.com/apidocs/introduction-to-the-b2-native-api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve backblaze     # run it
cauldron verify backblaze -v # check every claim
```
