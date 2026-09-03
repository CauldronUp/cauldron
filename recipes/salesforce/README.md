# Salesforce

Emulates the Salesforce API (v60.0), for local development and tests.

**9 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## Sources

- Documentation: https://developer.salesforce.com/docs/atlas.en-us.api_rest.meta/api_rest/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve salesforce     # run it
cauldron verify salesforce -v # check every claim
```
