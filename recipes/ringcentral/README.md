# RingCentral

Emulates the RingCentral API (v1.0), for local development and tests.

**17 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## Sources

- Documentation: https://developers.ringcentral.com/api-reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve ringcentral     # run it
cauldron verify ringcentral -v # check every claim
```
