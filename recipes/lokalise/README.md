# Lokalise

Emulates the Lokalise API (api2), for local development and tests.

**14 conformance cases, 12 checked against the live API on 2026-09-02.**

## What this Recipe found

It **checks the parameter before the caller**: a
malformed project identifier is refused for its shape with no credential
consulted, and identically with a wrong one.

## Sources

- Documentation: https://developers.lokalise.com/reference/lokalise-rest-api
- Machine-readable description: https://developers.lokalise.com/openapi/lokalise-api-without-branches.yml, last checked 2026-09-02
  `cauldron drift lokalise` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve lokalise     # run it
cauldron verify lokalise -v # check every claim
```
