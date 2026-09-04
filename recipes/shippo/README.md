# Shippo

Emulates the Shippo API (2018-02-08), for local development and tests.

**10 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Buying a label is asynchronous: the transaction comes back `QUEUED` with no `label_url`, and the URL appears later. Code that reads `label_url` straight off the create response gets nothing -- and in test mode, where labels succeed instantly, it can get a URL on a good day and ship broken code to production. Address validation is similarly advisory: an address can come back with `validation_results` saying it isn't valid, and the request still succeeds, so anything checking only the HTTP status never sees the warning.

## Sources

- Documentation: https://docs.goshippo.com/shippoapi/public-api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve shippo     # run it
cauldron verify shippo -v # check every claim
```
