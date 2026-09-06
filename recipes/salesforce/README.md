# Salesforce

Emulates the Salesforce API (v60.0), for local development and tests.

**13 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

No case here can carry a verified date, and that is itself worth stating: the REST Data API this Recipe models has no shared host to probe. Every real org answers at its own instance domain, issued only after an OAuth flow with a registered Connected App's client id and secret, which this file is not permitted to obtain. Checked 2026-09-05: login.salesforce.com, the one host every flow starts at, answers real OAuth errors at its token endpoint, but OAuth itself is out of this Recipe's scope; pointing the Data API path at that host directly instead answers a login-flow error of its own, evidence the data API genuinely lives elsewhere rather than something this file's surface can reach unauthenticated.

Every record carries an `attributes` object holding its type and URL, so code that iterates a record's own keys hits that first, and diffing two records finds them different because their URLs differ. Field names are capitalised (`Id`, `Name`), and a query answering `done: false` has to be followed through `nextRecordsUrl` -- Salesforce does not warn when a match exceeds the batch size, it just stops and hands back a prefix.

Creating a record answers with a lower-case `id`; reading the same record back answers `Id`. Same identifier, two spellings, one API.

## Sources

- Documentation: https://developer.salesforce.com/docs/atlas.en-us.api_rest.meta/api_rest/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve salesforce     # run it
cauldron verify salesforce -v # check every claim
```
